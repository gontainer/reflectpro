[![Go Reference](https://pkg.go.dev/badge/github.com/gontainer/reflectpro.svg)](https://pkg.go.dev/github.com/gontainer/reflectpro)
[![Tests](https://github.com/gontainer/reflectpro/actions/workflows/tests.yml/badge.svg)](https://github.com/gontainer/reflectpro/actions/workflows/tests.yml)
[![Coverage Status](https://coveralls.io/repos/github/gontainer/reflectpro/badge.svg?branch=main)](https://coveralls.io/github/gontainer/reflectpro?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/gontainer/reflectpro)](https://goreportcard.com/report/github.com/gontainer/reflectpro)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=gontainer_reflectpro&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=gontainer_reflectpro)

# Reflectpro

Simple, elegant, and intuitive [callers](caller), [copiers](copier), [getters](getter) and [setters](setter).
<br/>This package has been extracted from [github.com/gontainer/gontainer-helpers](https://github.com/gontainer/gontainer-helpers).

## Examples

## Caller

```go
type Person struct {
	name string
}

func (p *Person) SetName(n string) {
	p.name = n
}

func ExampleCallMethod() {
	p := &Person{}
	_, _ = caller.CallMethod(p, "SetName", []any{"Mary"}, false)
	fmt.Println(p.name)
	// Output: Mary
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

```go
type Person struct {
    Name string
}
p := Person{}
_ = setter.Set(&p, "Name", "Jane", false)
fmt.Println(p) // {Jane}
```
