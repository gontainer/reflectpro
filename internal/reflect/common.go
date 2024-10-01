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

func reducedStructValueOf(strct any) (reflect.Value, kindChain, error) {
	reflectVal := reflect.ValueOf(strct)

	chain, err := ValueToKindChain(reflectVal)
	if err != nil {
		return reflect.Value{}, nil, err
	}

	/*
		removes prepending duplicate [reflect.Ptr] & [reflect.Interface] elements
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

	return reflectVal, chain, nil
}
