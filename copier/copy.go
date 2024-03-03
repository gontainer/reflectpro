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

package copier

import (
	"fmt"
	"reflect"

	intReflect "github.com/gontainer/reflectpro/internal/reflect"
)

/*
Copy copies a value of `from` to `to`. The third parameter says whether the value should be converted whenever possible.

	from := 5
	b := 0
	Copy(from, &to, false)
	fmt.Println(to) // 5
*/
func Copy(from any, to any, convert bool) error {
	t := reflect.ValueOf(to)

	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("expected %s, %T given", reflect.Ptr.String(), to)
	}

	f, err := intReflect.ValueOf(from, t.Elem().Type(), convert)
	if err != nil {
		return err //nolint:wrapcheck
	}

	t.Elem().Set(f)

	return nil
}
