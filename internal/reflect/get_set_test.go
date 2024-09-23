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

package reflect_test

import (
	"fmt"
	stdReflect "reflect"
	"testing"

	"github.com/gontainer/reflectpro/internal/reflect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type Person struct {
		name string
	}

	t.Run("OK", func(t *testing.T) {
		t.Run("#1", func(t *testing.T) {
			t.Parallel()

			var p any = Person{name: "Mary"}
			n, err := reflect.Get(p, "name")
			assert.NoError(t, err)
			assert.Equal(t, "Mary", n)
		})
		t.Run("#2", func(t *testing.T) {
			t.Parallel()

			var s any = struct{ val int }{val: 5}
			x, err := reflect.Get(s, "val")
			assert.NoError(t, err)
			assert.Equal(t, 5, x)
		})
		t.Run("#3", func(t *testing.T) {
			t.Parallel()

			var p any = &Person{name: "Mary"}
			n, err := reflect.Get(p, "name")
			assert.NoError(t, err)
			assert.Equal(t, "Mary", n)
		})
		t.Run("#4", func(t *testing.T) {
			t.Parallel()

			var p any = &Person{name: "Mary"}
			n, err := reflect.Get(&p, "name")
			assert.NoError(t, err)
			assert.Equal(t, "Mary", n)
		})
		t.Run("#5", func(t *testing.T) {
			t.Parallel()

			p := Person{name: "Mary"}
			n, err := reflect.Get(p, "name")
			assert.NoError(t, err)
			assert.Equal(t, "Mary", n)
		})
		t.Run("#6", func(t *testing.T) {
			t.Parallel()

			p := &Person{name: "Mary"}
			n, err := reflect.Get(&p, "name")
			assert.NoError(t, err)
			assert.Equal(t, "Mary", n)
		})
		t.Run("#7", func(t *testing.T) {
			t.Parallel()

			p := &Person{name: "Mary"}
			var p2 any = &p
			n, err := reflect.Get(&p2, "name")
			assert.NoError(t, err)
			assert.Equal(t, "Mary", n)
		})
	})

	t.Run("Errors", func(t *testing.T) {
		t.Parallel()

		t.Run("#1", func(t *testing.T) {
			t.Parallel()

			_, err := reflect.Get(nil, "name")
			assert.EqualError(t, err, `get (<nil>)."name": expected struct, <nil> given`)
		})
		t.Run("#2", func(t *testing.T) {
			t.Parallel()

			_, err := reflect.Get(make([]int, 0), "name")
			assert.EqualError(t, err, `get ([]int)."name": expected struct, []int given`)
		})
		t.Run("#3", func(t *testing.T) {
			t.Parallel()

			_, err := reflect.Get(Person{}, "_")
			assert.EqualError(t, err, `get (reflect_test.Person)."_": "_" is not supported`)
		})
		t.Run("#4", func(t *testing.T) {
			t.Parallel()

			_, err := reflect.Get(Person{name: "Mary"}, "age")
			assert.EqualError(t, err, `get (reflect_test.Person)."age": field "age" does not exist`)
		})
		t.Run("#5", func(t *testing.T) {
			t.Parallel()

			var p any
			p = &p
			_, err := reflect.Get(p, "age")
			assert.EqualError(t, err, `get (*interface {})."age": unexpected pointer loop`)
		})
		t.Run("#6", func(t *testing.T) {
			t.Parallel()

			var a *struct {
				Name string
			}
			_, err := reflect.Get(&a, "name")
			assert.EqualError(t, err, `get (**struct { Name string })."name": pointer to nil struct given`)
		})
	})
}

