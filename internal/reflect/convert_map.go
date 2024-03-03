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

//nolint:cyclop
func convertMap(
	from reflect.Value,
	to reflect.Type,
	convert generalConverter,
) (
	_ reflect.Value,
	supports bool,
	_ error,
) {
	if from.Type().Kind() != reflect.Map || to.Kind() != reflect.Map {
		supports = false
		return //nolint:nakedret,nlreturn
	}

	if from.Len() == 0 { //nolint:nestif
		if !isAny(from.Type().Key()) || !isAny(to.Key()) {
			if _, err := convert(
				reflect.Zero(from.Type().Key()).Interface(),
				to.Key(),
			); err != nil {
				return reflect.Value{}, true, fmt.Errorf("non convertible keys: %w", err)
			}
		}

		if !isAny(from.Type().Elem()) || !isAny(to.Elem()) {
			if _, err := convert(
				reflect.Zero(from.Type().Elem()).Interface(),
				to.Elem(),
			); err != nil {
				return reflect.Value{}, true, fmt.Errorf("non convertible values: %w", err)
			}
		}
	}

	if from.IsNil() {
		return reflect.Zero(to), true, nil
	}

	mapType := reflect.MapOf(to.Key(), to.Elem())
	result := reflect.MakeMapWithSize(mapType, from.Len())

	iter := from.MapRange()
	for iter.Next() {
		newKey, err := convert(iter.Key().Interface(), to.Key())
		if err != nil {
			err = fmt.Errorf("map key: %w", err)

			return reflect.Value{}, true, err
		}

		newValue, err := convert(iter.Value().Interface(), to.Elem())
		if err != nil {
			err = fmt.Errorf("map value: %w", err)

			return reflect.Value{}, true, err
		}

		result.SetMapIndex(newKey, newValue)
	}

	return result, true, nil
}
