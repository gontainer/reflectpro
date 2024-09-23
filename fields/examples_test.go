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
