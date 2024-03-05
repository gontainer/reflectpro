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
	"errors"
	"fmt"
	"strings"
	"testing"

	errAssert "github.com/gontainer/grouperror/assert"
	"github.com/gontainer/reflectpro/caller"
	"github.com/gontainer/reflectpro/getter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCall(t *testing.T) {
	t.Parallel()

	t.Run("Given method", func(t *testing.T) {
		t.Parallel()

		p := person{}
		assert.Equal(t, "", p.name)
		_, err := caller.Call(p.setName, []any{"Mary"}, false)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "Mary", p.name)
	})

	t.Run("Given invalid functions", func(t *testing.T) {
		t.Parallel()

		scenarios := []struct {
			fn any
		}{
			{fn: 5},
			{fn: false},
			{fn: (*error)(nil)},
			{fn: struct{}{}},
		}
		const expectedRegexp = "\\Acannot call .*: expected func, .* given\\z"
		for i, tmp := range scenarios {
			s := tmp
			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				t.Parallel()

				_, err := caller.Call(s.fn, nil, false)
				assert.Error(t, err)
				assert.Regexp(t, expectedRegexp, err)
			})
		}
	})

	t.Run("Nil-func", func(t *testing.T) {
		t.Parallel()

		r, err := caller.Call(nil, nil, true)
		assert.EqualError(t, err, "cannot call <nil>: invalid func: <nil>")
		assert.Nil(t, r)
	})

	t.Run("Given invalid argument", func(t *testing.T) {
		t.Parallel()

		const msg = "cannot call func([]int): arg0: cannot convert struct {} to []int"
		callee := func([]int) {}
		params := []any{
			struct{}{},
		}

		_, err := caller.Call(callee, params, true)
		assert.EqualError(t, err, msg)
	})

	t.Run("Given invalid arguments", func(t *testing.T) {
		t.Parallel()

		callee := func([]int, *int) {}
		params := []any{
			struct{}{},
			"*int",
		}

		_, err := caller.Call(callee, params, true)

		expected := []string{
			"cannot call func([]int, *int): arg0: cannot convert struct {} to []int",
			"cannot call func([]int, *int): arg1: cannot convert string to *int",
		}
		errAssert.EqualErrorGroup(t, err, expected)
	})

	t.Run("Given too many arguments", func(t *testing.T) {
		t.Parallel()

		const msg = "too many input arguments"
		scenarios := []struct {
			fn   any
			args []any
		}{
			{
				fn:   strings.Join,
				args: []any{"1", "2", "3"},
			},
			{
				fn:   func() {},
				args: []any{1},
			},
		}
		for i, tmp := range scenarios {
			s := tmp
			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				t.Parallel()
				_, err := caller.Call(s.fn, s.args, true)
				assert.ErrorContains(t, err, msg)
			})
		}
	})

	t.Run("Given too few arguments", func(t *testing.T) {
		const msg = "not enough input arguments"
		scenarios := []struct {
			fn   any
			args []any
		}{
			{
				fn:   strings.Join,
				args: []any{},
			},
			{
				fn:   func(a int) {},
				args: []any{},
			},
		}
		for i, tmp := range scenarios {
			s := tmp
			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				t.Parallel()

				_, err := caller.Call(s.fn, s.args, true)
				assert.ErrorContains(t, err, msg)
			})
		}
	})

	t.Run("Given scenarios", func(t *testing.T) {
		scenarios := []struct {
			fn       any
			args     []any
			expected []any
		}{
			{
				fn: func(a, b int) int {
					return a + b
				},
				args:     []any{uint(1), uint(2)},
				expected: []any{int(3)},
			},
			{
				fn: func(a, b uint) uint {
					return a + b
				},
				args:     []any{int(7), int(3)},
				expected: []any{uint(10)},
			},
			{
				fn: func(vals ...uint) (result uint) {
					for _, v := range vals {
						result += v
					}

					return
				},
				args:     []any{int(1), int(2), int(3)},
				expected: []any{uint(6)},
			},
		}
		for i, tmp := range scenarios {
			s := tmp
			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				t.Parallel()

				r, err := caller.Call(s.fn, s.args, true)
				assert.NoError(t, err)
				assert.Equal(t, s.expected, r)
			})
		}
	})

	t.Run("Convert parameters", func(t *testing.T) {
		scenarios := map[string]struct {
			fn     any
			input  any
			output any
			error  string
		}{
			"[]any to []type": {
				fn: func(v []int) []int {
					return v
				},
				input:  []any{1, 2, 3},
				output: []int{1, 2, 3},
			},
			"[]struct{}{} to []type": {
				fn:    func([]int) {},
				input: []struct{}{},
				error: `cannot call func([]int): arg0: cannot convert []struct {} to []int: cannot convert struct {} to int`,
			},
			"nil to any": {
				fn: func(v any) any {
					return v
				},
				input:  nil,
				output: (any)(nil),
			},
		}

		for n, tmp := range scenarios {
			s := tmp
			t.Run(n, func(t *testing.T) {
				t.Parallel()

				r, err := caller.Call(s.fn, []any{s.input}, true)
				if s.error != "" {
					assert.EqualError(t, err, s.error)
					assert.Nil(t, r)

					return
				}

				assert.NoError(t, err)
				assert.Equal(t, r[0], s.output)
			})
		}
	})
}

