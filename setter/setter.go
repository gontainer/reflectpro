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

package setter

import (
	"reflect"

	intReflect "github.com/gontainer/reflectpro/internal/reflect"
)

/*
Set assigns the value `val` to the field `field` on the struct `strct`.
If the fourth argument equals true, it converts the type whenever it is possible.
Unexported fields are supported.

	type Person struct {
		Name string
	}
	p := Person{}
	_ = setter.Set(&p, "Name", "Jane", false)
	fmt.Println(p) // {Jane}
*/
func Set(strct any, field string, val any, convert bool) error {
	return intReflect.Set(strct, field, val, convert)
}

type CallbackOption func(*callbackConfig)

type callbackConfig struct {
	convertType      bool
	convertToPointer bool
}

func ConvertType(cfg *callbackConfig) {
	cfg.convertType = true
}

func ConvertToPointer(cfg *callbackConfig) {
	cfg.convertToPointer = true

}

// SetByCallback invokes callback on the all fields of the given struct.
// If you want to modify the value of the given field, the callback must return that value and `ok == true`.
func SetByCallback(strct any, callback func(field reflect.StructField) (_ any, ok bool), opts ...CallbackOption) error {
	var cfg callbackConfig

	for _, o := range opts {
		o(&cfg)
	}

	return intReflect.SetByCallback(strct, callback, cfg.convertType, cfg.convertToPointer)
}
