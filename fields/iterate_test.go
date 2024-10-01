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

package fields_test

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/gontainer/reflectpro/fields"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type CustomString string

type Person struct {
	Name string
}

type Employee struct {
	Person
	Role string
}

type TeamMeta struct {
	Name string
}

type Team struct {
	Lead Employee
	TeamMeta
}

type C struct {
	D string
}

type B struct {
	C C
}

type A struct {
	B B
}

type XX struct {
	_ int
	_ string
}

type YY struct {
	*XX
}

func setValueByFieldIndex(ptrStruct any, fieldIndex int, value any) {
	f := reflect.ValueOf(ptrStruct).Elem().Field(fieldIndex)
	f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
	f.Set(reflect.ValueOf(value))
}

func newXXWithBlankValues(t *testing.T, first int, second string) *XX { //nolint:thelper
	x := XX{}
	setValueByFieldIndex(&x, 0, first)
	setValueByFieldIndex(&x, 1, second)

	buff := bytes.NewBuffer(nil)
	_, err := fmt.Fprint(buff, x)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("{%d %s}", first, second), buff.String())

	return &x
}

//nolint:gocognit,goconst,lll
func TestIterate(t *testing.T) {
	t.Parallel()

	t.Run("Setter", func(t *testing.T) {
		scenarios := []struct {
			name    string
			options []fields.Option
			input   any
			output  any
			error   string
		}{
			{
				name: "Person OK",
				options: []fields.Option{
					fields.Setter(func(_ fields.Path, _ any) (_ any, set bool) {
						return "Jane", true
					}),
				},
				input: Person{},
				output: Person{
					Name: "Jane",
				},
				error: "",
			},
			{
				name: "Person OK (convert types)",
				options: []fields.Option{
					fields.Setter(func(_ fields.Path, _ any) (_ any, set bool) {
						return CustomString("Jane"), true
					}),
					fields.ConvertTypes(true),
				},
				input: Person{},
				output: Person{
					Name: "Jane",
				},
				error: "",
			},
			{
				name: "Person error (convert types)",
				options: []fields.Option{
					fields.Setter(func(_ fields.Path, value any) (_ any, set bool) {
						return CustomString("Jane"), true
					}),
				},
				input: Person{},
				output: Person{
					Name: "Jane",
				},
				error: "fields.Iterate: IterateFields: *interface {}: IterateFields: fields_test.Person: field 0 \"Name\": value of type fields_test.CustomString is not assignable to type string",
			},
			{
				name: "A.B.C.D OK",
				options: []fields.Option{
					fields.Setter(func(path fields.Path, value any) (_ any, set bool) {
						if path.EqualNames("B", "C", "D") {
							return "Hello", true
						}

						return nil, false
					}),
					fields.Recursive(true),
				},
				input: A{},
				output: A{
					B: B{
						C: C{
							D: "Hello",
						},
					},
				},
				error: "",
			},
			{
				name: "A.B.C.D error (convert types)",
				options: []fields.Option{
					fields.Setter(func(path fields.Path, value any) (_ any, set bool) {
						if path.EqualNames("B", "C", "D") {
							return 5, true
						}

						return nil, false
					}),
					fields.Recursive(true),
				},
				input:  A{},
				output: A{},
				error:  `fields.Iterate: B: C: IterateFields: *interface {}: IterateFields: fields_test.C: field 0 "D": value of type int is not assignable to type string`,
			},
			{
				name: "Employee (embedded)",
				options: []fields.Option{
					fields.Setter(func(path fields.Path, value any) (_ any, set bool) {
						switch {
						case path.EqualNames("Person", "Name"):
							return "Jane", true
						case path.EqualNames("Role"):
							return "Lead", true
						}

						return nil, false
					}),
					fields.Recursive(true),
				},
				input: Employee{},
				output: Employee{
					Person: Person{
						Name: "Jane",
					},
					Role: "Lead",
				},
				error: "",
			},
			{
				name: "Team #1",
				options: []fields.Option{
					fields.Setter(func(path fields.Path, value any) (_ any, set bool) {
						switch {
						case path.EqualNames("Lead", "Person", "Name"):
							return "Jane", true
						case path.EqualNames("Lead", "Role"):
							return "Lead", true
						case path.EqualNames("TeamMeta", "Name"):
							return "Hawkeye", true
						}

						return nil, false
					}),
					fields.Recursive(true),
				},
				input: Team{},
				output: Team{
					Lead: Employee{
						Person: Person{
							Name: "Jane",
						},
						Role: "Lead",
					},
					TeamMeta: TeamMeta{
						Name: "Hawkeye",
					},
				},
				error: "",
			},
			{
				name: "Team #2",
				options: []fields.Option{
					fields.Setter(func(path fields.Path, value any) (_ any, set bool) {
						switch {
						case path.EqualNames("Lead", "Role"):
							return "Lead", true
						case path.EqualNames("Lead"):
							return Employee{
								Person: Person{
									Name: "Jane",
								},
								Role: "Lead",
							}, true
						case path.EqualNames("TeamMeta", "Name"):
							return "Hawkeye", true
						}

						return nil, false
					}),
					fields.Recursive(true),
				},
				input: Team{},
				output: Team{
					Lead: Employee{
						Person: Person{
							Name: "Jane",
						},
						Role: "Lead",
					},
					TeamMeta: TeamMeta{
						Name: "Hawkeye",
					},
				},
				error: "",
			},
			{
				name: "YY",
				options: []fields.Option{
					fields.Setter(func(path fields.Path, value any) (_ any, set bool) {
						if path.EqualNames("XX") {
							return &XX{}, true
						}

						//nolint:exhaustive
						if path.EqualNames("XX", "_") {
							switch path[len(path)-1].Type.Kind() {
							case reflect.Int:
								return 5, true
							case reflect.String:
								return "five", true
							}
						}

						return nil, false
					}),
					fields.Recursive(true),
				},
				input: YY{},
				output: YY{
					XX: newXXWithBlankValues(t, 5, "five"),
				},
			},
			{
				name: "YY",
				options: []fields.Option{
					fields.Setter(func(path fields.Path, value any) (_ any, set bool) {
						if path.EqualNames("XX") {
							return &XX{}, true
						}

						//nolint:exhaustive
						if path.EqualNames("XX", "_") {
							switch path[len(path)-1].Type.Kind() {
							case reflect.Int:
								return 7, true
							case reflect.String:
								return "seven", true
							}
						}

						return nil, false
					}),
					fields.Recursive(true),
				},
				input: YY{},
				output: YY{
					XX: newXXWithBlankValues(t, 7, "seven"),
				},
			},
			{
				name:    "invalid input",
				options: nil,
				input:   100,
				output:  nil,
				error:   "fields.Iterate: IterateFields: expected struct or pointer to struct, *interface {} given",
			},
		}

		for _, s := range scenarios {
			s := s

			t.Run(s.name, func(t *testing.T) {
				t.Parallel()

				input := s.input
				err := fields.Iterate(&input, s.options...)

				if s.error != "" {
					require.EqualError(t, err, s.error)

					return
				}

				require.NoError(t, err)

				assert.Equal(t, s.output, input)
			})
		}
	})
}
