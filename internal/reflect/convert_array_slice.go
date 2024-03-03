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

func convertSlice(
	from reflect.Value,
	to reflect.Type,
	convert generalConverter,
) (
	_ reflect.Value,
	supports bool,
	_ error,
) {
	if isConvertibleSliceOrArray(from, to) {
		v, err := convertSliceOrArray(from, to, convert)

		return v, true, err
	}

	return reflect.Value{}, false, nil
}

func isArrayOrSlice(k reflect.Kind) bool {
	switch k { //nolint:exhaustive
	case
		reflect.Slice,
		reflect.Array:
		return true
	}

	return false
}

func isConvertibleSliceOrArray(from reflect.Value, to reflect.Type) bool {
	if !isArrayOrSlice(from.Kind()) || !isArrayOrSlice(to.Kind()) {
		return false
	}

	return true
}

//nolint:cyclop
func convertSliceOrArray(from reflect.Value, to reflect.Type, convert generalConverter) (reflect.Value, error) {
	// check whether slice values are convertible for len == 0
	if from.Len() == 0 && !isAny(from.Type().Elem()) && !isAny(to.Elem()) {
		if _, err := convert(
			reflect.Zero(from.Type().Elem()).Interface(),
			to.Elem(),
		); err != nil {
			return reflect.Value{}, err
		}
	}

	if from.Kind() == reflect.Slice && from.IsNil() {
		// zero value for slice == nil
		return reflect.Zero(to), nil
	}

	var (
		cp    reflect.Value
		toLen int
	)

	if to.Kind() == reflect.Array {
		toLen = to.Len()
		if from.Len() < toLen {
			toLen = from.Len()
		}

		cp = reflect.New(reflect.ArrayOf(to.Len(), to.Elem())).Elem()
	} else {
		toLen = from.Len()
		cp = reflect.MakeSlice(to, from.Len(), from.Len())
	}

	for i := 0; i < toLen; i++ {
		item := from.Index(i)
		for item.Kind() == reflect.Interface {
			item = item.Elem()
		}

		var currVal any

		if item.IsValid() {
			currVal = item.Interface()
		}

		curr, err := convert(currVal, to.Elem())
		if err != nil {
			return reflect.Value{}, fmt.Errorf("#%d: %w", i, err)
		}

		cp.Index(i).Set(curr)
	}

	return cp, nil
}
