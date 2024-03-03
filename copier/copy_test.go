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
	"reflect"
	"testing"

	"github.com/gontainer/reflectpro/copier"
	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	t.Parallel()

	t.Run("Convert value", func(t *testing.T) {
		t.Parallel()

		t.Run("Do not convert", func(t *testing.T) {
			t.Parallel()

			var (
				from int = 5
				to   uint
			)
			err := copier.Copy(from, &to, false)
			assert.Empty(t, to)
			assert.EqualError(t, err, "value of type int is not assignable to type uint")
		})
		t.Run("Convert", func(t *testing.T) {
			t.Parallel()

			var (
				from int = 5
				to   uint
			)
			err := copier.Copy(from, &to, true)
			assert.Equal(t, uint(5), to)
			assert.NoError(t, err)
		})
	})
	t.Run("Convert pointer", func(t *testing.T) {
		t.Parallel()

		t.Run("Do not convert", func(t *testing.T) {
			t.Parallel()

			var (
				from *int
				to   *uint
			)
			err := copier.Copy(from, &to, false)
			assert.Empty(t, to)
			assert.EqualError(t, err, "value of type *int is not assignable to type *uint")
		})
		t.Run("Convert", func(t *testing.T) {
			t.Parallel()

			var (
				from *int
				to   *uint
			)
			err := copier.Copy(from, &to, true)
			assert.Empty(t, to)
			assert.EqualError(t, err, "cannot convert *int to *uint")
		})
	})
	t.Run("Non-empty interface", func(t *testing.T) {
		t.Parallel()

		t.Run("Do not convert", func(t *testing.T) {
			t.Parallel()

			var (
				from interface{ Foo() }
				to   interface{ Bar() } // even tho interfaces differ, it does not return an error
			)
			err := copier.Copy(from, &to, false)
			assert.Empty(t, to)
			assert.NoError(t, err)
		})
		t.Run("Convert", func(t *testing.T) {
			t.Parallel()

			var (
				from interface{ Foo() }
				to   interface{ Bar() } // even tho interfaces differ, it does not return an error
			)
			err := copier.Copy(from, &to, true)
			assert.Empty(t, to)
			assert.NoError(t, err)
		})
	})
	t.Run("var from any", func(t *testing.T) {
		t.Parallel()

		var (
			from any = car{age: 5}
			to   car
		)

		assert.NoError(t, copier.Copy(from, &to, false))
		assert.Equal(t, 5, to.age)
	})
	t.Run("var to any", func(t *testing.T) {
		six := 6
		scenarios := []any{
			5,
			3.14,
			struct{}{},
			nil,
			&six,
			car{age: 10},
			&car{age: 10},
			(*car)(nil),
		}

		for id, tmp := range scenarios {
			d := tmp
			t.Run(fmt.Sprintf("%d: `%T`", id, d), func(t *testing.T) {
				t.Parallel()

				var to any
				assert.NoError(t, copier.Copy(d, &to, false))
				assert.Equal(t, d, to)
				if reflect.ValueOf(d).Kind() == reflect.Ptr {
					assert.Same(t, d, to)
				}
			})
		}
	})
	t.Run("Given errors", func(t *testing.T) {
		t.Parallel()

		t.Run("non-pointer value", func(t *testing.T) {
			t.Parallel()

			const msg = "expected ptr, int given"
			assert.EqualError(
				t,
				copier.Copy(5, 5, false),
				msg,
			)
		})
	})
	t.Run("Copy to nil", func(t *testing.T) {
		t.Parallel()

		assert.EqualError(
			t,
			copier.Copy(5, nil, false),
			"expected ptr, <nil> given",
		)
	})
	t.Run("Convert", func(t *testing.T) {
		t.Parallel()

		t.Run("[]int to []any", func(t *testing.T) {
			t.Parallel()

			var (
				from = []int{1, 2, 3}
				to   []any
			)

			err := copier.Copy(from, &to, true)
			assert.NoError(t, err)
			assert.Equal(t, []any{1, 2, 3}, to)
		})
		t.Run("[]any to []int", func(t *testing.T) {
			t.Parallel()

			var (
				from = []any{1, 2, 3}
				to   []int
			)

			err := copier.Copy(from, &to, true)
			assert.NoError(t, err)
			assert.Equal(t, []int{1, 2, 3}, to)
		})
		t.Run("[]int to [N]int", func(t *testing.T) {
			t.Parallel()

			var (
				from = []int{1, 2, 3}
				to   [3]int
			)

			err := copier.Copy(from, &to, true)
			assert.NoError(t, err)
			assert.Equal(t, [3]int{1, 2, 3}, to)
		})
		t.Run("[]int to [N]int #2", func(t *testing.T) {
			t.Parallel()

			var (
				from = []int{1, 2, 3}
				to   [2]int
			)

			err := copier.Copy(from, &to, true)
			assert.NoError(t, err)
			assert.Equal(t, [2]int{1, 2}, to)
		})
		t.Run("[N]int to [N-1]int", func(t *testing.T) {
			t.Parallel()

			var (
				from = [3]int{1, 2, 3}
				to   [2]int
			)

			err := copier.Copy(from, &to, true)
			assert.NoError(t, err)
			assert.Equal(t, [2]int{1, 2}, to)
		})
		t.Run("[N]int to [N+1]int", func(t *testing.T) {
			t.Parallel()

			var (
				from = [3]int{1, 2, 3}
				to   [4]int
			)

			err := copier.Copy(from, &to, true)
			assert.NoError(t, err)
			assert.Equal(t, [4]int{1, 2, 3, 0}, to)
		})
		t.Run("[N]any to [N+1]any", func(t *testing.T) {
			t.Parallel()

			var (
				from = [3]any{6, 7, 8}
				to   [4]any
			)

			err := copier.Copy(from, &to, true)
			assert.NoError(t, err)
			assert.Equal(t, [4]any{6, 7, 8, nil}, to)
		})
		t.Run("[N]int to []int", func(t *testing.T) {
			t.Parallel()

			var (
				from = [3]int{1, 2, 3}
				to   []int
			)

			err := copier.Copy(from, &to, true)
			assert.NoError(t, err)
			assert.Equal(t, []int{1, 2, 3}, to)
		})
		t.Run("[N]int to [N]uint", func(t *testing.T) {
			t.Parallel()

			var (
				from = [3]int{1, 2, 3}
				to   [3]uint
			)

			err := copier.Copy(from, &to, true)
			assert.NoError(t, err)
			assert.Equal(t, [3]uint{1, 2, 3}, to)
		})
	})
}

type car struct {
	age int
}
