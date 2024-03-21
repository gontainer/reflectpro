[![Go Reference](https://pkg.go.dev/badge/github.com/gontainer/reflectpro.svg)](https://pkg.go.dev/github.com/gontainer/reflectpro)
[![Tests](https://github.com/gontainer/reflectpro/actions/workflows/tests.yml/badge.svg)](https://github.com/gontainer/reflectpro/actions/workflows/tests.yml)
[![Coverage Status](https://coveralls.io/repos/github/gontainer/reflectpro/badge.svg?branch=main)](https://coveralls.io/github/gontainer/reflectpro?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/gontainer/reflectpro)](https://goreportcard.com/report/github.com/gontainer/reflectpro)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=gontainer_reflectpro&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=gontainer_reflectpro)

# Reflectpro

Simple, elegant, and intuitive [callers](caller), [copiers](copier), [getters](getter) and [setters](setter).

## Examples

## Caller

In the following example, we have a pointer to `any` that stores a `struct`,
instead of having a direct pointer to a `struct`.
The receiver is a pointer, so eventually we cannot call the given method.
`Caller` handles that by creating a pointer to a copy of that value.

```
type Person struct {
	name string
	age  int
}

func (p *Person) SetName(n string) {
	p.name = n
}

func Example() {
	var p any
	p := &Person{age: 25}
	_, _ = caller.CallMethod(p, "SetName", []any{"Mary"}, false)
	fmt.Printf("%+v\n", p)
	// Output: &{name:Mary age:25}
}
```

## Copier

```go
var (
    from = []any{int(1), uint(2), float32(3), float64(4)}
    to   []uint64
)
_ = copier.Copy(from, &to, true) // convert
fmt.Printf("%#v\n", to)
// Output: []uint64{0x1, 0x2, 0x3, 0x4}
```

## Getter

In the following example, we read an unexported field of the given struct.

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

## Setter

In the following example, we have a pointer to `any` that stores a `struct`,
instead of having a direct pointer to a `struct`. Since it is an `unaddressable value`,
the `reflect` package from the standard library does not allow assigning a new value to this field.
`Setter` handles that by creating an addressable copy.

```go
var person any
person = struct {
    name string
}{}
_ = setter.Set(&person, "name", "Mary", false)
fmt.Println(person)
// Output: {Mary}
```
