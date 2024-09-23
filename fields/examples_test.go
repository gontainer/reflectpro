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
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gontainer/reflectpro/copier"
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

func ExampleSet() {
	p := TrainingPlan{}

	_ = fields.Iterate(
		&p,
		fields.Setter(func(path fields.Path, value any) (_ any, ok bool) {
			switch {
			case path.EqualNames("TrainingPlanMeta", "Name"):
				return "My training plan", true
			case path.EqualNames("Monday", "Name"):
				return "pushups", true
			case path.EqualNames("Tuesday", "name"):
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
		fields.Setter(func(path fields.Path, value any) (_ any, ok bool) {
			if path.EqualNames("os") {
				return "Android", true
			}

			return nil, false
		}),
	)

	fmt.Println(p.os)

	// Output:
	// Android
}

type MyCache struct {
	TTL time.Duration
}

type MyConfig struct {
	MyCache *MyCache
}

func ExamplePrefillNilStructs() {
	cfg := MyConfig{}

	_ = fields.Iterate(
		&cfg,
		fields.Setter(func(path fields.Path, value any) (_ any, ok bool) {
			if path.EqualNames("MyCache", "TTL") {
				return time.Minute, true
			}

			return nil, false
		}),
		fields.PrefillNilStructs(true),
		fields.Recursive(),
	)

	fmt.Println(cfg.MyCache.TTL)

	// Output:
	// 1m0s
}

type CTO struct {
	Salary int
}

type Company struct {
	CTO CTO
}

func ExampleGetter() {
	c := Company{
		CTO: CTO{
			Salary: 1000000,
		},
	}

	var salary int

	_ = fields.Iterate(
		c,
		fields.Getter(func(p fields.Path, value any) {
			if p.EqualNames("CTO", "Salary") {
				_ = copier.Copy(value, &salary, false)
			}
		}),
		fields.Recursive(),
	)

	fmt.Println(salary)

	// Output:
	// 1000000
}

func ExampleConvertToPointers() {
	var cfg struct {
		TTL *time.Duration // expect a pointer
	}

	_ = fields.Iterate(
		&cfg,
		fields.Setter(func(path fields.Path, value any) (_ any, ok bool) {
			if path.EqualNames("TTL") {
				return time.Minute, true // return a value
			}

			return nil, false
		}),
		fields.ConvertToPointers(), // this line will instruct the library to convert values to pointers
	)

	fmt.Println(*cfg.TTL)

	// Output:
	// 1m0s
}

func ExampleReadJSON() {
	var person struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Age       uint   `json:"age"`
		Bio       string `json:"-"`
	}

	// read data from JSON...
	js := `
{
	"firstname": "Jane",
	"lastname": "Doe",
	"age": 30,
	"bio": "bio..."
}`
	var data map[string]any
	_ = json.Unmarshal([]byte(js), &data)

	// populate the data from JSON to the `person` variable,
	// use struct tags, to determine the correct relations
	_ = fields.Iterate(
		&person,
		fields.Setter(func(p fields.Path, value any) (_ any, ok bool) {
			tag, ok := p[len(p)-1].Tag.Lookup("json")
			if !ok {
				return nil, false
			}

			name := strings.Split(tag, ",")[0]
			if name == "-" {
				return nil, false
			}

			if fromJSON, ok := data[name]; ok {
				return fromJSON, true
			}

			return nil, false
		}),
		fields.ConvertTypes(),
	)

	fmt.Printf("%+v\n", person)

	// Output:
	// {Firstname:Jane Lastname:Doe Age:30 Bio:}
}