func TestCallProvider(t *testing.T) {
	t.Parallel()

	t.Run("Given scenarios", func(t *testing.T) {
		t.Parallel()

		scenarios := []struct {
			provider any
			params   []any
			expected any
		}{
			{
				provider: func() any {
					return nil
				},
				expected: nil,
			},
			{
				provider: func(vals ...int) (int, error) {
					result := 0
					for _, v := range vals {
						result += v
					}

					return result, nil
				},
				params:   []any{10, 100, 200},
				expected: 310,
			},
		}

		for i, tmp := range scenarios {
			s := tmp
			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				t.Parallel()

				r, executed, err := caller.CallProvider(s.provider, s.params, false)
				require.NoError(t, err)
				assert.True(t, executed)
				assert.Equal(t, s.expected, r)
			})
		}
	})

	t.Run("Provider error", func(t *testing.T) {
		t.Parallel()

		type myError struct {
			error
		}
		p := func() (any, error) {
			return nil, &myError{errors.New("my error")}
		}
		_, executed, err := caller.CallProvider(p, nil, false)
		require.EqualError(t, err, "provider returned error: my error")
		assert.True(t, executed)

		var providerErr *caller.ProviderError
		require.True(t, errors.As(err, &providerErr))
		assert.EqualError(t, providerErr, "my error")

		var myErr *myError
		require.True(t, errors.As(err, &myErr))
		assert.EqualError(t, myErr, "my error")
	})

	t.Run("Given errors", func(t *testing.T) {
		t.Parallel()

		scenarios := []struct {
			provider any
			params   []any
			executed bool
			err      string
		}{
			{
				provider: func() {},
				executed: false,
				err:      "cannot call provider func(): provider must return 1 or 2 values, given function returns 0 values",
			},
			{
				provider: func() (any, any, any) {
					return nil, nil, nil
				},
				executed: false,
				err:      "cannot call provider func() (interface {}, interface {}, interface {}): provider must return 1 or 2 values, given function returns 3 values",
			},
			{
				provider: func() (any, any) {
					return nil, nil
				},
				executed: false,
				err:      "cannot call provider func() (interface {}, interface {}): second value returned by provider must implement error interface, interface {} given",
			},
			{
				provider: func() (any, int) {
					return nil, 0
				},
				executed: false,
				err:      "cannot call provider func() (interface {}, int): second value returned by provider must implement error interface, int given",
			},
			{
				provider: func() (any, Person) {
					return nil, Person{}
				},
				executed: false,
				err:      "cannot call provider func() (interface {}, caller_test.Person): second value returned by provider must implement error interface, caller_test.Person given",
			},
			{
				provider: func() (any, error) {
					return nil, errors.New("test error")
				},
				executed: true,
				err:      "provider returned error: test error",
			},
			{
				provider: func() any {
					return nil
				},
				params:   []any{1, 2, 3},
				executed: false,
				err:      "cannot call provider func() interface {}: too many input arguments",
			},
		}

		for i, tmp := range scenarios {
			s := tmp
			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				t.Parallel()

				r, executed, err := caller.CallProvider(s.provider, s.params, false)
				require.EqualError(t, err, s.err)
				assert.Equal(t, s.executed, executed)
				assert.Nil(t, r)
			})
		}
	})

	t.Run("Given invalid provider", func(t *testing.T) {
		t.Parallel()

		_, executed, err := caller.CallProvider(5, nil, false)
		require.EqualError(t, err, "cannot call provider int: expected func, int given")
		assert.False(t, executed)
	})

	t.Run("Given provider panics", func(t *testing.T) {
		t.Parallel()

		defer func() {
			assert.Equal(t, "panic!", recover())
		}()

		_, _, _ = caller.CallProvider(
			func() any {
				panic("panic!")
			},
			nil,
			false,
		)
	})
}

