# Caller

This package provides functions that allow calling other functions with unknown arguments.
The `convertArgs` argument instructs the caller whether you allow for converting the type.

```go
sum := func(a, b int) int {
    return a + b
}

returns, _ := caller.Call(sum, []any{2, 3}, false)
fmt.Println(returns) // [5]
```

See [examples](examples_test.go).
