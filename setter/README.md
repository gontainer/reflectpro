# Setter

Package setter allows for write operations on fields of any struct.

```go
person := struct {
    name string
}{}
_ = setter.Set(&person, "name", "Mary", false)
fmt.Println(person.name)
// Output: Mary
```

See [examples](examples_test.go).