type mockProvider struct {
	fn          func() any
	fnWithError func() (any, error)
}

func (m *mockProvider) Provider() any { //nolint:ireturn
	return m.fn()
}

func (m *mockProvider) ProviderWithError() (any, error) { //nolint:ireturn
	return m.fnWithError()
}

func (m *mockProvider) NotProvider() (any, any) { //nolint:ireturn
	return nil, nil
}

func TestCallProviderMethod(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		t.Run("#1", func(t *testing.T) {
			t.Parallel()

			p := mockProvider{
				fn: func() any {
					return "value #1"
				},
			}
			r, executed, err := caller.CallProviderMethod(&p, "Provider", nil, false)
			assert.NoError(t, err)
			assert.True(t, executed)
			assert.Equal(t, "value #1", r)
		})
		t.Run("#2", func(t *testing.T) {
			t.Parallel()

			p := mockProvider{
				fnWithError: func() (any, error) {
					return "value #2", nil
				},
			}
			r, executed, err := caller.CallProviderMethod(&p, "ProviderWithError", nil, false)
			assert.NoError(t, err)
			assert.True(t, executed)
			assert.Equal(t, "value #2", r)
		})
	})
	t.Run("Errors", func(t *testing.T) {
		t.Parallel()

		t.Run("#1", func(t *testing.T) {
			t.Parallel()

			r, executed, err := caller.CallProviderMethod(nil, "MyProvider", nil, false)
			assert.Nil(t, r)
			assert.False(t, executed)
			assert.EqualError(t, err, `cannot call provider (<nil>)."MyProvider": cannot call method (<nil>)."MyProvider": invalid method receiver: <nil>`)
		})
		t.Run("#2", func(t *testing.T) {
			t.Parallel()

			p := mockProvider{
				fnWithError: func() (any, error) {
					return "error value", errors.New("my error")
				},
			}
			r, executed, err := caller.CallProviderMethod(&p, "ProviderWithError", nil, false)
			assert.Equal(t, "error value", r)
			assert.True(t, executed)
			assert.EqualError(t, err, "provider returned error: my error")
		})
		t.Run("#3", func(t *testing.T) {
			t.Parallel()

			r, executed, err := caller.CallProviderMethod(&mockProvider{}, "NotProvider", nil, false)
			assert.Nil(t, r)
			assert.False(t, executed)
			assert.EqualError(t, err, `cannot call provider (*caller_test.mockProvider)."NotProvider": cannot call method (*caller_test.mockProvider)."NotProvider": second value returned by provider must implement error interface, interface {} given`)
		})
	})
}

func TestCallMethod(t *testing.T) {
	t.Parallel()

	t.Run("Pointer loop", func(t *testing.T) {
		t.Parallel()

		var a any
		a = &a
		r, err := caller.CallMethod(a, "method", nil, false)
		assert.EqualError(t, err, `cannot call method (*interface {})."method": unexpected pointer loop`)
		assert.Nil(t, r)
	})
}

