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

// generalConverter is an alias for a signature of the [convert] function.
type generalConverter = func(value any, to reflect.Type) (reflect.Value, error)

type converter interface {
	convert(from reflect.Value, to reflect.Type, convert generalConverter) (_ reflect.Value, supports bool, _ error)
}

type converterFn func(
	from reflect.Value,
	to reflect.Type,
	convert generalConverter,
) (
	_ reflect.Value,
	supports bool,
	_ error,
)

func (fn converterFn) convert(
	from reflect.Value,
	to reflect.Type,
	convert generalConverter,
) (
	_ reflect.Value,
	supports bool,
	_ error,
) {
	return fn(from, to, convert)
}

//nolint:gochecknoglobals
var (
	converters = []converter{
		converterFn(convertBuiltIn),
		converterFn(convertSlice),
		converterFn(convertMap),
	}
)

// convert converts given value to given type whenever it is possible.
// In opposition to built-in reflect package it can convert slices and maps.
func convert(value any, to reflect.Type) (reflect.Value, error) {
	from := reflect.ValueOf(value)
	if !from.IsValid() {
		return zeroForNilable(value, to)
	}

	for _, c := range converters {
		if v, supports, err := c.convert(from, to, convert); supports {
			if err != nil {
				err = fmt.Errorf("cannot convert %s to %s: %w", from.Type().String(), to.String(), err)
			}

			return v, err
		}
	}

	return reflect.Value{}, fmt.Errorf("cannot convert %s to %s", from.Type().String(), to.String())
}
