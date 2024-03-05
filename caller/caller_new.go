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

package caller

import (
	"fmt"

	"github.com/gontainer/grouperror"
	"github.com/gontainer/reflectpro/caller/internal/caller"
)

// NewCallProviderMethod will replace [CallProviderMethod] and [ForceCallProviderMethod].
//
// Deprecated: it'll be renamed to `CallProviderMethod`.
func NewCallProviderMethod( //nolint:ireturn
	object any,
	method string,
	args []any,
	convertArgs bool,
) (
	_ any,
	executed bool,
	err error,
) {
	results, err := callMethod(object, method, args, convertArgs, caller.ValidatorProvider)
	if err != nil {
		//nolint:wrapcheck
		return nil, false, grouperror.Prefix(fmt.Sprintf(providerMethodInternalErrPrefix, object, method), err)
	}

	var e error

	if len(results) > 1 {
		// do not panic when results[1] == nil
		e, _ = results[1].(error)
	}

	if e != nil {
		e = grouperror.Prefix(providerExternalErrPrefix, newProviderError(e))
	}

	return results[0], true, e //nolint:wrapcheck
}

// NewCallWither will replace [CallWither] and [ForceCallWither].
//
// Deprecated: it'll be renamed to `CallWither`.
func NewCallWither( //nolint:ireturn
	object any,
	wither string,
	args []any,
	convertArgs bool,
) (
	_ any,
	err error,
) {
	defer func() {
		if err != nil {
			err = grouperror.Prefix(fmt.Sprintf("cannot call wither (%T).%+q: ", object, wither), err)
		}
	}()

	results, err := callMethod(object, wither, args, convertArgs, caller.ValidatorWither)
	if err != nil {
		return nil, err
	}

	return results[0], nil
}
