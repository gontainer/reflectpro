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
	"testing"

	"github.com/gontainer/reflectpro/caller"
	"github.com/gontainer/reflectpro/getter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Character struct {
	Name string
}

func (c *Character) SetName(n string) {
	c.Name = n
}

type (
	//nolint
	// Deprecated: TODO remove it
	Nuller struct{}
)

func (*Nuller) Null() {
	fmt.Println("null") //nolint
}

func TestNewCallMethod(t *testing.T) {
	t.Parallel()

	t.Run("No pointer", func(t *testing.T) {
		t.Parallel()

		c := Character{}
		_, err := caller.NewCallMethod(c, "SetName", []any{"Thor"}, true)
		require.EqualError(
			t,
			err,
			`cannot call method (caller_test.Character)."SetName": (caller_test.Character)."SetName": invalid method`,
		)
		assert.Empty(t, c.Name)
	})

	t.Run("Pointer", func(t *testing.T) {
		t.Parallel()

		c := Character{}
		_, err := caller.NewCallMethod(&c, "SetName", []any{"Thor"}, true)
		require.NoError(t, err)
		assert.Equal(t, c.Name, "Thor")
	})

	t.Run("var c any = &Character{}; caller.NewCallMethod(&c", func(t *testing.T) {
		t.Parallel()

		var c any = &Character{}

		_, err := caller.NewCallMethod(&c, "SetName", []any{"Thor"}, true)
		require.NoError(t, err)
		actual, _ := getter.Get(c, "Name")
		assert.Equal(t, actual, "Thor")
	})

	t.Run("var c = Character{}; caller.NewCallMethod(&c", func(t *testing.T) {
		t.Parallel()

		var c any = Character{}

		_, err := caller.NewCallMethod(&c, "SetName", []any{"Thor"}, true)
		require.NoError(t, err)
		actual, _ := getter.Get(c, "Name")
		assert.Equal(t, actual, "Thor")
	})

	t.Run("var c = Character{}; var d any = &c; caller.NewCallMethod(&d", func(t *testing.T) {
		t.Parallel()

		var c any = Character{}

		var d any = &c

		_, err := caller.NewCallMethod(d, "SetName", []any{"Thor"}, true)
		require.NoError(t, err)
		actual, _ := getter.Get(c, "Name")
		assert.Equal(t, actual, "Thor")
	})

	t.Run("invalid method", func(t *testing.T) {
		t.Parallel()

		_, err := caller.NewCallMethod(struct{}{}, "Do", nil, true)
		//nolint
		// TODO re-format errors to remove the duplication
		require.EqualError(t, err, `cannot call method (struct {})."Do": (struct {})."Do": invalid method`)
	})

	t.Run("Nuller", func(t *testing.T) {
		t.Parallel()

		var x *Nuller

		_, err := caller.NewCallMethod(x, "Null", nil, true)
		require.NoError(t, err)
	})
}
