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

package caller //nolint:testpackage

import (
	"errors"
	"testing"

	"github.com/gontainer/grouperror"
	errAssert "github.com/gontainer/grouperror/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderError_Collection(t *testing.T) {
	t.Parallel()

	err := newProviderError(newCallerError(
		grouperror.Prefix("prefix: ", errors.New("error 1"), errors.New("error 2")),
	))
	expected := []string{
		`prefix: error 1`,
		`prefix: error 2`,
	}
	errAssert.EqualErrorGroup(t, err, expected)
}

type providerFunc func() (any, error)

func (p providerFunc) Provide() (any, error) { //nolint:ireturn
	return p()
}

func TestProviderError_Unwrap(t *testing.T) {
	t.Parallel()

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		originalErr := errors.New("my error")

		_, err := CallProvider(
			func() (any, error) {
				return nil, originalErr
			},
			nil,
			true,
		)

		require.EqualError(t, err, "provider returned error: my error")
		var providerErr *ProviderError
		require.True(t, errors.As(err, &providerErr))
		require.Same(t, originalErr, providerErr.Unwrap())
	})

	t.Run("Method", func(t *testing.T) {
		t.Parallel()

		originalErr := errors.New("my error")

		_, err := CallProviderMethod(
			providerFunc(func() (any, error) {
				return nil, originalErr
			}),
			"Provide",
			nil,
			true,
		)

		require.EqualError(t, err, "provider returned error: my error")
		var providerErr *ProviderError
		require.True(t, errors.As(err, &providerErr))
		require.Same(t, originalErr, providerErr.Unwrap())
	})
}
