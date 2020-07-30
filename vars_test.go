package main

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type Person struct {
	Name             string
	Age              int
	Height           float32
	somePrivateField string
}

func (p Person) ValueMethodOne() {}
func (p *Person) PtrMethodOne()  {}
func (p Person) ValueMethodTwo() {}

func TestVariables(t *testing.T) {
	tt := []struct {
		variable interface{}
		expected vari
	}{
		{
			Person{Name: "John"}, vari{
				Name:      "github.com/komuw/dir.Person",
				Kind:      reflect.Struct,
				Signature: "main.Person",
				Fields:    []string{"Name", "Age", "Height"},
				// TODO: `Methods` should include `PtrMethodOne`
				Methods: []string{"ValueMethodOne func(main.Person)", "ValueMethodTwo func(main.Person)"},
			},
		},
		{

			&Person{Name: "Jane"}, vari{
				Name:      ".Person",
				Kind:      reflect.Ptr,
				Signature: "*main.Person",
				Fields:    []string{},
				// TODO: `Methods` should be unified with that of Person{} above
				Methods: []string{"PtrMethodOne func(*main.Person)", "ValueMethodOne func(*main.Person)", "ValueMethodTwo func(*main.Person)"},
			},
		},
	}

	for _, v := range tt {
		res := newVari(v.variable)

		if !cmp.Equal(res, v.expected) {
			t.Errorf("\ngot %#+v \nwanted %#+v", res, v.expected)
		}
	}
}
