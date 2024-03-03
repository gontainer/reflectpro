# Getter

Package getter allows for read operations on fields of any struct.

```go
person := struct {
    name string
}{
    name: "Mary",
}
v, _ := getter.Get(person, "name")
fmt.Println(v)
// Output: Mary
```

See [examples](examples_test.go).
