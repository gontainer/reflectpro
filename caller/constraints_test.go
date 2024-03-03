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
	"testing"

	"github.com/gontainer/reflectpro/caller"
	"github.com/stretchr/testify/assert"
)

type book struct {
	title string
}

func (b *book) SetTitle(t string) {
	b.title = t
}

func (b *book) setTitle(t string) { //nolint:unused
	b.title = t
}

func (b book) WithTitle(t string) book {
	b.title = t

	return b
}

func TestConstraint(t *testing.T) {
	t.Parallel()

	const (
		harryPotterTitle = "Harry Potter"
	)

	var (
		harryPotter = book{title: harryPotterTitle}
		emptyBook   = book{}
	)

	// https://github.com/golang/go/wiki/MethodSets#interfaces

	// Method with a pointer receiver requires explicit definition of the pointer:
	// v := &book{}; CallMethod(v, ...
	// var v any = &book{}; CallMethod(v, ...
	// v := book{}; CallMethod(&v, ...
	//
	// Creating variable as a value will not work:
	// v := book{}; CallMethod(v, ...
	// var v interface{} = book{}; CallMethod(&v, ...
	t.Run("Call a method", func(t *testing.T) {
		t.Parallel()

		t.Run("Pointer receiver", func(t *testing.T) {
			t.Parallel()

			t.Run("Given errors", func(t *testing.T) {
				t.Parallel()

				t.Run("v := book{}; CallMethod(v, ...", func(t *testing.T) {
					t.Parallel()

					b := book{}
					r, err := caller.CallMethod(b, "SetTitle", []any{harryPotterTitle}, false)
					assert.EqualError(t, err, `cannot call method (caller_test.book)."SetTitle": (caller_test.book)."SetTitle": invalid method`)
					assert.Nil(t, r)
					assert.Zero(t, b)
				})
				t.Run("var v any = book{}; CallMethod(&v, ...", func(t *testing.T) {
					t.Parallel()

					//nolint
					// Since all the constraint don't exist anymore the following test does not make any sense
					// TODO remove this file
					//t.Run("CallMethod", func(t *testing.T) {
					//	t.Parallel()
					//
					//	var b any = book{}
					//	r, err := caller.CallMethod(&b, "SetTitle", []any{harryPotterTitle}, false)
					//	assert.EqualError(t, err, `cannot call method (*interface {})."SetTitle": (*interface {})."SetTitle": invalid method`)
					//	assert.Nil(t, r)
					//	assert.Equal(t, emptyBook, b)
					//})
					t.Run("ForceCallMethod", func(t *testing.T) {
						t.Parallel()

						var b any = book{}
						r, err := caller.ForceCallMethod(&b, "SetTitle", []any{harryPotterTitle}, false)
						assert.NoError(t, err)
						assert.Nil(t, r)
						assert.Equal(t, harryPotter, b)
					})
				})
			})
			t.Run("Given scenarios", func(t *testing.T) {
				t.Parallel()

				t.Run("v := book{}; CallMethod(&v, ...", func(t *testing.T) {
					t.Parallel()

					b := book{}
					r, err := caller.CallMethod(&b, "SetTitle", []any{harryPotterTitle}, false)
					assert.NoError(t, err)
					assert.Len(t, r, 0)
					assert.Equal(t, harryPotter, b)
				})
				t.Run("v := &book{}; CallMethod(&v, ...", func(t *testing.T) {
					t.Parallel()

					b := &book{}
					r, err := caller.CallMethod(&b, "SetTitle", []any{harryPotterTitle}, false)
					assert.NoError(t, err)
					assert.Len(t, r, 0)
					assert.Equal(t, &harryPotter, b)
				})
				t.Run("v := &book{}; CallMethod(v, ...", func(t *testing.T) {
					t.Parallel()

					b := &book{}
					r, err := caller.CallMethod(b, "SetTitle", []any{harryPotterTitle}, false)
					assert.NoError(t, err)
					assert.Len(t, r, 0)
					assert.Equal(t, &harryPotter, b)
				})
				t.Run("var v any = &book{}; CallMethod(v, ...", func(t *testing.T) {
					t.Parallel()

					var b any = &book{}
					r, err := caller.CallMethod(b, "SetTitle", []any{harryPotterTitle}, false)
					assert.NoError(t, err)
					assert.Len(t, r, 0)
					assert.Equal(t, &harryPotter, b)
				})
				t.Run("var v any = &book{}; CallMethod(&v, ...", func(t *testing.T) {
					t.Parallel()

					var b any = &book{}
					r, err := caller.CallMethod(&b, "SetTitle", []any{harryPotterTitle}, false)
					assert.NoError(t, err)
					assert.Len(t, r, 0)
					assert.Equal(t, &harryPotter, b)
				})
				t.Run("var v interface{ SetTitle(string) } = &book{}; CallMethod(v, ...", func(t *testing.T) {
					t.Parallel()

					var b interface{ SetTitle(string) } = &book{} //nolint:inamedparam
					r, err := caller.CallMethod(b, "SetTitle", []any{harryPotterTitle}, false)
					assert.NoError(t, err)
					assert.Len(t, r, 0)
					assert.Equal(t, &harryPotter, b)
				})
			})
		})
		// Methods with a value receiver do not have any constraints
		t.Run("Value receiver", func(t *testing.T) {
			t.Parallel()

			t.Run("b := book{}", func(t *testing.T) {
				t.Parallel()

				b := book{}
				r, err := caller.CallWither(b, "WithTitle", []any{harryPotterTitle}, false)
				assert.NoError(t, err)
				assert.Equal(t, harryPotter, r)
				assert.Zero(t, b)
			})
			t.Run("b := &book{}", func(t *testing.T) {
				t.Parallel()

				b := &book{}
				r, err := caller.CallWither(b, "WithTitle", []any{harryPotterTitle}, false)
				assert.NoError(t, err)
				assert.Equal(t, harryPotter, r)
				assert.Equal(t, &emptyBook, b)
			})
			t.Run("var b any = book{}", func(t *testing.T) {
				t.Parallel()

				var b any = book{}
				r, err := caller.CallWither(b, "WithTitle", []any{harryPotterTitle}, false)
				assert.NoError(t, err)
				assert.Equal(t, harryPotter, r)
				assert.Equal(t, emptyBook, b)
			})
			t.Run("var b any = &book{}", func(t *testing.T) {
				t.Parallel()

				var b any = &book{}
				r, err := caller.CallWither(b, "WithTitle", []any{harryPotterTitle}, false)
				assert.NoError(t, err)
				assert.Equal(t, harryPotter, r)
				assert.Equal(t, &emptyBook, b)
			})
		})
		t.Run("Unexported method", func(t *testing.T) {
			t.Parallel()

			t.Run("CallMethod", func(t *testing.T) {
				t.Parallel()

				b := book{}
				_, err := caller.CallMethod(&b, "setTitle", []any{harryPotter}, false)
				assert.EqualError(t, err, `cannot call method (*caller_test.book)."setTitle": (*caller_test.book)."setTitle": invalid method`)
			})
			t.Run("ForceCallMethod", func(t *testing.T) {
				t.Parallel()

				b := book{}
				_, err := caller.ForceCallMethod(&b, "setTitle", []any{harryPotter}, false)
				assert.EqualError(t, err, `cannot call method (*caller_test.book)."setTitle": (*caller_test.book)."setTitle": invalid method`)
			})
		})
	})
}
