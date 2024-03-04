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
	"github.com/gontainer/reflectpro/getter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type character struct {
	name string
}

func (c *character) SetName(n string) {
	c.name = n
}

func (c character) GetName() string {
	return c.name
}

type Caller struct{}

func (*Caller) Call(fn func()) {
	fn()
}

func TestNewCallMethod_error(t *testing.T) {
	t.Parallel()

	t.Run("No pointer, expected pointer receiver", func(t *testing.T) {
		t.Parallel()

		c := character{}
		_, err := caller.NewCallMethod(c, "SetName", []any{"Thor"}, true)
		require.EqualError(
			t,
			err,
			`cannot call method (caller_test.character)."SetName": (caller_test.character)."SetName": invalid method`,
		)
		assert.Empty(t, c.name)
	})

	t.Run("Invalid method", func(t *testing.T) {
		t.Parallel()

		object := &struct{}{}
		_, err := caller.NewCallMethod(&object, "Do", nil, true)
		require.EqualError(t, err, `cannot call method (**struct {})."Do": (*struct {})."Do": invalid method`)
	})

	t.Run("Cannot convert args", func(t *testing.T) {
		t.Parallel()

		c := character{}
		_, err := caller.NewCallMethod(&c, "SetName", []any{struct{}{}}, true)
		require.EqualError(
			t,
			err,
			`cannot call method (*caller_test.character)."SetName": arg0: cannot convert struct {} to string`,
		)
	})

	t.Run("Not assignable args", func(t *testing.T) {
		t.Parallel()

		c := character{}
		_, err := caller.NewCallMethod(&c, "SetName", []any{struct{}{}}, false)
		require.EqualError(
			t,
			err,
			`cannot call method (*caller_test.character)."SetName": arg0: value of type struct {} is not assignable to type string`, //nolint:lll
		)
	})

	t.Run("Too many input arguments", func(t *testing.T) {
		t.Parallel()

		c := character{}
		_, err := caller.NewCallMethod(&c, "SetName", []any{"Bernhard", "Riemann"}, false)
		require.EqualError(
			t,
			err,
			`cannot call method (*caller_test.character)."SetName": too many input arguments`,
		)
	})

	t.Run("Not enough input arguments", func(t *testing.T) {
		t.Parallel()

		c := character{}
		_, err := caller.NewCallMethod(&c, "SetName", []any{}, false)
		require.EqualError(
			t,
			err,
			`cannot call method (*caller_test.character)."SetName": not enough input arguments`,
		)
	})
}

func TestNewCallMethod_okPointer(t *testing.T) {
	t.Parallel()

	t.Run("Pointer", func(t *testing.T) {
		t.Parallel()

		c := character{}
		_, err := caller.NewCallMethod(&c, "SetName", []any{"Thor"}, true)
		require.NoError(t, err)
		assert.Equal(t, c.name, "Thor")
	})

	t.Run("Pointer #2", func(t *testing.T) {
		t.Parallel()

		c := &character{}
		_, err := caller.NewCallMethod(c, "SetName", []any{"Thor"}, true)
		require.NoError(t, err)
		assert.Equal(t, c.name, "Thor")
	})

	t.Run("Pointer to pointer", func(t *testing.T) {
		t.Parallel()

		c := &character{}
		_, err := caller.NewCallMethod(&c, "SetName", []any{"Thor"}, true)
		require.NoError(t, err)
		assert.Equal(t, c.name, "Thor")
	})

	t.Run("Pointer to pointer to pointer", func(t *testing.T) {
		t.Parallel()

		tmp := &character{}
		c := &tmp

		_, err := caller.NewCallMethod(&c, "SetName", []any{"Thor"}, true)
		require.NoError(t, err)
		actual, _ := getter.Get(c, "name")
		assert.Equal(t, actual, "Thor")
	})

	t.Run("Pointer to any", func(t *testing.T) {
		t.Parallel()

		var c any = character{}

		_, err := caller.NewCallMethod(&c, "SetName", []any{"Thor"}, true)
		require.NoError(t, err)
		actual, _ := getter.Get(c, "name")
		assert.Equal(t, actual, "Thor")
	})

	t.Run("Pointer to any over pointer", func(t *testing.T) {
		t.Parallel()

		var c any = &character{}

		_, err := caller.NewCallMethod(&c, "SetName", []any{"Thor"}, true)
		require.NoError(t, err)
		actual, _ := getter.Get(c, "name")
		assert.Equal(t, actual, "Thor")
	})

	t.Run("Any over pointer", func(t *testing.T) {
		t.Parallel()

		var c any = &character{}

		_, err := caller.NewCallMethod(c, "SetName", []any{"Thor"}, true)
		require.NoError(t, err)
		actual, _ := getter.Get(c, "name")
		assert.Equal(t, actual, "Thor")
	})

	// Golang allows for executing methods over nil-receivers.
	t.Run("Nil pointer", func(t *testing.T) {
		t.Parallel()

		var (
			c        *Caller
			executed = false
			arg      = func() {
				executed = true
			}
		)

		_, err := caller.NewCallMethod(c, "Call", []any{arg}, false)
		require.NoError(t, err)
		assert.True(t, executed)
	})

	t.Run("Getter over value", func(t *testing.T) {
		t.Parallel()

		c := character{name: "Thor"}
		r, err := caller.NewCallMethod(c, "GetName", nil, false)
		require.NoError(t, err)
		assert.Equal(t, []any{"Thor"}, r)
	})

	t.Run("Getter over pointer", func(t *testing.T) {
		t.Parallel()

		c := character{name: "Thor"}
		r, err := caller.NewCallMethod(&c, "GetName", nil, false)
		require.NoError(t, err)
		assert.Equal(t, []any{"Thor"}, r)
	})
}
