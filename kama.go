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
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"
)

var (
	cfg = Config{ //nolint:gochecknoglobals
		MaxLength:         14,
		ShowPrivateFields: false,
		MaxIndentLevel:    10,
	}

	onceCfg = &sync.Once{} //nolint:gochecknoglobals
)

// Config controls how printing is going to be done.
type Config struct {
	// MaxLength is the length of slices/maps/strings that is going to be dumped.
	MaxLength int
	// ShowPrivateFields dictates whether private struct fields will be dumped.
	ShowPrivateFields bool
	// MaxIndentLevel is the maximum level of indentation/recursiveness to dump to.
	// This is especially important to set if the thing you are dumping has circular references.
	MaxIndentLevel int
}

// Dirp prints (to stdout) exported information of types, variables, packages, modules, imports.
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
//	kama.Dirp(http.Request{}, Config{999})
func Dirp(i interface{}, c ...Config) {
	fmt.Println(Dir(i, c...))
}

// Dir returns exported information of types, variables, packages, modules, imports.
func Dir(i interface{}, c ...Config) string {
	if len(c) > 0 {
		onceCfg.Do(func() {
			cfg = c[0]
			if cfg.MaxLength < 1 {
				cfg.MaxLength = 1
			}
			if cfg.MaxLength > 10_000 {
				// the upper limit of a slice is some significant fraction of the address space of a process.
				// https://github.com/golang/go/issues/38673#issuecomment-643885108
				cfg.MaxLength = 10_000
			}

			if cfg.MaxIndentLevel < 1 {
				cfg.MaxIndentLevel = 10 // ie, the default
			}
			if cfg.MaxIndentLevel > 100 {
				cfg.MaxIndentLevel = 100
			}
		})
	}

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
			specialErr := "is not in GOROOT"
			for _, eMsg := range []string{
				// We check if it is truly a module error. We check all the errors that can be returned.
				// Unfortunately `ImportMissingError` is an internal error so we cant use errors.Is/As
				// https://github.com/golang/go/blob/go1.20.5/src/cmd/go/internal/modload/import.go#L57-L96
				//
				// This list will need to be kept uptodate with Go versions
				specialErr,
				"cannot find module providing package",
				"cannot find module",
				"is replaced but not required",
				"no required module",
				"but not at required version",
				"missing module",
			} {
				if strings.Contains(eMsg, specialErr) {
					// `package heyWorld is not in GOROOT`
					// This error message means that this is really not a module.
					// This means that someone is most probably trying to `Dir` a string variable.
					continue
				}
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

// Stackp prints to standard error the colorized stack trace.
//
// Stack trace from the runtime/stdlib is colored blue, third party libraries is yellow
// whereas your code is red.
func Stackp() {
	stackp(os.Stderr)
}

// Stack returns the colorized stack trace.
//
// Stack trace from the runtime/stdlib is colored blue, third party libraries is yellow
// whereas your code is red.
func Stack() string {
	w := &bytes.Buffer{}
	stackp(w)
	return w.String()
}

// Diffp prints a formatted diff showing the minimum line-level additions and removals that would turn old into new.
func Diffp(old, new interface{}, c ...Config) {
	fmt.Println(
		Diff(old, new, c...),
	)
}

// Diff returns a formatted diff showing the minimum line-level additions and removals that would turn old into new.
func Diff(old, new interface{}, c ...Config) string {
	return diff(
		Dir(old, c...),
		Dir(new, c...),
	)
}