//nolint:maintidx
func TestSet(t *testing.T) {
	t.Parallel()

	t.Run("blank identifier", func(t *testing.T) {
		t.Parallel()

		assert.EqualError(
			t,
			reflect.Set(struct{ _ int }{}, "_", 123, false),
			`set (struct { _ int })."_": "_" is not supported`,
		)
	})
	t.Run("anonymous struct", func(t *testing.T) {
		t.Parallel()

		p := struct {
			color string
		}{}
		assert.NoError(t, reflect.Set(&p, "color", "red", false))
		assert.Equal(t, "red", p.color)
	})
	t.Run("anonymous *struct", func(t *testing.T) {
		t.Parallel()

		p := &struct {
			color string
		}{}
		assert.NoError(t, reflect.Set(&p, "color", "brown", false))
		assert.Equal(t, "brown", p.color)
	})
	t.Run("***struct", func(t *testing.T) {
		t.Parallel()

		p := &struct {
			color string
		}{}
		p2 := &p
		p3 := &p2
		assert.NoError(t, reflect.Set(&p3, "color", "brown", false))
		assert.Equal(t, "brown", p.color)
	})
	t.Run("var a any", func(t *testing.T) {
		t.Parallel()

		t.Run("*struct{}", func(t *testing.T) {
			t.Parallel()

			const color = "red"
			p := struct {
				color string
			}{}
			var obj any = &p
			assert.Equal(t, "", p.color)
			assert.NoError(t, reflect.Set(obj, "color", color, false))
			assert.Equal(t, color, p.color)
		})
		t.Run("**struct{}", func(t *testing.T) {
			t.Parallel()

			const color = "blue"
			p := struct {
				color string
			}{}
			p2 := &p
			var obj any = &p2
			assert.Equal(t, "", p.color)
			assert.NoError(t, reflect.Set(obj, "color", color, false))
			assert.Equal(t, color, p.color)
		})
		t.Run("***struct{}", func(t *testing.T) {
			t.Parallel()

			const color = "yellow"
			p := struct {
				color string
			}{}
			p2 := &p
			p3 := &p2
			var obj any = &p3
			assert.Equal(t, "", p.color)
			assert.NoError(t, reflect.Set(obj, "color", color, false))
			assert.Equal(t, color, p.color)
		})
		t.Run("****struct{}", func(t *testing.T) {
			t.Parallel()

			const color = "green"
			p := struct {
				color string
			}{}
			p2 := &p
			p3 := &p2
			p4 := &p3
			var obj any = &p4
			assert.Equal(t, "", p.color)
			assert.NoError(t, reflect.Set(obj, "color", color, false))
			assert.Equal(t, color, p.color)
		})
	})
	t.Run("struct", func(t *testing.T) {
		t.Parallel()

		p := person{}
		assert.NoError(t, reflect.Set(&p, "Name", "Jane", false))
		assert.NoError(t, reflect.Set(&p, "age", 30, true))
		assert.Equal(t, person{Name: "Jane", age: 30}, p)
	})
	t.Run("*struct", func(t *testing.T) {
		t.Parallel()

		p := &person{}
		assert.NoError(t, reflect.Set(&p, "Name", "Mary", false))
		assert.NoError(t, reflect.Set(&p, "age", uint(33), true))
		assert.Equal(t, &person{Name: "Mary", age: 33}, p)
	})
	t.Run("var a any = &struct{}", func(t *testing.T) {
		t.Parallel()

		var p any = &person{}
		assert.NoError(t, reflect.Set(&p, "Name", "Bernhard Riemann", false))
		assert.NoError(t, reflect.Set(&p, "age", 45, true))
		assert.Equal(t, &person{Name: "Bernhard Riemann", age: 45}, p)
	})
	t.Run("var a any = struct{}", func(t *testing.T) {
		t.Parallel()

		var p any = person{}
		assert.NoError(t, reflect.Set(&p, "Name", "Jane", false))
		assert.Equal(t, person{Name: "Jane"}, p)
	})
	t.Run("var a any = struct{}; a2 := &a; setter.Set(&a2...", func(t *testing.T) {
		t.Parallel()

		var p any = person{}
		p2 := &p
		assert.NoError(t, reflect.Set(&p2, "Name", "Jane", false))
		assert.Equal(t, person{Name: "Jane"}, p)
	})
	t.Run("var a1 any = struct{}; var a2 any = &a1; var a3 any = &a2; ...; setter.Set(&aN...", func(t *testing.T) {
		t.Parallel()

		var p any = person{}
		p2 := &p
		var p3 any = &p2
		var p4 any = &p3
		var p5 any = &p4
		assert.NoError(t, reflect.Set(&p5, "Name", "Jane", false))
		assert.Equal(t, person{Name: "Jane"}, p)
	})
	t.Run("loop #1", func(t *testing.T) {
		t.Parallel()

		var p any
		p = &p
		assert.EqualError(
			t,
			reflect.Set(&p, "Name", "Jane", false),
			`set (*interface {})."Name": unexpected pointer loop`,
		)
	})
	t.Run("loop #2", func(t *testing.T) {
		t.Parallel()

		var a, b any
		a = &b
		b = &a
		assert.EqualError(
			t,
			reflect.Set(a, "Name", "Jane", false),
			`set (*interface {})."Name": unexpected pointer loop`,
		)
	})
	t.Run("unexported type of field", func(t *testing.T) {
		t.Parallel()

		p := person{}
		assert.NoError(t, reflect.Set(&p, "wallet", wallet{amount: 400}, false))
		assert.Equal(t, wallet{amount: 400}, p.wallet)
	})
	t.Run("convert []any to []type", func(t *testing.T) {
		t.Parallel()

		s := storage{}
		assert.NoError(
			t,
			reflect.Set(&s, "wallets", []any{wallet{100}, wallet{200}}, true),
		)
		assert.Equal(
			t,
			[]wallet{{100}, {200}},
			s.wallets,
		)
	})
	t.Run("convert int16 to uint", func(t *testing.T) {
		t.Parallel()

		var p struct {
			Age uint
		}
		err := reflect.Set(&p, "Age", int16(20), true)
		assert.NoError(t, err)
		assert.Equal(t, uint(20), p.Age)
	})
	t.Run("Given errors", func(t *testing.T) {
		t.Parallel()

		t.Run("Field does not exist", func(t *testing.T) {
			t.Parallel()

			p := person{}
			err := reflect.Set(&p, "FirstName", "Mary", false)
			assert.EqualError(t, err, `set (*reflect_test.person)."FirstName": field "FirstName" does not exist`)
		})
		t.Run("Invalid pointer dest", func(t *testing.T) {
			t.Parallel()

			p := 5
			err := reflect.Set(&p, "FirstName", "Mary", false)
			assert.EqualError(t, err, `set (*int)."FirstName": expected pointer to struct, *int given`)
		})
		t.Run("Invalid type of value", func(t *testing.T) {
			t.Parallel()

			t.Run("Convert", func(t *testing.T) {
				t.Parallel()

				p := person{}
				err := reflect.Set(&p, "Name", struct{}{}, true)
				assert.EqualError(t, err, `set (*reflect_test.person)."Name": cannot convert struct {} to string`)
			})
			t.Run("Do not convert", func(t *testing.T) {
				t.Parallel()

				p := person{}
				err := reflect.Set(&p, "Name", struct{}{}, false)
				assert.EqualError(t, err, `set (*reflect_test.person)."Name": value of type struct {} is not assignable to type string`) //nolint:lll
			})
		})
		t.Run("Invalid type of value (var p any = person{})", func(t *testing.T) {
			t.Parallel()

			var p any = person{}
			err := reflect.Set(&p, "Name", struct{}{}, true)
			assert.EqualError(t, err, `set (*interface {})."Name": cannot convert struct {} to string`)
		})
	})
	t.Run("Invalid struct #1", func(t *testing.T) {
		t.Parallel()

		err := reflect.Set(nil, "Name", "Jane", true)
		assert.EqualError(t, err, `set (<nil>)."Name": expected pointer to struct, <nil> given`)
	})
	t.Run("Invalid struct #2", func(t *testing.T) {
		t.Parallel()

		var a *struct {
			Name string
		}
		err := reflect.Set(&a, "Name", "Jane", true)
		assert.EqualError(t, err, `set (**struct { Name string })."Name": pointer to nil struct given`)
	})
	t.Run("Invalid struct #3", func(t *testing.T) {
		t.Parallel()

		var a *struct {
			Name string
		}
		var b any = a
		err := reflect.Set(&b, "Name", "Jane", true)
		assert.EqualError(t, err, `set (*interface {})."Name": expected pointer to struct, *interface {} given`)
	})
	t.Run("Invalid struct #4", func(t *testing.T) {
		t.Parallel()

		var a *struct {
			Name string
		}
		var b any = a
		err := reflect.Set(b, "Name", "Jane", true)
		assert.EqualError(t, err, `set (*struct { Name string })."Name": pointer to nil struct given`)
	})
}

type person struct {
	Name   string
	age    uint8
	wallet wallet
}

type wallet struct {
	amount uint
}

type storage struct {
	wallets []wallet
}

func TestIterateFields(t *testing.T) {
	t.Parallel()

	t.Run("Set", func(t *testing.T) {
		t.Parallel()

		scenarios := []struct {
			strct        any
			callback     reflect.FieldCallback
			convert      bool
			convertToPtr bool

			expected any
			error    string
		}{
			{
				strct: person{},
				callback: func(f stdReflect.StructField, value interface{}) (_ any, set bool) {
					if f.Name == "Name" {
						return "Jane", true
					}

					if f.Name == "age" {
						return uint(30), true
					}

					return nil, false
				},
				convert:      true,
				convertToPtr: false,
				expected: person{
					Name: "Jane",
					age:  30,
				},
			},
		}

		for i, s := range scenarios {
			s := s

			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				t.Parallel()

				strct := s.strct
				err := reflect.IterateFields(&strct, s.callback, s.convert, s.convertToPtr)

				if s.error != "" {
					require.EqualError(t, err, s.error)

					return
				}

				require.NoError(t, err)
				assert.Equal(t, s.expected, strct)
			})
		}
	})
}
