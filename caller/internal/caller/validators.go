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
)

//nolint:gochecknoglobals
var (
	DontValidate      FuncValidator = nil //nolint
	ValidatorWither                 = ChainValidator{validateWither}
	ValidatorProvider               = ChainValidator{validateProvider}
)

//nolint:gochecknoglobals
var (
	errorInterface = reflect.TypeOf((*error)(nil)).Elem()
)

type FuncValidator interface {
	Validate(reflect.Value) error //nolint:inamedparam
}

type ChainValidator []func(reflect.Value) error

func (f ChainValidator) Validate(fn reflect.Value) error {
	for _, v := range f {
		if err := v(fn); err != nil {
			return err
		}
	}

	return nil
}

func validateWither(fn reflect.Value) error {
	if fn.Type().NumOut() != 1 {
		return fmt.Errorf("wither must return 1 value, given function returns %d values", fn.Type().NumOut())
	}

	return nil
}

func validateProvider(fn reflect.Value) error {
	if fn.Type().NumOut() == 0 || fn.Type().NumOut() > 2 {
		return fmt.Errorf(
			"provider must return 1 or 2 values, given function returns %d values",
			fn.Type().NumOut(),
		)
	}

	if fn.Type().NumOut() == 2 && !fn.Type().Out(1).Implements(errorInterface) {
		return fmt.Errorf(
			"second value returned by provider must implement error interface, %s given",
			fn.Type().Out(1).String(),
		)
	}

	return nil
}
