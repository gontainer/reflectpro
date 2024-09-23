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
	"errors"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/gontainer/grouperror"
)

func ptrToNilStructError(v any) error {
	t := reflect.TypeOf(v)
	if t != nil {
		for t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		if t.Kind() == reflect.Struct {
			return errors.New("pointer to nil struct given")
		}
	}

	return nil
}

func fieldNotSupportedError(field string) error {
	if field == "_" {
		return errors.New(`"_" is not supported`)
	}

	return nil
}

func fieldByName(strct reflect.Value, field string) (f reflect.Value, err error) {
	f = strct.FieldByName(field)
	if !f.IsValid() {
		err = fmt.Errorf("field %+q does not exist", field)
	}

	return
}

func Get(strct any, field string) (_ any, err error) { //nolint:cyclop,ireturn
	defer func() {
		if err != nil {
			err = grouperror.Prefix(fmt.Sprintf("get (%T).%+q: ", strct, field), err)
		}
	}()

	if err := fieldNotSupportedError(field); err != nil {
		return nil, err
	}

	reflectVal := reflect.ValueOf(strct)

	chain, err := ValueToKindChain(reflectVal)
	if err != nil {
		return nil, err
	}

	for len(chain) > 1 {
		chain = chain[1:]
		reflectVal = reflectVal.Elem()
	}

	if reflectVal.Kind() != reflect.Struct {
		if !reflectVal.IsValid() {
			if err := ptrToNilStructError(strct); err != nil {
				return nil, err
			}
		}

		return nil, fmt.Errorf("expected struct, %T given", strct)
	}

	f, err := fieldByName(reflectVal, field)
	if err != nil {
		return nil, err
	}

	if !f.CanSet() { // handle unexported fields
		if !f.CanAddr() {
			tmpReflectVal := reflect.New(reflectVal.Type()).Elem()
			tmpReflectVal.Set(reflectVal)
			f = tmpReflectVal.FieldByName(field)
		}

		f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
	}

	return f.Interface(), nil
}

