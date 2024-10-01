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

//nolint:lll
func TestIterateFields(t *testing.T) {
	t.Parallel()

	t.Run("set", func(t *testing.T) {
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
				callback: func(f stdReflect.StructField, value any) reflect.FieldCallbackResult {
					if f.Name == "Name" {
						return reflect.FieldCallbackResultSet("Jane")
					}

					if f.Name == "age" {
						return reflect.FieldCallbackResultSet(uint(30))
					}

					return reflect.FieldCallbackResultDontSet()
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
				callback: func(f stdReflect.StructField, value any) reflect.FieldCallbackResult {
					if f.Name == "Name" {
						return reflect.FieldCallbackResultSet("Jane")
					}

					if f.Name == "age" {
						return reflect.FieldCallbackResultSet(uint(30))
					}

					return reflect.FieldCallbackResultDontSet()
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
