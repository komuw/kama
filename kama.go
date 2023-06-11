// Package kama exposes exported information of types, variables, packages, modules, imports etc
// It also pretty prints data structures.
//
// It can be used to aid debugging and testing.
//
//	import "github.com/komuw/kama"
//
//	kama.Dirp("compress/flate")
//	kama.Dirp(&http.Request{})
//	kama.Dirp("github.com/pkg/errors")
package kama

import (
	"fmt"
	"reflect"
	"strings"
)

// Dirp prints (to stdout) exported information of types, variables, packages, modules, imports
// It also pretty prints data structures.
//
// examples:
//
//	import "github.com/komuw/kama"
//
//	kama.Dirp("compress/flate")
//	kama.Dirp(&http.Request{})
//	kama.Dirp("github.com/pkg/errors")
//	kama.Dirp(http.Request{})
func Dirp(i interface{}) {
	fmt.Println(Dir(i))
}

// Dir returns exported information of types, variables, packages, modules, imports
func Dir(i interface{}) string {
	iType := reflect.TypeOf(i)
	if iType == nil {
		res := newVari(i)
		return res.String()
	} else if iType.Kind() == reflect.String {
		pat, ok := i.(string)
		if !ok {
			res := newVari(i)
			return res.String()
		}

		res, err := newPak(pat)
		if err != nil {
			for _, eMsg := range []string{
				// We check if it is truly a module error. We check all the errors that can be returned.
				// Unfortunately `ImportMissingError` is an internal error so we cant use errors.Is/As
				// https://github.com/golang/go/blob/go1.20.5/src/cmd/go/internal/modload/import.go#L57-L96
				//
				// This list will need to be kept uptodate with Go versions
				"is not in GOROOT",
				"cannot find module providing package",
				"cannot find module",
				"is replaced but not required",
				"no required module",
				"but not at required version",
				"missing module",
			} {
				if strings.Contains(err.Error(), eMsg) {
					return err.Error()
				}
			}

			// If it is not a module error, then probably `i` is a variable of type string.
			// Thus we need to create a `kama.vari`
			res := newVari(i)
			return res.String()
		}
		return res.String()
	} else {
		res := newVari(i)
		return res.String()
	}
}

// Stackp prints to standard error the colorized stack trace returned by debug.Stack.
//
// Stack trace from the runtime/stdlib is colored cyan, third party libraries is magenta
// whereas your code is green.
func Stackp() {
	stackp()
}