func TestCallWither(t *testing.T) {
	t.Parallel()

	t.Run("Given scenarios", func(t *testing.T) {
		t.Parallel()

		var emptyPerson any = person{}

		scenarios := []struct {
			object any
			wither string
			params []any
			output any
		}{
			{
				object: make(ints, 0),
				wither: "Append",
				params: []any{5},
				output: ints{5},
			},
			{
				object: person{name: "Mary"},
				wither: "WithName",
				params: []any{"Jane"},
				output: person{name: "Jane"},
			},
			{
				object: &person{name: "Mary"},
				wither: "WithName",
				params: []any{"Jane"},
				output: person{name: "Jane"},
			},
			{
				object: emptyPerson,
				wither: "WithName",
				params: []any{"Kaladin"},
				output: person{name: "Kaladin"},
			},
			{
				object: &emptyPerson,
				wither: "WithName",
				params: []any{"Shallan"},
				output: person{name: "Shallan"},
			},
		}

		for i, tmp := range scenarios {
			s := tmp
			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				t.Parallel()

				result, err := caller.CallWither(s.object, s.wither, s.params, false)
				assert.NoError(t, err)
				assert.Equal(t, s.output, result)
			})
		}
	})

	t.Run("Given errors", func(t *testing.T) {
		t.Parallel()

		scenarios := []struct {
			object any
			wither string
			params []any
			error  string
		}{
			{
				object: person{},
				wither: "withName",
				params: nil,
				error:  `cannot call wither (caller_test.person)."withName": cannot call method (caller_test.person)."withName": (caller_test.person)."withName": invalid method`,
			},
			{
				object: person{},
				wither: "Clone",
				params: nil,
				error:  `cannot call wither (caller_test.person)."Clone": cannot call method (caller_test.person)."Clone": wither must return 1 value, given function returns 2 values`,
			},
			{
				object: person{},
				wither: "WithName",
				params: nil,
				error:  `cannot call wither (caller_test.person)."WithName": cannot call method (caller_test.person)."WithName": not enough input arguments`,
			},
		}

		for i, tmp := range scenarios {
			s := tmp
			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				t.Parallel()

				o, err := caller.CallWither(s.object, s.wither, s.params, false)
				assert.Nil(t, o)
				assert.EqualError(t, err, s.error)
			})
		}
	})

	t.Run("Pointer loop", func(t *testing.T) {
		t.Parallel()

		var a any
		a = &a
		r, err := caller.CallWither(a, "method", nil, false)
		assert.EqualError(t, err, `cannot call wither (*interface {})."method": cannot call method (*interface {})."method": unexpected pointer loop`)
		assert.Nil(t, r)
	})

	t.Run("Nil pointer receiver", func(t *testing.T) {
		t.Parallel()

		var p *person
		r, err := caller.CallWither(p, "Empty", nil, false)
		assert.NoError(t, err)
		assert.Nil(t, p)
		assert.Equal(t, person{}, r)
	})
}

type ints []int

func (i ints) Append(v int) ints {
	return append(i, v)
}

type person struct {
	name string
}

func (p person) Clone() (person, error) {
	return p, nil
}

func (p person) WithName(n string) person {
	return person{name: n}
}

func (p person) withName(n string) person { //nolint:unused
	return person{name: n}
}

func (p *person) setName(n string) {
	p.name = n
}

func (p *person) SetName(n string) {
	p.name = n
}

func (p *person) Empty() person {
	return person{}
}

type character struct {
	name string
}

func (c *character) SetName(n string) {
	c.name = n
}

func (c character) WithName(n string) character {
	c.name = n

	return c
}

func (c character) GetName() string {
	return c.name
}

type callerStrct struct{}

func (*callerStrct) Call(fn func()) {
	fn()
}

