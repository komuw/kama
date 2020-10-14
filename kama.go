// Package kama exposes exported information of types, variables, packages, modules, imports etc
// It also pretty prints data structures.
//
// It can be used to aid debugging and testing.
//
//     import "github.com/komuw/kama"
//
//     kama.Dirp("compress/flate")
//     kama.Dirp(&http.Request{})
//     kama.Dirp("github.com/pkg/errors")
//
package kama

import (
	"fmt"
	"reflect"
)

// Dirp prints exported information of types, variables, packages, modules, imports
// It also pretty prints data structures.
//
// It is almost similar to Python's builtin dir function
//
// examples:
//
//     import "github.com/komuw/kama"
//
//     kama.Dirp("compress/flate")
//     kama.Dirp(&http.Request{})
//     kama.Dirp("github.com/pkg/errors")
//     kama.Dirp(http.Request{})
//
func Dirp(i interface{}) {
	fmt.Println(Dir(i))
}

// Dir returns exported information of types, variables, packages, modules, imports
//
func Dir(i interface{}) string {
	iType := reflect.TypeOf(i)
	if iType == nil {
		res := newVari(i)
		return res.String()
	} else if iType.Kind() == reflect.String {
		i := i.(string)
		res, err := newPak(i)
		if err != nil {
			res := newVari(i)
			return res.String()
		}
		return res.String()
	} else {
		res := newVari(i)
		return res.String()
	}
}
