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

// TODO: maybe add syntax highlighting, maybe make it optional??

// TODO: clean up

// TODO: add of `dir` documentation

// TODO: maybe we should show docs when someone requests for something specific.
// eg if they do `dir(http)` we do not show docs, but if they do `dir(&http.Request{})` we show docs.
// An alternative is only show docs, if someone requests. `dir(i interface{}, config ...dir.Config)`; config is `...` so that it is optional
// where config is a `type Config struct {}`

// TODO: add a command line api.
//   eg; `dir http.Request` or `dir http`
// have a look at `golang.org/x/tools/cmd/godex`

// TODO: this will stutter; `dir.dir(23)`
// maybe it is okay??
// TODO: surface all info for both the type and its pointer.
// currently `dir(&http.Client{})` & `dir(http.Client{})` produces different output; they should NOT
func dir(i interface{}) {
	var res interface{}
	var err error

	if reflect.TypeOf(i).Kind() == reflect.String {
		i := i.(string)
		res, err = newPaki(i)
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
