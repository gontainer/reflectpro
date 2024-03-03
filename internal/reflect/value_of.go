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

package reflect

import (
	"fmt"
	"reflect"
)

/*
ValueOf is an extended version of [reflect.ValueOf].

Built-in [reflect.ValueOf](nil) returns the zero [reflect.Value].
ValueOf for the `i == nil` and a nil-able [reflect.Kind] of `to` returns a zero value from `to`.

If the third argument equals true, it converts the type whenever it is possible,
If `result.Type()` is not assignable to `to` it returns an error.
*/
func ValueOf(i any, to reflect.Type, convertVal bool) (reflect.Value, error) {
	if convertVal {
		return convert(i, to)
	}

	r := reflect.ValueOf(i)
	if !r.IsValid() {
		return zeroForNilable(i, to)
	}

	if err := assignable(r.Type(), to); err != nil {
		return reflect.Value{}, err
	}

	return r, nil
}

func zeroForNilable(i any, t reflect.Type) (reflect.Value, error) {
	if i == nil && isNilable(t.Kind()) {
		return reflect.Zero(t), nil
	}

	return reflect.Value{}, fmt.Errorf("cannot convert %T to %s", i, t.String())
}

func isNilable(k reflect.Kind) bool {
	switch k { //nolint:exhaustive
	case
		reflect.Chan,
		reflect.Func,
		reflect.Map,
		reflect.Ptr,
		reflect.UnsafePointer,
		reflect.Interface,
		reflect.Slice:
		return true
	}

	return false
}

func assignable(from, to reflect.Type) error {
	if !from.AssignableTo(to) {
		return fmt.Errorf("value of type %s is not assignable to type %s", from.String(), to.String())
	}

	return nil
}
