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

package setter_test

import (
	"fmt"

	"github.com/gontainer/reflectpro/setter"
)

func ExampleSet_ok() {
	person := struct {
		name string
	}{}
	err := setter.Set(&person, "name", "Mary", false)
	fmt.Println(person.name)
	fmt.Println(err)
	// Output:
	// Mary
	// <nil>
}

func ExampleSet_errFieldDoesNotExists() {
	person := struct {
		name string
	}{}
	err := setter.Set(&person, "firstname", "Mary", false)
	fmt.Println(err)
	// Output: set (*struct { name string })."firstname": field "firstname" does not exist
}

func ExampleSet_errNilPtr() {
	type Person struct {
		name string //nolint:unused
	}

	var person Person

	err := setter.Set(person, "name", "Mary", false)
	fmt.Println(err)
	// Output: set (setter_test.Person)."name": pointer to nil struct given
}

func ExampleSet_typeMismatchingError() {
	type name string

	type Person struct {
		name string //nolint:unused
	}

	var person Person

	err := setter.Set(&person, "name", name("Jane"), false)
	fmt.Println(err)
	// Output: set (*setter_test.Person)."name": value of type setter_test.name is not assignable to type string
}

func ExampleSet_typeMismatchingConvert() {
	type name string

	type Person struct {
		name string
	}

	var person Person

	_ = setter.Set(&person, "name", name("Jane"), true)
	fmt.Println(person.name)
	// Output: Jane
}
