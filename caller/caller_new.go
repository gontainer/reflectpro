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
	"github.com/gontainer/reflectpro/caller/internal/caller"
	intReflect "github.com/gontainer/reflectpro/internal/reflect"
)

//nolint
// TODO format errors

func isPtr(v any) bool {
	chain, err := intReflect.ValueToKindChain(reflect.ValueOf(v))

	return err == nil && chain.Prefixed(reflect.Ptr)
}

//nolint:wrapcheck
func callMethod(object any, method string, args []any, convertArgs bool, v caller.FuncValidator) (_ []any, err error) {
	defer func() {
		if err != nil {
			err = grouperror.Prefix(fmt.Sprintf("cannot call method (%T).%+q: ", object, method), err)
		}
	}()

	fn, err := caller.Method(object, method)
	if err != nil {
		if errors.Is(err, caller.ErrInvalidMethod) && isPtr(object) {
			return caller.ValidateAndForceCallMethod(object, method, args, convertArgs, v)
		}

		return nil, err
	}

	if v != nil {
		if err := v.Validate(fn); err != nil {
			return nil, err
		}
	}

	return caller.CallFunc(fn, args, convertArgs)
}

// NewCallMethod replace will replace [CallMethod] and [ForceCallMethod].
func NewCallMethod(object any, method string, args []any, convertArgs bool) (_ []any, err error) {
	return callMethod(object, method, args, convertArgs, caller.DontValidate)
}

// NewCallProviderMethod will replace [CallProviderMethod] and [ForceCallProviderMethod].
func NewCallProviderMethod( //nolint:ireturn
	object any,
	method string,
	args []any,
	convertArgs bool,
) (
	_ any,
	err error,
) {
	//nolint
	// TODO different error type for internal and external errors
	results, err := callMethod(object, method, args, convertArgs, caller.ValidatorProvider)
	if err != nil {
		//nolint:wrapcheck
		return nil, grouperror.Prefix(fmt.Sprintf(providerMethodInternalErrPrefix, object, method), err)
	}

	var e error

	if len(results) > 1 {
		// do not panic when results[1] == nil
		e, _ = results[1].(error)
	}

	if e != nil {
		e = grouperror.Prefix(providerExternalErrPrefix, newProviderError(e))
	}

	return results[0], e //nolint:wrapcheck
}
