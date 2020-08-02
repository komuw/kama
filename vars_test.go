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

func (p Person) somePrivateMethodOne()  {}
func (p *Person) somePrivateMethodTwo() {}
func (p Person) ValueMethodOne()        {}
func (p *Person) PtrMethodOne()         {}
func (p Person) ValueMethodTwo()        {}
func (p *Person) PtrMethodTwo() float32 { return p.Height }

func ThisFunction(arg1 string, arg2 int) (string, error) {
	return "", nil
}

var thisFunctionVar = ThisFunction

type customerID uint16

func (c customerID) Id() uint16 { return uint16(c) }

func TestBasicVariables(t *testing.T) {
	tt := []struct {
		variable interface{}
		expected vari
	}{
		{
			Person{Name: "John"}, vari{
				Name:      "github.com/komuw/kama.Person",
				Kind:      reflect.Struct,
				Signature: []string{"main.Person", "*main.Person"},
				Fields:    []string{"Name", "Age", "Height"},
				Methods:   []string{"ValueMethodOne func(main.Person)", "ValueMethodTwo func(main.Person)", "PtrMethodOne func(*main.Person)", "PtrMethodTwo func(*main.Person) float32"},
			},
		},
		{

			&Person{Name: "Jane"}, vari{
				Name:      ".Person",
				Kind:      reflect.Struct,
				Signature: []string{"*main.Person", "main.Person"},
				Fields:    []string{"Name", "Age", "Height"},
				Methods:   []string{"ValueMethodOne func(main.Person)", "ValueMethodTwo func(main.Person)", "PtrMethodOne func(*main.Person)", "PtrMethodTwo func(*main.Person) float32"},
			},
		},
		{
			ThisFunction, vari{
				Name:      "github.com/komuw/kama.ThisFunction",
				Kind:      reflect.Func,
				Signature: []string{"func(string, int) (string, error)"},
				Fields:    []string{},
				Methods:   []string{},
			},
		},
		{
			thisFunctionVar, vari{
				Name:      "github.com/komuw/kama.ThisFunction",
				Kind:      reflect.Func,
				Signature: []string{"func(string, int) (string, error)"},
				Fields:    []string{},
				Methods:   []string{},
			},
		},
		{
			customerID(9), vari{
				Name:      "github.com/komuw/kama.customerID",
				Kind:      reflect.Uint16,
				Signature: []string{"main.customerID"},
				Fields:    []string{},
				Methods:   []string{"Id func(main.customerID) uint16"},
			},
		},
	}

	for _, v := range tt {
		res := newVari(v.variable)

		if !cmp.Equal(res, v.expected) {
			t.Errorf("\ngot \n\t%#+v \nwanted \n\t%#+v", res, v.expected)
		}
	}
}
