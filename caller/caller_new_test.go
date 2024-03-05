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
	"database/sql"
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

type callerStrct struct{}

func (*callerStrct) Call(fn func()) {
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

	t.Run("Invalid method #2", func(t *testing.T) {
		t.Parallel()

		object := struct{}{}
		_, err := caller.NewCallMethod(&object, "Do", nil, true)
		require.EqualError(t, err, `cannot call method (*struct {})."Do": (*struct {})."Do": invalid method`)
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
			`cannot call method (*caller_test.character)."SetName": `+
				`arg0: value of type struct {} is not assignable to type string`,
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

	t.Run("Pointer to any(nil)", func(t *testing.T) {
		t.Parallel()

		var c any
		_, err := caller.NewCallMethod(&c, "Do", nil, true)
		require.EqualError(t, err, `cannot call method (*interface {})."Do": invalid object`)
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

	t.Run("Pointer to pointer to pointer to pointer", func(t *testing.T) {
		t.Parallel()

		tmp := &character{}
		tmp2 := &tmp
		c := &tmp2

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
			c        *callerStrct
			executed = false

			callee = func() {
				executed = true
			}
		)

		_, err := caller.NewCallMethod(c, "Call", []any{callee}, false)
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

type provider struct {
	fn func(...any) (any, error)
}

func (p provider) Provide(args ...any) (any, error) { //nolint:ireturn
	return p.fn(args...)
}

func TestNewCallProviderMethod_ok(t *testing.T) {
	t.Parallel()

	t.Run("2+2=4", func(t *testing.T) {
		t.Parallel()

		provider := provider{
			fn: func(args ...any) (any, error) {
				require.Equal(t, []any{2, 2}, args)

				return 4, nil
			},
		}

		result, executed, err := caller.NewCallProviderMethod(provider, "Provide", []any{2, 2}, true)
		require.NoError(t, err)
		assert.True(t, executed)
		assert.Equal(t, 4, result)
	})
}

func TestNewCallProviderMethod_error(t *testing.T) {
	t.Parallel()

	t.Run("Not assignable args", func(t *testing.T) {
		t.Parallel()

		db := &sql.DB{}
		result, executed, err := caller.NewCallProviderMethod(
			db,
			"BeginTx",
			[]any{"start", nil},
			false,
		)
		require.EqualError(
			t,
			err,
			`cannot call provider (*sql.DB)."BeginTx": cannot call method (*sql.DB)."BeginTx": `+
				`arg0: value of type string is not assignable to type context.Context`,
		)
		require.False(t, executed)
		assert.Nil(t, result)
	})

	t.Run("Cannot convert args", func(t *testing.T) {
		t.Parallel()

		db := &sql.DB{}
		result, executed, err := caller.NewCallProviderMethod(
			db,
			"BeginTx",
			[]any{"start", nil},
			true,
		)
		require.EqualError(
			t,
			err,
			`cannot call provider (*sql.DB)."BeginTx": cannot call method (*sql.DB)."BeginTx": `+
				`arg0: cannot convert string to context.Context`,
		)
		require.False(t, executed)
		assert.Nil(t, result)
	})
}
