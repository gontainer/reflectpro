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
	"fmt"
	"reflect"

	intReflect "github.com/gontainer/reflectpro/internal/reflect"
)

func Func(fn any) (reflect.Value, error) {
	v := reflect.ValueOf(fn)
	if !v.IsValid() {
		return reflect.Value{}, fmt.Errorf("invalid func: %T", fn)
	}

	if v.Kind() != reflect.Func {
		return reflect.Value{}, fmt.Errorf("expected %s, %T given", reflect.Func.String(), fn)
	}

	return v, nil
}

const (
	invalidMethodErr = "invalid func (%T).%+q"
)

func Method(object any, method string) (reflect.Value, error) {
	obj := reflect.ValueOf(object)
	if !obj.IsValid() {
		return reflect.Value{}, fmt.Errorf("invalid method receiver: %T", object)
	}

	fn := obj.MethodByName(method)

	_, err := intReflect.ValueToKindChain(obj)
	if err != nil {
		return reflect.Value{}, err //nolint:wrapcheck
	}

	for !fn.IsValid() && (obj.Kind() == reflect.Ptr || obj.Kind() == reflect.Interface) {
		obj = obj.Elem()
		fn = obj.MethodByName(method)
	}

	if !fn.IsValid() {
		return reflect.Value{}, fmt.Errorf(invalidMethodErr, object, method)
	}

	return fn, nil
}

func MethodByName(val reflect.Value, method string) (reflect.Value, error) {
	fn := val.MethodByName(method)
	if !fn.IsValid() {
		return reflect.Value{}, fmt.Errorf(invalidMethodErr, val.Interface(), method)
	}

	return fn, nil
}
