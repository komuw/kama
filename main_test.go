package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"

	pkgErrors "github.com/pkg/errors"
)

/*
  go test -timeout 1m -race -cover -v ./...

  This tests are from/inspired by;
  https://github.com/sanity-io/litter/blob/b3546bd0a12c8738436e565b9e016bcd1876403d/dump_test.go

  That repo as of that commit has an MIT license;
  https://github.com/sanity-io/litter/blob/b3546bd0a12c8738436e565b9e016bcd1876403d/LICENSE
*/

const acceptableCodeCoverage = 0.5

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags

	eCode := m.Run()

	if eCode == 0 && testing.CoverMode() != "" {
		coverage := testing.Coverage()
		// note: for some reason the value of `coverage` is always less
		// than the one reported on the terminal by go test

		if coverage < acceptableCodeCoverage {
			fmt.Printf("\n\tThe test code coverage has fallen below the acceptable value of %v. The current value is %v. \n", acceptableCodeCoverage, coverage)
			eCode = -1
		}
	}

	os.Exit(eCode)
}

type BlankStruct struct{}
type BasicStruct struct {
	Public  int
	private int
}
type IntAlias int

func SomeFunction(arg1 string, arg2 int) (string, error) {
	return "", nil
}

func OkayFunc(arg1 string, arg2 int) {

}

func TestPrimitives(t *testing.T) {
	tt := []interface{}{
		false,
		true,
		7,
		12.3535,
		int8(10),
		int16(10),
		int32(10),
		int64(10),
		uint8(10),
		uint16(10),
		uint32(10),
		uint64(10),
		uint(10),
		float32(12.3),
		float64(12.3),
		complex64(12 + 10.5i),
		complex128(-1.2 - 0.1i),
		(func(v int) *int { return &v })(10),
		"string with \"quote\"",
		[]int{1, 2, 3},
		interface{}("hello from interface"),
		BlankStruct{},
		&BlankStruct{},
		BasicStruct{1, 2},
		IntAlias(10),
		(func(v IntAlias) *IntAlias { return &v })(10),
		SomeFunction,

		func(arg string) (bool, error) { return false, nil },
		nil,
		interface{}(nil),
	}

	for _, v := range tt {
		dir(v)

	}
}

type myHandler struct{}

func (h myHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
}

func TestStdlibTypes(t *testing.T) {
	tt := []interface{}{
		errors.New,
		reflect.Value{},
		http.Handle,
		http.HandleFunc,
		http.Handler(myHandler{}),
	}

	for _, v := range tt {
		dir(v)

	}
}

func TestThirdPartyTypes(t *testing.T) {
	tt := []interface{}{
		pkgErrors.New,
	}

	for _, v := range tt {
		dir(v)
	}
}
