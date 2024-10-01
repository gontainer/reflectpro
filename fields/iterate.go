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

package fields

import (
	"fmt"
	"reflect"

	intReflect "github.com/gontainer/reflectpro/internal/reflect"
)

type config struct {
	setter            func(_ Path, value any) (_ any, ok bool)
	getter            func(_ Path, value any)
	prefillNilStructs bool
	convertTypes      bool
	convertToPtr      bool
	recursive         bool
}

func newConfig(opts ...Option) *config {
	c := &config{
		setter:            nil,
		getter:            nil,
		prefillNilStructs: false,
		convertTypes:      false,
		convertToPtr:      false,
		recursive:         false,
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

type Option func(*config)

func PrefillNilStructs(v bool) Option {
	return func(c *config) {
		c.prefillNilStructs = v
	}
}

func Setter(fn func(path Path, value any) (_ any, set bool)) Option {
	return func(c *config) {
		c.setter = fn
	}
}

func Getter(fn func(_ Path, value any)) Option {
	return func(c *config) {
		c.getter = fn
	}
}

func ConvertTypes(v bool) Option {
	return func(c *config) {
		c.convertTypes = v
	}
}

func ConvertToPointers(v bool) Option {
	return func(c *config) {
		c.convertToPtr = v
	}
}

func Recursive(v bool) Option {
	return func(c *config) {
		c.recursive = v
	}
}

func Iterate(strct any, opts ...Option) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("fields.Iterate: %w", err)
		}
	}()

	return iterate(strct, newConfig(opts...), nil)
}

//nolint:wrapcheck
func iterate(strct any, cfg *config, path []reflect.StructField) error {
	var fn intReflect.FieldCallback

	var finalErr error

	fn = func(f reflect.StructField, value any) (_ any, set bool) {
		if finalErr != nil {
			return nil, false
		}

		// call getter
		if cfg.getter != nil {
			cfg.getter(append(path, f), value)
		}

		var setterHasBeenTriggered bool

		value, setterHasBeenTriggered = trySetValue(f, value, cfg, path)

		if cfg.recursive && isStructOrNonNilStructPtr(f.Type, value) {
			original := value

			if err := iterate(&value, cfg, append(path, f)); err != nil {
				finalErr = fmt.Errorf("%s: %w", f.Name, err)

				return nil, false
			}

			if !reflect.DeepEqual(original, value) {
				setterHasBeenTriggered = true
			}
		}

		if setterHasBeenTriggered {
			return value, true
		}

		return nil, false
	}

	err := intReflect.IterateFields(
		strct,
		fn,
		cfg.convertTypes,
		cfg.convertToPtr,
	)

	if err != nil {
		return err
	}

	if finalErr != nil {
		return finalErr
	}

	return nil
}

func trySetValue( //nolint:ireturn
	f reflect.StructField,
	value any,
	cfg *config,
	path []reflect.StructField,
) (
	_ any,
	set bool,
) {
	// Call setter
	if cfg.setter != nil {
		if newVal, ok := cfg.setter(append(path, f), value); ok {
			return newVal, true
		}
	}

	// Set pointer to a zero-value struct
	if cfg.prefillNilStructs &&
		f.Type.Kind() == reflect.Ptr && f.Type.Elem().Kind() == reflect.Struct &&
		reflect.ValueOf(value).IsZero() {
		return reflect.New(f.Type.Elem()).Interface(), true
	}

	return value, false
}

// isStructOrNonNilStructPtr checks if the given type is a struct or a non-nil pointer to a struct.
func isStructOrNonNilStructPtr(t reflect.Type, v any) bool {
	return t.Kind() == reflect.Struct ||
		(t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct && !reflect.ValueOf(v).IsZero())
}
