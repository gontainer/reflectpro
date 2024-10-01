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
	"unsafe"
)

type FieldCallbackResult struct {
	value any
	set   bool
	stop  bool
}

func FieldCallbackResultSet(value any) FieldCallbackResult {
	return FieldCallbackResult{
		value: value,
		set:   true,
		stop:  false,
	}
}

func FieldCallbackResultDontSet() FieldCallbackResult {
	return FieldCallbackResult{
		value: nil,
		set:   false,
		stop:  false,
	}
}

func FieldCallbackResultStop() FieldCallbackResult {
	return FieldCallbackResult{
		value: nil,
		set:   false,
		stop:  true,
	}
}

type FieldCallback = func(_ reflect.StructField, value any) FieldCallbackResult

// IterateFields traverses the fields of a struct, applying the callback function.
// Parameters:
//   - strct: The struct to iterate over
//   - callback: Function to call for each field
//   - convert: If true, attempts type conversion
//   - convertToPtr: If true, converts values returned by the callback to pointers when required
func IterateFields(strct any, callback FieldCallback, convert bool, convertToPtr bool) (err error) {
	strType := ""

	defer func() {
		if err != nil {
			if strType != "" {
				err = fmt.Errorf("%s: %w", strType, err)
			}

			err = fmt.Errorf("IterateFields: %w", err)
		}
	}()

	reflectVal, chain, err := ReducedValueOf(strct)
	if err != nil {
		return err
	}

	var iterator func(
		reflectVal reflect.Value,
		callback FieldCallback,
		convert bool,
		convertToPtr bool,
	) error

	switch {
	case chain.equalTo(reflect.Struct):
		strType = fmt.Sprintf("%T", reflect.Zero(reflectVal.Type()).Interface())
		iterator = iterateStruct

	case chain.equalTo(reflect.Ptr, reflect.Struct):
		strType = fmt.Sprintf("%T", reflect.Zero(reflectVal.Elem().Type()).Interface())
		iterator = iteratePtrStruct

	case chain.equalTo(reflect.Ptr, reflect.Interface, reflect.Struct):
		strType = fmt.Sprintf("%T", reflect.Zero(reflectVal.Type()).Interface())
		iterator = iteratePtrInterfaceStruct

	default:
		if err := ptrToNilStructError(strct); err != nil {
			return err
		}

		return fmt.Errorf("expected struct or pointer to struct, %T given", strct)
	}

	if err := iterator(reflectVal, callback, convert, convertToPtr); err != nil {
		return err
	}

	return nil
}

func valueFromField(strct reflect.Value, i int) any { //nolint:ireturn
	f := strct.Field(i)

	if !f.CanSet() { // handle unexported fields
		if !f.CanAddr() {
			tmpReflectVal := reflect.New(strct.Type()).Elem()
			tmpReflectVal.Set(strct)
			f = tmpReflectVal.Field(i)
		}

		f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
	}

	return f.Interface()
}

//nolint:revive
func iterateStruct(reflectVal reflect.Value, callback FieldCallback, convert bool, convertToPtr bool) error {
	for i := 0; i < reflectVal.Type().NumField(); i++ {
		result := callback(reflectVal.Type().Field(i), valueFromField(reflectVal, i))

		if result.set {
			return fmt.Errorf("pointer is required to set fields")
		}

		if result.stop {
			return nil
		}
	}

	return nil
}

func iteratePtrStruct(reflectVal reflect.Value, callback FieldCallback, convert bool, convertToPtr bool) error {
	for i := 0; i < reflectVal.Elem().Type().NumField(); i++ {
		result := callback(reflectVal.Elem().Type().Field(i), valueFromField(reflectVal.Elem(), i))

		if result.set {
			f := reflectVal.Elem().Field(i)
			if !f.CanSet() {
				f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
			}

			newVal := result.value

			newRefVal, err := func() (reflect.Value, error) {
				if convertToPtr && f.Kind() == reflect.Ptr && (newVal != nil || reflect.ValueOf(newVal).Kind() != reflect.Ptr) {
					val, err := ValueOf(newVal, f.Type().Elem(), convert)
					if err != nil {
						return reflect.Value{}, err
					}

					ptr := reflect.New(val.Type())
					ptr.Elem().Set(val)

					return ptr, nil
				}

				return ValueOf(newVal, f.Type(), convert)
			}()

			if err != nil {
				return fmt.Errorf("field %d %+q: %w", i, reflectVal.Elem().Type().Field(i).Name, err)
			}

			f.Set(newRefVal)
		}

		if result.stop {
			return nil
		}
	}

	return nil
}

func iteratePtrInterfaceStruct(
	reflectVal reflect.Value,
	callback FieldCallback,
	convert bool,
	convertToPtr bool,
) error {
	v := reflectVal.Elem()
	tmp := reflect.New(v.Elem().Type())
	tmp.Elem().Set(v.Elem())

	var (
		stop = false
		set  = false
	)

	newCallback := func(f reflect.StructField, value any) FieldCallbackResult {
		if stop {
			return FieldCallbackResult{
				value: nil,
				set:   false,
				stop:  true,
			}
		}

		result := callback(f, value)

		if result.stop {
			stop = true
		}

		if result.set {
			set = true
		}

		return result
	}

	if err := IterateFields(tmp.Interface(), newCallback, convert, convertToPtr); err != nil {
		return err
	}

	if set {
		v.Set(tmp.Elem())
	}

	return nil
}
