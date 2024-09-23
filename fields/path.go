package fields

import (
	"reflect"
)

type Path []reflect.StructField

func (p Path) Names() []string {
	if p == nil {
		return nil
	}

	r := make([]string, len(p))
	for i, x := range p {
		r[i] = x.Name
	}

	return r
}

func (p Path) HasSuffix(path ...string) bool {
	if len(p) < len(path) {
		return false
	}

	for i := 0; i < len(path); i++ {
		if p[len(p)-1-i].Name != path[len(path)-1-i] {
			return false
		}
	}

	return true
}

func (p Path) CompareStrings(path ...string) bool {
	if len(p) != len(path) {
		return false
	}

	for i := 0; i < len(p); i++ {
		if p[i].Name != path[i] {
			return false
		}
	}

	return true
}