func TestCallMethod_error(t *testing.T) {
	t.Parallel()

	t.Run("No pointer, expected pointer receiver", func(t *testing.T) {
		t.Parallel()

		c := character{}
		_, err := caller.CallMethod(c, "SetName", []any{"Thor"}, true)
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
		_, err := caller.CallMethod(&object, "Do", nil, true)
		require.EqualError(t, err, `cannot call method (**struct {})."Do": (*struct {})."Do": invalid method`)
	})

	t.Run("Invalid method #2", func(t *testing.T) {
		t.Parallel()

		object := struct{}{}
		_, err := caller.CallMethod(&object, "Do", nil, true)
		require.EqualError(t, err, `cannot call method (*struct {})."Do": (*struct {})."Do": invalid method`)
	})

	t.Run("Cannot convert args", func(t *testing.T) {
		t.Parallel()

		c := character{}
		_, err := caller.CallMethod(&c, "SetName", []any{struct{}{}}, true)
		require.EqualError(
			t,
			err,
			`cannot call method (*caller_test.character)."SetName": arg0: cannot convert struct {} to string`,
		)
	})

	t.Run("Not assignable args", func(t *testing.T) {
		t.Parallel()

		c := character{}
		_, err := caller.CallMethod(&c, "SetName", []any{struct{}{}}, false)
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
		_, err := caller.CallMethod(&c, "SetName", []any{"Bernhard", "Riemann"}, false)
		require.EqualError(
			t,
			err,
			`cannot call method (*caller_test.character)."SetName": too many input arguments`,
		)
	})

	t.Run("Not enough input arguments", func(t *testing.T) {
		t.Parallel()

		c := character{}
		_, err := caller.CallMethod(&c, "SetName", []any{}, false)
		require.EqualError(
			t,
			err,
			`cannot call method (*caller_test.character)."SetName": not enough input arguments`,
		)
	})

	t.Run("Pointer to any(nil)", func(t *testing.T) {
		t.Parallel()

		var c any
		_, err := caller.CallMethod(&c, "Do", nil, true)
		require.EqualError(t, err, `cannot call method (*interface {})."Do": invalid object`)
	})
}

