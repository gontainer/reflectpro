// Copyright (c) 2023–present Bartłomiej Krukowski
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package caller

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gontainer/grouperror"
	intReflect "github.com/gontainer/reflectpro/internal/reflect"
)

type reflectType struct {
	reflect.Type
}

// In works almost same as reflect.Type.In,
// but it returns t.In(t.NumIn() - 1).Elem() for t.isVariadic() && i >= t.NumIn().
func (t reflectType) In(i int) reflect.Type {
	last := t.NumIn() - 1
	if i > last {
		i = last
	}

	r := t.Type.In(i)

	if t.IsVariadic() && i == last {
		r = r.Elem()
	}

	return r
}

// ValidateAndCallFunc validates and calls the given func.
//
// fn.Kind() MUST BE equal to [reflect.Func].
func ValidateAndCallFunc(fn reflect.Value, args []any, convertArgs bool, v FuncValidator) ([]any, error) {
	if v != nil {
		if err := v.Validate(fn); err != nil {
			return nil, err //nolint:wrapcheck
		}
	}

	return CallFunc(fn, args, convertArgs)
}

// CallFunc calls the given func.
//
// fn.Kind() MUST BE equal to [reflect.Func].
func CallFunc(fn reflect.Value, args []any, convertArgs bool) ([]any, error) {
	fnType := reflectType{fn.Type()}

	if len(args) > fnType.NumIn() && !fnType.IsVariadic() {
		return nil, errors.New("too many input arguments")
	}

	minParams := fnType.NumIn()
	if fnType.IsVariadic() {
		minParams--
	}

	if len(args) < minParams {
		return nil, errors.New("not enough input arguments")
	}

	var (
		argsVals = make([]reflect.Value, len(args))
		errs     = make([]error, 0, len(args))
	)

	for i, p := range args {
		var (
			convertTo = fnType.In(i)
			err       error
		)

		argsVals[i], err = intReflect.ValueOf(p, convertTo, convertArgs)

		if err != nil {
			errs = append(errs, grouperror.Prefix(fmt.Sprintf("arg%d: ", i), err))
		}
	}

	if len(errs) > 0 {
		return nil, grouperror.Join(errs...) //nolint:wrapcheck
	}

	var result []any

	if fn.Type().NumOut() > 0 {
		result = make([]any, fn.Type().NumOut())
	}

	for i, v := range fn.Call(argsVals) {
		result[i] = v.Interface()
	}

	return result, nil
}

//nolint:wrapcheck,cyclop
func ValidateAndForceCallMethod(
	object any,
	method string,
	args []any,
	convertArgs bool,
	v FuncValidator,
) (
	[]any,
	error,
) {
	val := reflect.ValueOf(object)
	if val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("expected %s, %T given", reflect.Ptr.String(), object)
	}

	chain, err := intReflect.ValueToKindChain(val)
	if err != nil {
		return nil, err
	}

	// see [intReflect.Set]
	for {
		switch {
		case chain.Prefixed(reflect.Ptr, reflect.Ptr):
			val = val.Elem()
			chain = chain[1:]

			continue
		case chain.Prefixed(reflect.Ptr, reflect.Interface, reflect.Ptr):
			val = val.Elem().Elem()
			chain = chain[2:]

			continue
		}

		break
	}

	if len(chain) == 2 && chain.Prefixed(reflect.Ptr) {
		fn, err := MethodByName(val, method)
		if err != nil {
			return nil, err
		}

		return ValidateAndCallFunc(fn, args, convertArgs, v)
	}

	if len(chain) == 3 && chain.Prefixed(reflect.Ptr, reflect.Interface) {
		if chain[2] == reflect.Invalid {
			return nil, errors.New("invalid object")
		}

		cp := reflect.New(val.Elem().Elem().Type())
		cp.Elem().Set(val.Elem().Elem())

		fn, err := MethodByName(cp, method)
		if err != nil {
			return nil, err
		}

		res, err := ValidateAndCallFunc(fn, args, convertArgs, v)
		if err == nil {
			val.Elem().Set(cp.Elem())
		}

		return res, err
	}

	panic("ValidateAndForceCallMethod: unexpected error") // this should be unreachable
}
