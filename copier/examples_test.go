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

package copier_test

import (
	"fmt"

	"github.com/gontainer/reflectpro/copier"
)

func ExampleCopy_convertOK() {
	var (
		from = int(5) // uint is not assignable to int,
		to   uint     // but [copier.Copy] can convert the type
	)

	err := copier.Copy(from, &to, true)

	fmt.Println(to)
	fmt.Println(err)

	// Output:
	// 5
	// <nil>
}

func ExampleCopy_convertMap() {
	var (
		from = map[int64]any{0: "Jane", 1: "John"}
		to   map[int32]string // let's convert keys and values
	)

	err := copier.Copy(from, &to, true)

	fmt.Println(to)
	fmt.Println(err)

	// Output:
	// map[0:Jane 1:John]
	// <nil>
}

func ExampleCopy_ok() {
	var (
		from = 5 // the type of the variable `to` can be different from the type of the variable `from`
		to   any // as long as the value of the `from` is assignable to the `to`
	)

	err := copier.Copy(from, &to, false)

	fmt.Println(to)
	fmt.Println(err)

	// Output:
	// 5
	// <nil>
}

func ExampleCopy_err1() {
	var (
		from = int(5)
		to   uint
	)

	err := copier.Copy(from, &to, false)

	fmt.Println(to)
	fmt.Println(err)

	// Output:
	// 0
	// value of type int is not assignable to type uint
}

func ExampleCopy_err2() {
	var (
		from float32 = 5
		to   uint    = 0
	)

	err := copier.Copy(from, &to, false)

	fmt.Println(to)
	fmt.Println(err)

	// Output:
	// 0
	// value of type float32 is not assignable to type uint
}

func ExampleCopy_err3() {
	var (
		from *int
		to   *uint
	)

	err := copier.Copy(from, &to, false)

	fmt.Println(to)
	fmt.Println(err)

	// Output:
	// <nil>
	// value of type *int is not assignable to type *uint
}
