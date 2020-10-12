package main

import (
	"net/http"

	"github.com/komuw/kama"
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

func main() {
	kama.Dirp("archive/tar")
	kama.Dirp("compress/flate")
	kama.Dirp(&http.Request{})
	kama.Dirp(http.Request{})
	kama.Dirp("github.com/pkg/errors")
}