type FieldCallback = func(_ reflect.StructField, value any) (_ any, ok bool)

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

	reflectVal := reflect.ValueOf(strct)
	if !reflectVal.IsValid() {
		return fmt.Errorf("expected struct, %T given", strct)
	}

	chain, err := ValueToKindChain(reflectVal)
	if err != nil {
		return err
	}

	// TODO
	//if chain[len(chain)-1] != reflect.Struct {
	//	return fmt.Errorf("expected struct, %T given", strct)
	//}

	// see [Set]
	for {
		switch {
		case chain.Prefixed(reflect.Ptr, reflect.Ptr):
			reflectVal = reflectVal.Elem()
			chain = chain[1:]

			continue
		case chain.Prefixed(reflect.Ptr, reflect.Interface, reflect.Ptr):
			reflectVal = reflectVal.Elem().Elem()
			chain = chain[2:]

			continue
		}

		break
	}

	valueFromField := func(strct reflect.Value, i int) any {
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

	switch {
	case chain.equalTo(reflect.Struct):
		strType = fmt.Sprintf("%T", reflect.Zero(reflectVal.Type()).Interface())
		for i := 0; i < reflectVal.Type().NumField(); i++ {
			if _, ok := callback(reflectVal.Type().Field(i), valueFromField(reflectVal, i)); ok {
				return fmt.Errorf("pointer is required to set fields")
			}
		}

	case chain.equalTo(reflect.Ptr, reflect.Struct):
		strType = fmt.Sprintf("%T", reflect.Zero(reflectVal.Elem().Type()).Interface())
		for i := 0; i < reflectVal.Elem().Type().NumField(); i++ {
			if newVal, ok := callback(reflectVal.Elem().Type().Field(i), valueFromField(reflectVal.Elem(), i)); ok {
				f := reflectVal.Elem().Field(i)
				if !f.CanSet() {
					f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
				}

				newRefVal, err := func() (reflect.Value, error) {
					if convertToPtr && f.Kind() == reflect.Ptr && (newVal != nil || reflect.ValueOf(newVal).Kind() != reflect.Ptr) {
						val, err := ValueOf(newVal, f.Elem().Type(), convert)
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
		}

	case chain.equalTo(reflect.Ptr, reflect.Interface, reflect.Struct):
		strType = fmt.Sprintf("%T", reflect.Zero(reflectVal.Type()).Interface())
		// TODO remove recursion
		v := reflectVal.Elem()
		tmp := reflect.New(v.Elem().Type())
		tmp.Elem().Set(v.Elem())
		if err := IterateFields(tmp.Interface(), callback, convert, convertToPtr); err != nil {
			return err
		}
		v.Set(tmp.Elem())

	default:
		if err := ptrToNilStructError(strct); err != nil {
			return err
		}

		return fmt.Errorf("expected struct or pointer to struct, %T given", strct)
	}

	return nil
}

//nolint:cyclop
func Set(strct any, field string, val any, convert bool) (err error) {
	defer func() {
		if err != nil {
			err = grouperror.Prefix(fmt.Sprintf("set (%T).%+q: ", strct, field), err)
		}
	}()

	if err := fieldNotSupportedError(field); err != nil {
		return err
	}

	reflectVal := reflect.ValueOf(strct)

	chain, err := ValueToKindChain(reflectVal)
	if err != nil {
		return err
	}

	/*
		removes prepending duplicate Ptr & Interface elements
		e.g.:
			s := &struct{ val int }{}
			Set(&s, ... // chain == {Ptr, Ptr, Struct}

		or:
			var s any = &struct{ val int }{}
			var s2 any = &s
			var s3 any = &s
			Set(&s3, ... // chain == {Ptr, Interface, Ptr, Interface, Ptr, Interface, Struct}
	*/
	for {
		switch {
		case chain.Prefixed(reflect.Ptr, reflect.Ptr):
			reflectVal = reflectVal.Elem()
			chain = chain[1:]

			continue
		case chain.Prefixed(reflect.Ptr, reflect.Interface, reflect.Ptr):
			reflectVal = reflectVal.Elem().Elem()
			chain = chain[2:]

			continue
		}

		break
	}

	switch {
	// s := struct{ val int }{}
	// Set(&s...
	case chain.equalTo(reflect.Ptr, reflect.Struct):
		return setOnValue(
			reflectVal.Elem(),
			field,
			val,
			convert,
		)

	// var s any = struct{ val int }{}
	// Set(&s...
	case chain.equalTo(reflect.Ptr, reflect.Interface, reflect.Struct):
		v := reflectVal.Elem()
		tmp := reflect.New(v.Elem().Type()).Elem()
		tmp.Set(v.Elem())

		if err := setOnValue(tmp, field, val, convert); err != nil {
			return err
		}

		v.Set(tmp)

		return nil

	default:
		if err := ptrToNilStructError(strct); err != nil {
			return err
		}

		return fmt.Errorf("expected pointer to struct, %T given", strct)
	}
}

func setOnValue(strct reflect.Value, field string, val any, convert bool) error {
	f, err := fieldByName(strct, field)
	if err != nil {
		return err
	}

	v, err := ValueOf(val, f.Type(), convert)
	if err != nil {
		return err
	}

	if !f.CanSet() { // handle unexported fields
		f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
	}

	f.Set(v)

	return nil
}

type kindChain []reflect.Kind

func (c kindChain) equalTo(kinds ...reflect.Kind) bool {
	if len(c) != len(kinds) {
		return false
	}

	for i := 0; i < len(c); i++ {
		if c[i] != kinds[i] {
			return false
		}
	}

	return true
}

func (c kindChain) Prefixed(kinds ...reflect.Kind) bool {
	if len(c) < len(kinds) {
		return false
	}

	return c[:len(kinds)].equalTo(kinds...)
}

func ValueToKindChain(v reflect.Value) (kindChain, error) { //nolint // unexported-return
	var r kindChain
	r = make(kindChain, 0, 5) //nolint:gomnd
	ptrs := make(map[uintptr]struct{})

	for {
		if v.Kind() == reflect.Ptr && !v.IsNil() {
			ptr := v.Elem().UnsafeAddr()
			if _, ok := ptrs[ptr]; ok {
				return nil, errors.New("unexpected pointer loop")
			}

			ptrs[ptr] = struct{}{}
		}

		r = append(r, v.Kind())

		if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
			v = v.Elem()

			continue
		}

		break
	}

	return r, nil
}
