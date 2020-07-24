package main

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

/*
  go test -timeout 1m -race -cover -v ./...

  This tests are from/inspired by;
  https://github.com/sanity-io/litter/blob/b3546bd0a12c8738436e565b9e016bcd1876403d/dump_test.go

  That repo as of that commit has an MIT license;
  https://github.com/sanity-io/litter/blob/b3546bd0a12c8738436e565b9e016bcd1876403d/LICENSE
*/

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

func TestMore(t *testing.T) {
	tt := []interface{}{
		reflect.Value{},
		errors.New,
		http.Handle,
		http.HandleFunc,
	}

	for _, v := range tt {
		dir(v)

	}
}
