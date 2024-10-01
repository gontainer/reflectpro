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
	"reflect"
)

// Path is built over [reflect.StructField],
// that exports us useful details like:
//   - [reflect.StructField.Name]
//   - [reflect.StructField.Anonymous]
//   - [reflect.StructField.Tag]
//   - [reflect.StructField.Type]
type Path []reflect.StructField

func (p Path) Names() []string {
	if p == nil {
		return nil
	}

	r := make([]string, len(p))
	for i, x := range p {
		r[i] = x.Name
	}

	return r
}

func (p Path) HasSuffix(path ...string) bool {
	if len(p) < len(path) {
		return false
	}

	for i := 0; i < len(path); i++ {
		if p[len(p)-1-i].Name != path[len(path)-1-i] {
			return false
		}
	}

	return true
}

func (p Path) EqualNames(path ...string) bool {
	if len(p) != len(path) {
		return false
	}

	for i := 0; i < len(p); i++ {
		if p[i].Name != path[i] {
			return false
		}
	}

	return true
}
