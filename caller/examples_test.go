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

package caller_test

import (
	"fmt"

	"github.com/gontainer/reflectpro/caller"
)

type Person struct {
	name string
}

func NewPerson(name string) *Person {
	return &Person{name: name}
}

func (p *Person) SetName(n string) {
	p.name = n
}

func (p Person) WithName(n string) Person {
	p.name = n

	return p
}

func ExampleCall_ok() {
	p := &Person{}
	_, _ = caller.Call(p.SetName, []any{"Mary"}, false)
	fmt.Println(p.name)
	// Output: Mary
}

func ExampleCall_returnValue() {
	fn := func(a int, b int) int {
		return a * b
	}
	r, _ := caller.Call(fn, []any{3, 2}, false)
	fmt.Println(r[0])
	// Output: 6
}

func ExampleCall_error() {
	fn := func(a int, b int) int {
		return a * b
	}
	_, err := caller.Call(fn, []any{"2", "2"}, true)

	fmt.Println(err)
	// Output:
	// cannot call func(int, int) int: arg0: cannot convert string to int
	// cannot call func(int, int) int: arg1: cannot convert string to int
}

func ExampleCall_error2() {
	fn := func(a int, b int) int {
		return a * b
	}
	_, err := caller.Call(fn, []any{"2", "2"}, false)
	fmt.Println(err)
	// Output:
	// cannot call func(int, int) int: arg0: value of type string is not assignable to type int
	// cannot call func(int, int) int: arg1: value of type string is not assignable to type int
}

func ExampleCallProvider() {
	p, _, _ := caller.CallProvider(NewPerson, []any{"Mary"}, false)
	fmt.Printf("%+v", p)
	// Output: &{name:Mary}
}

func ExampleCallMethod() {
	p := &Person{}
	_, _ = caller.CallMethod(p, "SetName", []any{"Mary"}, false)
	fmt.Println(p.name)
	// Output: Mary
}

func ExampleCallWither() {
	p := Person{}
	p2, _ := caller.CallWither(p, "WithName", []any{"Mary"}, false)
	fmt.Printf("%+v", p2)
	// Output: {name:Mary}
}
