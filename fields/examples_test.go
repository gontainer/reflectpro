package fields_test

import (
	"fmt"
	"reflect"

	"github.com/davecgh/go-spew/spew"
	"github.com/gontainer/reflectpro/fields"
)

type Exercise struct {
	Name string
}

type TrainingPlanMeta struct {
	Name string
}

type TrainingPlan struct {
	TrainingPlanMeta

	Monday  Exercise
	Tuesday Exercise
}

func comparePaths(reflectPath []reflect.StructField, names ...string) bool {
	if len(reflectPath) != len(names) {
		return false
	}

	for i := 0; i < len(reflectPath); i++ {
		if reflectPath[i].Name != names[i] {
			return false
		}
	}

	return true
}

func ExampleSet() {
	p := TrainingPlan{}

	_ = fields.Iterate(
		&p,
		fields.Setter(func(path []reflect.StructField, value any) (_ any, ok bool) {
			switch {
			case comparePaths(path, "TrainingPlanMeta", "Name"):
				return "My training plan", true
			case comparePaths(path, "Monday", "Name"):
				return "pushups", true
			case comparePaths(path, "Tuesday", "name"):
				return "pullups", true
			}

			return nil, false
		}),
		fields.Recursive(),
	)

	spew.Dump(p)

	// Output:
	// (fields_test.TrainingPlan) {
	//  TrainingPlanMeta: (fields_test.TrainingPlanMeta) {
	//   Name: (string) (len=16) "My training plan"
	//  },
	//  Monday: (fields_test.Exercise) {
	//   Name: (string) (len=7) "pushups"
	//  },
	//  Tuesday: (fields_test.Exercise) {
	//   Name: (string) ""
	//  }
	// }
}

type Phone struct {
	os string
}

func ExampleSetUnexported() {
	p := Phone{}
	_ = fields.Iterate(
		&p,
		fields.Setter(func(path []reflect.StructField, value any) (_ any, ok bool) {
			if path[len(path)-1].Name == "os" {
				return "Android", true
			}

			return nil, false
		}),
	)

	fmt.Println(p.os)

	// Output:
	// Android
}
