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
	setter            func(path []reflect.StructField, value any) (_ any, ok bool)
	getter            func(path []reflect.StructField, value any)
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

func Setter(fn func(path []reflect.StructField, value any) (_ any, ok bool)) Option {
	return func(c *config) {
		c.setter = fn
	}
}

func Getter(fn func(path []reflect.StructField, value any)) Option {
	return func(c *config) {
		c.getter = fn
	}
}

func ConvertTypes() Option {
	return func(c *config) {
		c.convertTypes = true
	}
}

func ConvertToPointers() Option {
	return func(c *config) {
		c.convertToPtr = true
	}
}

func Recursive() Option {
	return func(c *config) {
		c.recursive = true
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

func iterate(strct any, cfg *config, path []reflect.StructField) error {
	var fn intReflect.FieldCallback

	var stopErr error

	fn = func(f reflect.StructField, value interface{}) (_ interface{}, ok bool) {
		if stopErr != nil {
			return nil, false
		}

		// call getter
		if cfg.getter != nil {
			cfg.getter(append(path, f), value)
		}

		setterHasBeenTriggered := false

		// call setter
		if cfg.setter != nil {
			newVal, ok := cfg.setter(append(path, f), value)
			if ok {
				value, setterHasBeenTriggered = newVal, true
			}
		}

		// set poiner to a zero-value
		if !setterHasBeenTriggered &&
			cfg.prefillNilStructs &&
			value == nil &&
			f.Type.Kind() == reflect.Ptr && f.Type.Elem().Kind() == reflect.Struct {
			value = reflect.New(f.Type.Elem())
		}

		if cfg.recursive {
			if f.Type.Kind() == reflect.Struct || // value is a struct
				(f.Type.Kind() == reflect.Ptr && f.Type.Elem().Kind() == reflect.Struct && value != nil) { // value is a pointer to a non-nil struct

				original := value

				newPath := append(path, f) // TODO maybe copy instead of append
				if err := iterate(&value, cfg, newPath); err != nil {
					// TODO decorate it
					stopErr = err

					return nil, false
				}

				if !reflect.DeepEqual(original, value) {
					setterHasBeenTriggered = true
				}
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

	if stopErr != nil {
		return stopErr
	}

	return nil
}