func TestCallMethod_ok(t *testing.T) {
	t.Parallel()

	t.Run("Pointer", func(t *testing.T) {
		t.Parallel()

		c := character{}
		_, err := caller.CallMethod(&c, "SetName", []any{"Thor"}, true)
		require.NoError(t, err)
		assert.Equal(t, c.name, "Thor")
	})

	t.Run("Pointer #2", func(t *testing.T) {
		t.Parallel()

		c := &character{}
		_, err := caller.CallMethod(c, "SetName", []any{"Thor"}, true)
		require.NoError(t, err)
		assert.Equal(t, c.name, "Thor")
	})

	t.Run("Pointer to pointer", func(t *testing.T) {
		t.Parallel()

		c := &character{}
		_, err := caller.CallMethod(&c, "SetName", []any{"Thor"}, true)
		require.NoError(t, err)
		assert.Equal(t, c.name, "Thor")
	})

	t.Run("Pointer to pointer to pointer", func(t *testing.T) {
		t.Parallel()

		tmp := &character{}
		c := &tmp

		_, err := caller.CallMethod(&c, "SetName", []any{"Thor"}, true)
		require.NoError(t, err)
		actual, _ := getter.Get(c, "name")
		assert.Equal(t, actual, "Thor")
	})

	t.Run("Pointer to pointer to pointer to pointer", func(t *testing.T) {
		t.Parallel()

		tmp := &character{}
		tmp2 := &tmp
		c := &tmp2

		_, err := caller.CallMethod(&c, "SetName", []any{"Thor"}, true)
		require.NoError(t, err)
		actual, _ := getter.Get(c, "name")
		assert.Equal(t, actual, "Thor")
	})

	t.Run("Pointer to any", func(t *testing.T) {
		t.Parallel()

		var c any = character{}

		_, err := caller.CallMethod(&c, "SetName", []any{"Thor"}, true)
		require.NoError(t, err)
		actual, _ := getter.Get(c, "name")
		assert.Equal(t, actual, "Thor")
	})

	t.Run("Pointer to any over pointer", func(t *testing.T) {
		t.Parallel()

		var c any = &character{}

		_, err := caller.CallMethod(&c, "SetName", []any{"Thor"}, true)
		require.NoError(t, err)
		actual, _ := getter.Get(c, "name")
		assert.Equal(t, actual, "Thor")
	})

	t.Run("Any over pointer", func(t *testing.T) {
		t.Parallel()

		var c any = &character{}

		_, err := caller.CallMethod(c, "SetName", []any{"Thor"}, true)
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

		_, err := caller.CallMethod(c, "Call", []any{callee}, false)
		require.NoError(t, err)
		assert.True(t, executed)
	})

	t.Run("Getter over value", func(t *testing.T) {
		t.Parallel()

		c := character{name: "Thor"}
		r, err := caller.CallMethod(c, "GetName", nil, false)
		require.NoError(t, err)
		assert.Equal(t, []any{"Thor"}, r)
	})

	t.Run("Getter over pointer", func(t *testing.T) {
		t.Parallel()

		c := character{name: "Thor"}
		r, err := caller.CallMethod(&c, "GetName", nil, false)
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

func TestCallProviderMethod_ok(t *testing.T) {
	t.Parallel()

	t.Run("2+2=4", func(t *testing.T) {
		t.Parallel()

		provider := provider{
			fn: func(args ...any) (any, error) {
				require.Equal(t, []any{2, 2}, args)

				return 4, nil
			},
		}

		result, executed, err := caller.CallProviderMethod(provider, "Provide", []any{2, 2}, true)
		require.NoError(t, err)
		assert.True(t, executed)
		assert.Equal(t, 4, result)
	})
}

func TestCallProviderMethod_error(t *testing.T) {
	t.Parallel()

	t.Run("Not assignable args", func(t *testing.T) {
		t.Parallel()

		db := &sql.DB{}
		result, executed, err := caller.CallProviderMethod(
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
		result, executed, err := caller.CallProviderMethod(
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

func TestCallWither_ok(t *testing.T) {
	t.Parallel()

	c := character{}
	decorated, err := caller.CallWither(c, "WithName", []any{"Jane"}, false)
	require.NoError(t, err)
	require.Equal(
		t,
		character{name: "Jane"},
		decorated,
	)
}

func TestCallWither_errors(t *testing.T) {
	t.Parallel()

	var loop any
	loop = &loop

	scenarios := []struct {
		name        string
		object      any
		wither      string
		args        []any
		convertArgs bool
		error       string
	}{
		{
			name:        "Loop",
			object:      loop,
			wither:      "Loop",
			args:        nil,
			convertArgs: false,
			error:       `cannot call wither (*interface {})."Loop": cannot call method (*interface {})."Loop": unexpected pointer loop`,
		},
		{
			name:        "Setter",
			object:      &character{},
			wither:      "SetName",
			args:        []any{"Jane"},
			convertArgs: false,
			error:       `cannot call wither (*caller_test.character)."SetName": cannot call method (*caller_test.character)."SetName": wither must return 1 value, given function returns 0 values`,
		},
		{
			name:        "Not assignable arguments",
			object:      character{},
			wither:      "WithName",
			args:        []any{123},
			convertArgs: false,
			error:       `cannot call wither (caller_test.character)."WithName": cannot call method (caller_test.character)."WithName": arg0: value of type int is not assignable to type string`,
		},
		{
			name:        "Cannot convert arguments",
			object:      character{},
			wither:      "WithName",
			args:        []any{struct{}{}},
			convertArgs: true,
			error:       `cannot call wither (caller_test.character)."WithName": cannot call method (caller_test.character)."WithName": arg0: cannot convert struct {} to string`,
		},
	}

	for _, s := range scenarios {
		s := s

		t.Run(s.name, func(t *testing.T) {
			t.Parallel()

			decorated, err := caller.CallWither(s.object, s.wither, s.args, s.convertArgs)
			require.EqualError(t, err, s.error)
			require.Nil(t, decorated)
		})
	}
}
