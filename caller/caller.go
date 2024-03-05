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

// Call calls the given function with the given arguments.
// It returns values returned by the function in a slice.
// If the third argument equals true, it converts types whenever it is possible.
//
//nolint:wrapcheck
func Call(fn any, args []any, convertArgs bool) (_ []any, err error) {
	defer func() {
		if err != nil {
			err = grouperror.Prefix(fmt.Sprintf("cannot call %T: ", fn), err)
		}
	}()

	v, err := caller.Func(fn)
	if err != nil {
		return nil, err
	}

	return caller.CallFunc(v, args, convertArgs)
}

const (
	providerInternalErrPrefix       = "cannot call provider %T: "
	providerMethodInternalErrPrefix = "cannot call provider (%T).%+q: "
	providerExternalErrPrefix       = "provider returned error: "
)

/*
CallProvider works similar to [Call] with the difference it requires a provider as the first argument.
Provider is a function which returns 1 or 2 values.
The second return value which is optional must be a type of error.
See [ProviderError].

	p := func() (any, error) {
	    db, err := sql.Open("mysql", "user:password@/dbname")
	    if err != nil {
	         return nil, err
	    }

	    db.SetConnMaxLifetime(time.Minute * 3)
	    db.SetMaxOpenConns(10)
	    db.SetMaxIdleConns(10)

	    return db, nil
	}

	db, err := caller.CallProvider(p, nil, false)
*/
//nolint:wrapcheck
func CallProvider(provider any, args []any, convertArgs bool) (_ any, executed bool, err error) { //nolint:ireturn
	defer func() {
		if !executed && err != nil {
			err = grouperror.Prefix(fmt.Sprintf(providerInternalErrPrefix, provider), err)
		}
	}()

	fn, err := caller.Func(provider)
	if err != nil {
		return nil, false, err
	}

	if err := caller.ValidatorProvider.Validate(fn); err != nil {
		return nil, false, err
	}

	results, err := caller.CallFunc(fn, args, convertArgs)
	if err != nil {
		return nil, false, err
	}

	r := results[0]

	var e error

	if len(results) > 1 {
		// do not panic when results[1] == nil
		e, _ = results[1].(error)
	}

	if e != nil {
		e = grouperror.Prefix(providerExternalErrPrefix, newProviderError(e))
	}

	return r, true, e
}

/*
CallProviderMethod works similar to [CallProvider], but the provider must be a method on the given object.

	db, _ := sql.Open("mysql", "user:password@/dbname")
	tx, err := caller.CallProviderMethod(db, "Begin", nil, false)
*/
func CallProviderMethod( //nolint:ireturn
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

/*
CallMethod works similar to [Call] with the difference it calls the method by the name over the given receiver.

	type Person struct {
		Name string
	}

	func (p *Person) SetName(n string) {
		p.Name = n
	}

	func main() {
		p := &Person{}
		_, _ = caller.CallMethod(p, "SetName", []any{"Mary"}, false)
		fmt.Println(p.name)
		// Output: Mary
	}
*/
func CallMethod(object any, method string, args []any, convertArgs bool) (_ []any, err error) {
	return callMethod(object, method, args, convertArgs, caller.DontValidate)
}

/*
CallWither works similar to [CallMethod] with the difference the method must be a wither.

	type Person struct {
	    name string
	}

	func (p Person) WithName(n string) Person {
	    p.name = n
	    return p
	}

	func main() {
	    p := Person{}
	    p2, _ := caller.CallWither(p, "WithName", "Mary")
	    fmt.Printf("%+v", p2) // {name:Mary}
	}
*/
func CallWither(object any, wither string, args []any, convertArgs bool) (_ any, err error) { //nolint:ireturn
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
