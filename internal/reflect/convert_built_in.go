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
	"reflect"
)

func convertBuiltIn(from reflect.Value, to reflect.Type, _ generalConverter) (_ reflect.Value, supports bool, _ error) {
	// avoid a panic, see [reflect.Type.ConvertibleTo]
	if to.Kind() == reflect.Array &&
		(from.Kind() == reflect.Slice || from.Kind() == reflect.Array) &&
		from.Len() < to.Len() {
		return reflect.Value{}, false, nil
	}

	if from.Type().ConvertibleTo(to) {
		return from.Convert(to), true, nil
	}

	return reflect.Value{}, false, nil
}

/*
isAny returns true for any interface with zero methods:

	type any1 interface{}
	type any2 = interface{}
*/
func isAny(v reflect.Type) bool {
	return v.Kind() == reflect.Interface && v.NumMethod() == 0
}
