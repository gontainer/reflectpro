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

type callerError struct {
	error
}

func newCallerError(error error) *callerError {
	return &callerError{error: error}
}

func (e *callerError) Collection() []error {
	return []error{e.error}
}

func (e *callerError) Unwrap() error {
	return e.error
}

//nolint
// TODO change the doc comment for [ProviderError]

/*
ProviderError wraps errors returned by providers in [CallProvider].

	type myError struct {
		error
	}

	p := func() (any, error) {
		return nil, &myError{errors.New("my error")}
	}

	_, err := caller.CallProvider(p, nil, false)
	if err != nil {
		var providerErr *caller.ProviderError
		if errors.As(err, &providerErr) {
			fmt.Println("provider returned error:", providerErr)
		} else {
			fmt.Println("provider wasn't invoked:", err)
		}
	}
*/
type ProviderError struct {
	*callerError
}

func newProviderError(err error) *ProviderError {
	return &ProviderError{callerError: newCallerError(err)}
}
