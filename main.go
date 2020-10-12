package main

import (
	"fmt"
	"net/http"

	"reflect"
)

// TODO: If someone passes in, say a struct;
// we should show them its type, methods etc
// but also print it out and its contents
// basically, do what `litter.Dump` would have done

// TODO: clean up

// TODO: fuzz test

// TODO: add documentation for `kama`

// TODO: add a command line api.
//   eg; `kama http.Request` or `kama http`
// have a look at `golang.org/x/tools/cmd/godex`

func dir(i interface{}) {
	var res interface{}
	var err error

	iType := reflect.TypeOf(i)
	if iType == nil {
		res = newVari(i)
	} else if iType.Kind() == reflect.String {
		i := i.(string)
		res, err = newPak(i)
		if err != nil {
			res = newVari(i)
		}
	} else {
		res = newVari(i)
	}

	fmt.Println(res)

}

func main() {
	defer panicHandler()

	dir("archive/tar")
	dir("compress/flate")
	dir(&http.Request{})
	dir(http.Request{})
	dir("github.com/pkg/errors")
}
