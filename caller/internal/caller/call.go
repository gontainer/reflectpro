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

//nolint:wrapcheck
func CallMethod(
	object any,
	method string,
	args []any,
	convertArgs bool,
	validator FuncValidator,
) (
	_ []any,
	err error,
) {
	defer func() {
		if err != nil {
			err = grouperror.Prefix(fmt.Sprintf("cannot call method (%T).%+q: ", object, method), err)
		}
	}()

	fn, err := Method(object, method)
	if err != nil {
		if errors.Is(err, ErrInvalidMethod) && isPtr(object) {
			return validateAndForceCallMethod(object, method, args, convertArgs, validator)
		}

		return nil, err
	}

	if validator != nil {
		if err := validator.Validate(fn); err != nil {
			return nil, err
		}
	}

	return CallFunc(fn, args, convertArgs)
}

// validateAndCallFunc validates and calls the given func.
//
// fn.Kind() MUST BE equal to [reflect.Func].
func validateAndCallFunc(fn reflect.Value, args []any, convertArgs bool, v FuncValidator) ([]any, error) {
	if v != nil {
		if err := v.Validate(fn); err != nil {
			return nil, err //nolint:wrapcheck
		}
	}

	return CallFunc(fn, args, convertArgs)
}

//nolint:wrapcheck,cyclop
func validateAndForceCallMethod(
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

	val, chain, err := intReflect.ReducedValue(val)
	if err != nil {
		return nil, err
	}

	// see [intReflect.set]
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
		fn, err := MethodByReflect(val, method)
		if err != nil {
			return nil, err
		}

		return validateAndCallFunc(fn, args, convertArgs, v)
	}

	if len(chain) == 3 && chain.Prefixed(reflect.Ptr, reflect.Interface) {
		if chain[2] == reflect.Invalid {
			return nil, ErrInvalidObject
		}

		cp := reflect.New(val.Elem().Elem().Type())
		cp.Elem().Set(val.Elem().Elem())

		fn, err := MethodByReflect(cp, method)
		if err != nil {
			return nil, err
		}

		res, err := validateAndCallFunc(fn, args, convertArgs, v)
		if err == nil {
			val.Elem().Set(cp.Elem())
		}

		return res, err
	}

	panic("validateAndForceCallMethod: unexpected error") // this should be unreachable
}

func isPtr(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Ptr
}
