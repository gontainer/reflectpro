package caller

import (
	"reflect"
)

type reflectType struct {
	reflect.Type
}

// In works almost same as [reflect.Type.In],
// but it returns t.In(t.NumIn() - 1).Elem() for t.isVariadic() && i >= t.NumIn().
func (t reflectType) In(i int) reflect.Type {
	last := t.NumIn() - 1
	if i > last {
		i = last
	}

	r := t.Type.In(i)

	if t.IsVariadic() && i == last {
		r = r.Elem()
	}

	return r
}
