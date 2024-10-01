package reflect_test

import (
	"fmt"
	stdReflect "reflect"
	"testing"

	"github.com/gontainer/reflectpro/internal/reflect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
				callback: func(f stdReflect.StructField, value any) (_ any, set bool) {
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
			{
				strct: person{},
				callback: func(f stdReflect.StructField, value any) (_ any, set bool) {
					if f.Name == "Name" {
						return "Jane", true
					}

					if f.Name == "age" {
						return uint(30), true
					}

					return nil, false
				},
				convert:      false,
				convertToPtr: false,
				error:        `IterateFields: *interface {}: IterateFields: reflect_test.person: field 1 "age": value of type uint is not assignable to type uint8`,
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
