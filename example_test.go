package kama_test

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/komuw/kama"
)

type myHandler struct{ Logger *log.Logger }

func (h myHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
}

// func ExampleDir() {
// 	h := myHandler{Logger: log.New(os.Stderr, "", 0)}
// 	kama.Dir(h)
// 	// kama.Dir("compress/flate")
// 	// kama.Dir(&http.Request{})
// 	// kama.Dir(http.Request{})
// 	// kama.Dir("github.com/pkg/errors")

// 	// Output:
// 	//[
// 	//NAME: github.com/komuw/kama_test.myHandler
// 	//KIND: struct
// 	//SIGNATURE: [kama_test.myHandler *kama_test.myHandler]
// 	//FIELDS: [
// 	//	Logger
// 	//	]
// 	//METHODS: []
// 	//]
// }

func TestVars(t *testing.T) {
	tt := []struct {
		variable interface{}
		expected string
	}{
		{
			myHandler{Logger: log.New(os.Stderr, "", 0)}, `
[
NAME: github.com/komuw/kama_test.myHandler
KIND: struct
SIGNATURE: [kama_test.myHandler *kama_test.myHandler]
FIELDS: [
	Logger 
	]
METHODS: [
	ServeHTTP func(*kama_test.myHandler, http.ResponseWriter, *http.Request) 
	ServeHTTP func(kama_test.myHandler, http.ResponseWriter, *http.Request) 
	]
]
`,
		},
	}

	for _, v := range tt {
		res := kama.Dir(v.variable)

		if !cmp.Equal(res, v.expected) {
			t.Errorf("\ngot \n\t%#+v \nwanted \n\t%#+v", res, v.expected)
		}
	}
}
