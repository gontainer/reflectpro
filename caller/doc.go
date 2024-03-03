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

/*
Package caller provides functions that allow calling other functions with unknown arguments.

# Example

	sum := func(a, b int) int {
		return a + b
	}

	returns, _ := caller.Call(sum, []any{2, 3}, false)
	fmt.Println(returns) // [5]

# Provider

It is a function that returns 1 or 2 values. The first value is the desired output of the provider.
The optional second value may contain information about a potential error.

Provider that does not return any error:

	func NewHttpClient(timeout time.Duration) *http.Client {
		return &http.Client{
			Timeout: timeout,
		}
	}

	// httpClient, _ := caller.CallProvider(NewHttpClient, time.Minute)

Provider that may return an error:

	func NewDB(username string, password string) (any, error) {
	    db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/test", username, password))
	    if err != nil {
	         return nil, err
	    }

	    db.SetConnMaxLifetime(time.Minute * 3)
	    db.SetMaxOpenConns(10)
	    db.SetMaxIdleConns(10)

	    return db, nil
	}

	// db, err := caller.CallProvider(NewDB, []any{"root", "root"}, false)
	// if err != nil {
	// 	panic(err)
	// }

# Wither

It is a method that returns one value always:

	type Person struct {
		Name string
	}

	func (p Person) WithName(n string) Person { // it is a wither
		p.Name = n
		return p
	}

	// p, _ := caller.CallWither(caller.Person{}, "WithName", []any{"Jane"}, false)
	// fmt.Println(p) // {Jane}
*/
package caller
