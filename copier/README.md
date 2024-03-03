# Copier

Package copier allows for copying a value to a variable with an unknown type.
The third argument instructs the copier whether you allow for converting the type.

**Copy value**

```go
var (
	from = 5 // the type of the variable `to` can be different from the type of the variable `from`
	to   any // as long as the value of the `from` is assignable to the `to`
)
_ = copier.Copy(from, &to, false)
fmt.Println(to)
// Output: 5
```

**Convert & copy value**

```go
var (
	from = int(5) // uint is not assignable to int,
	to   uint     // but [copier.Copy] can convert the type
)
_ = copier.Copy(from, &to, true)
fmt.Println(to)
// Output: 5
```

**More sophisticated examples**

```go
package main

import (
	"fmt"

	"github.com/gontainer/gontainer-helpers/v3/copier"
)

type Person struct {
	Name string
}

func main() {
	{
		var (
			from = []any{Person{Name: "Jane"}, Person{Name: "John"}}
			to   []Person
		)
		_ = copier.Copy(from, &to, true)
		fmt.Println(to)
		// Output: [{Jane} {John}]
	}

	{
		var (
			from = []any{int(1), uint(2), float32(3), float64(4)}
			to   []uint64
		)
		_ = copier.Copy(from, &to, true) // convert
		fmt.Printf("%#v\n", to)
		// Output: []uint64{0x1, 0x2, 0x3, 0x4}

		err := copier.Copy(from, &to, false) // don't convert
		fmt.Println(err)
		// Output: value of type []interface {} is not assignable to type []uint64
	}

	{
		var (
			from = []any{1, "2", 3, 4}
			to   []int
		)
		err := copier.Copy(from, &to, true)
		fmt.Println(err)
		// Output: cannot convert []interface {} to []int: #1: cannot convert string to int
	}
	
	{
		var (
			from = map[int64]any{0: "Jane", 1: "John"}
			to   map[int32]string // let's convert keys and values
		)
		_ = copier.Copy(from, &to, true)
		fmt.Println(to)
		// Output:
		// map[0:Jane 1:John]
    }
}
```
