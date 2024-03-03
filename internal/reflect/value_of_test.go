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
	"math"
	"reflect"
	"testing"

	intReflect "github.com/gontainer/reflectpro/internal/reflect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type myMap map[string]any

func (m myMap) Foo() {
}

type myMap2 map[string]interface{ Foo() }

//nolint:maintidx
func TestValueOf(t *testing.T) {
	t.Parallel()

	t.Run("Do not convert", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			var x int = 5 //nolint
			v, err := intReflect.ValueOf(x, reflect.TypeOf(x), false)
			require.NoError(t, err)
			assert.Equal(t, x, v.Interface())
		})
		t.Run("Error", func(t *testing.T) {
			t.Parallel()

			var x int = 5 //nolint
			_, err := intReflect.ValueOf(x, reflect.TypeOf(uint(0)), false)
			assert.EqualError(t, err, "value of type int is not assignable to type uint")
		})
	})
	t.Run("Empty maps & slices", func(t *testing.T) {
		t.Run("Nil slice", func(t *testing.T) {
			t.Run("Do not convert", func(t *testing.T) {
				t.Parallel()

				var (
					from []int
					to   []int
				)

				r, err := intReflect.ValueOf(from, reflect.TypeOf(to), false)
				require.NoError(t, err)
				assert.Nil(t, r.Interface())
			})
			t.Run("Convert", func(t *testing.T) {
				var (
					from []int
					to   []uint
				)

				r, err := intReflect.ValueOf(from, reflect.TypeOf(to), true)
				require.NoError(t, err)
				assert.Nil(t, r.Interface())
			})
		})
		t.Run("Non-nil slice", func(t *testing.T) {
			var (
				from = make([]int, 0)
				to   []int
			)

			r, err := intReflect.ValueOf(from, reflect.TypeOf(to), true)
			require.NoError(t, err)
			assert.NotNil(t, r.Interface())
			assert.Len(t, r.Interface(), 0)
		})
		t.Run("Nil map", func(t *testing.T) {
			t.Run("Do not convert", func(t *testing.T) {
				var (
					from map[string]any
					to   map[string]any
				)

				r, err := intReflect.ValueOf(from, reflect.TypeOf(to), false)
				require.NoError(t, err)
				assert.Nil(t, r.Interface())
			})
			t.Run("Convert", func(t *testing.T) {
				var (
					from map[int]any
					to   map[uint]any
				)

				r, err := intReflect.ValueOf(from, reflect.TypeOf(to), true)
				require.NoError(t, err)
				assert.Nil(t, r.Interface())
			})
		})
		t.Run("Non-nil map", func(t *testing.T) {
			var (
				from = make(map[string]any)
				to   map[string]any
			)

			r, err := intReflect.ValueOf(from, reflect.TypeOf(to), true)
			require.NoError(t, err)
			assert.NotNil(t, r.Interface())
		})
	})

	t.Run("Nil", func(t *testing.T) {
		t.Parallel()

		var to map[string]string
		r, err := intReflect.ValueOf(nil, reflect.TypeOf(to), false)
		require.NoError(t, err)
		assert.IsType(t, to, r.Interface())
	})

	t.Run("Maps", func(t *testing.T) {
		t.Parallel()

		t.Run("OK", func(t *testing.T) {
			t.Parallel()

			t.Run("Convert keys and values", func(t *testing.T) {
				type myStr string
				from := map[myStr]int32{
					"2^8":  256,
					"2^9":  512,
					"2^10": 1024,
				}
				var to map[string]int64
				v, err := intReflect.ValueOf(from, reflect.TypeOf(to), true)
				require.NoError(t, err)
				assert.Equal(
					t,
					map[string]int64{
						"2^8":  256,
						"2^9":  512,
						"2^10": 1024,
					},
					v.Interface(),
				)
			})
		})
		t.Run("Errors", func(t *testing.T) {
			scenarios := []struct {
				input any
				to    reflect.Type
				error string
			}{
				{
					input: (map[struct{}]any)(nil),
					to:    reflect.TypeOf((map[bool]any)(nil)),
					error: "cannot convert map[struct {}]interface {} to map[bool]interface {}: non convertible keys: cannot convert struct {} to bool", //nolint:lll
				},
				{
					input: (map[string]int32)(nil),
					to:    reflect.TypeOf((map[string]struct{})(nil)),
					error: "cannot convert map[string]int32 to map[string]struct {}: non convertible values: cannot convert int32 to struct {}", //nolint:lll
				},
				{
					input: map[string]any{"pi": "3.14"},
					to:    reflect.TypeOf((map[string]float64)(nil)),
					error: "cannot convert map[string]interface {} to map[string]float64: map value: cannot convert string to float64",
				},
				{
					input: map[any]any{true: "true"},
					to:    reflect.TypeOf((map[uint]any)(nil)),
					error: "cannot convert map[interface {}]interface {} to map[uint]interface {}: map key: cannot convert bool to uint", //nolint:lll
				},
			}

			for i, s := range scenarios {
				s := s
				t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
					_, err := intReflect.ValueOf(s.input, s.to, true)
					assert.EqualError(t, err, s.error)
				})
			}
		})
	})
	t.Run("Recursion", func(t *testing.T) {
		t.Run("Slice refers to itself", func(t *testing.T) {
			t.Parallel()

			x := make([]any, 1)
			x[0] = x
			y, err := intReflect.ValueOf(x, reflect.TypeOf(x), true)
			require.NoError(t, err)
			assert.Equal(t, x, y.Interface())
		})
		t.Run("Map refers to itself", func(t *testing.T) {
			t.Run("#1", func(t *testing.T) {
				x := make(map[string]any)
				x["self"] = x
				y, err := intReflect.ValueOf(x, reflect.TypeOf(x), true)
				require.NoError(t, err)
				assert.Equal(t, x, y.Interface())
			})
			t.Run("#2", func(t *testing.T) {
				x := make(map[string]any)
				x["child"] = map[string]any{
					"grandchild": map[string]any{
						"grandgrandchild": map[string]any{
							"x": x,
						},
					},
				}
				y, err := intReflect.ValueOf(x, reflect.TypeOf(x), true)
				require.NoError(t, err)
				assert.Equal(t, x, y.Interface())
			})
			t.Run("#3", func(t *testing.T) {
				x := make(myMap)
				x["self"] = x
				y, err := intReflect.ValueOf(x, reflect.TypeOf((map[string]any)(nil)), true)
				require.NoError(t, err)
				assert.Equal(t, (map[string]any)(x), y.Interface())
			})
			t.Run("#4", func(t *testing.T) {
				x := make(map[string]any)
				x["self"] = x
				y, err := intReflect.ValueOf(x, reflect.TypeOf((myMap)(nil)), true)
				require.NoError(t, err)
				assert.Equal(t, (myMap)(x), y.Interface())
			})
			t.Run("#5", func(t *testing.T) {
				x := make(map[string]any)
				x["self"] = x
				_, err := intReflect.ValueOf(x, reflect.TypeOf((myMap2)(nil)), true)
				assert.EqualError(t, err, "cannot convert map[string]interface {} to reflect_test.myMap2: map value: cannot convert map[string]interface {} to interface { Foo() }") //nolint:lll
			})
			t.Run("#6", func(t *testing.T) {
				x := make(map[float32]any)
				x[0] = x
				v, err := intReflect.ValueOf(x, reflect.TypeOf((map[float64]any)(nil)), true)
				require.NoError(t, err)

				val, ok := v.Interface().(map[float64]any)
				require.True(t, ok)
				assert.Equal(t, x, val[0])
			})
		})
	})
	t.Run("Convert parameters", func(t *testing.T) {
		float64Val := float64(5)

		scenarios := map[string]struct {
			input  any
			output any
			error  string
		}{
			`nil to interface`: {
				input:  nil,
				output: (*interface{ Bar() })(nil),
			},
			`[3]int to [2]int`: {
				input:  [3]int{1, 2, 3},
				output: [2]int{1, 2},
			},
			`[]int (len == 3) to [2]int`: {
				input:  []int{1, 2, 3},
				output: [2]int{1, 2},
			},
			`[2]int to [3]int`: {
				input:  [2]int{1, 2},
				output: [3]int{1, 2, 0},
			},
			`[]int (len == 2) to [3]int`: {
				input:  []int{1, 2},
				output: [3]int{1, 2, 0},
			},
			`[]any{[]int{1, 2, 3}} to [][2]int{}`: {
				input:  []any{[]int{1, 2, 3}},
				output: [][2]int{{1, 2}},
			},
			`[]any to [0]int`: {
				input:  []any{},
				output: [0]int{},
			},
			`[0]int to []any`: {
				input:  [0]int{},
				output: []any{},
			},
			`[][3]int to [][2]int`: {
				input:  [][3]int{{3, 4, 5}},
				output: [][2]int{{3, 4}},
			},
			`[][3]int to [][3]int`: {
				input:  [][3]int{{5, 5, 5}, {6, 6, 6}},
				output: [][3]int{{5, 5, 5}, {6, 6, 6}},
			},
			`[][3]int to [][3]uint`: {
				input:  [][3]int{{1, 2, 3}},
				output: [][3]uint{{1, 2, 3}},
			},
			`[][3]int to [][3]any`: {
				input:  [][3]int{{2, 2, 2}},
				output: [][3]any{{2, 2, 2}},
			},
			`[][3]any to [][3]int`: {
				input:  [][3]any{{3, 5, 7}},
				output: [][3]int{{3, 5, 7}},
			},
			`[]any{[2]int{}} to [][3]int`: {
				input:  []any{[2]int{5, 5}},
				output: [][3]int{{5, 5, 0}},
			},
			`[][]any to [][]int`: {
				input:  [][]any{{1, 2, 3}},
				output: [][]int{{1, 2, 3}},
			},
			`[][]any to [][]int (invalid)`: {
				input:  [][]any{{1, false, 3}},
				output: [][]int{{1, 2, 3}},
				error:  "cannot convert [][]interface {} to [][]int: #0: cannot convert []interface {} to []int: #1: cannot convert bool to int", //nolint:lll
			},
			`[][]int to [][]any`: {
				input:  [][]int{{1, 2, 3}},
				output: [][]any{{1, 2, 3}},
			},
			`[][]uint to [][]int`: {
				input:  [][]uint{{1, 2, 3}},
				output: [][]int{{1, 2, 3}},
			},
			`[]any to []int`: {
				input:  []any{1, 2, 3},
				output: []int{1, 2, 3},
			},
			`[]any to []int (invalid #1)`: {
				input:  []any{1, 2, nil},
				output: []int{},
				error:  "cannot convert []interface {} to []int: #2: cannot convert <nil> to int",
			},
			`[]any to []int (invalid #2)`: {
				input:  []any{1, 2, 3, struct{}{}},
				output: []int{},
				error:  "cannot convert []interface {} to []int: #3: cannot convert struct {} to int",
			},
			`[]any to []*int`: {
				input:  []any{nil, nil},
				output: []*int{nil, nil},
			},
			`[]int to []any`: {
				input:  []int{1, 2, 3},
				output: []any{1, 2, 3},
			},
			`[]any to []any`: {
				input:  []any{1, 2, 3},
				output: []any{1, 2, 3},
			},
			`[]int8 to []int`: {
				input:  []int8{1, 2, 3},
				output: []int{1, 2, 3},
			},
			`[]int to []int8`: {
				input:  []int{1, 2, 256},
				output: []int8{1, 2, 0},
			},
			`[]struct{}{} to []type`: {
				input:  []struct{}{},
				output: []int{},
				error:  `cannot convert []struct {} to []int: cannot convert struct {} to int`,
			},
			`float64 to int`: {
				input:  float64(math.Pi),
				output: 3,
			},
			`nil to int`: {
				input:  nil,
				output: 0,
				error:  "cannot convert <nil> to int",
			},
			`nil to *int`: {
				input:  nil,
				output: (*int)(nil),
			},
			`*float64 to *int`: {
				input:  &float64Val,
				output: (*int)(nil),
				error:  "cannot convert *float64 to *int",
			},
			`*float64 to *float32`: {
				input:  &float64Val,
				output: (*float32)(nil),
				error:  "cannot convert *float64 to *float32",
			},
			`*float64 to *float64`: {
				input:  &float64Val,
				output: &float64Val,
			},
			`int to float64`: {
				input:  int(5),
				output: float64(5),
			},
			`string to []byte`: {
				input:  `hello`,
				output: []byte(`hello`),
			},
			`[]byte to string`: {
				input:  []byte(`hello`),
				output: `hello`,
			},
			`string to int`: { // cannot convert string to int
				input:  `5`,
				output: int(5),
				error:  "cannot convert string to int",
			},
			`int to string`: { // but reverse conversion is possible, isn't worth to unify this behavior?
				input:  5,
				output: "\x05",
			},
			`nil to *any`: {
				input:  nil,
				output: (*any)(nil),
			},
		}

		for n, s := range scenarios {
			s := s
			t.Run(n, func(t *testing.T) {
				t.Parallel()

				v, err := intReflect.ValueOf(s.input, reflect.TypeOf(s.output), true)
				if s.error != "" {
					assert.EqualError(t, err, s.error, v)

					return
				}

				require.NoError(t, err)
				assert.Equal(t, s.output, v.Interface())
			})
		}
	})
}
