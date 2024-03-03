package caller_test

import (
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
}
