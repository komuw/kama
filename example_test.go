package kama_test

import (
	"log"
	"net/http"
	"os"

	"github.com/komuw/kama"
)

type myHandler struct{ Logger *log.Logger }

func (h myHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
}

func ExampleDirp() {
	h := myHandler{Logger: log.New(os.Stderr, "", 0)}
	kama.Dirp(h)

	// Output:
	// [
	// NAME: github.com/komuw/kama_test.myHandler
	// KIND: struct
	// SIGNATURE: [kama_test.myHandler *kama_test.myHandler]
	// FIELDS: [
	// 	Logger
	// 	]
	// METHODS: [
	// 	ServeHTTP func(*kama_test.myHandler, http.ResponseWriter, *http.Request)
	// 	ServeHTTP func(kama_test.myHandler, http.ResponseWriter, *http.Request)
	// 	]
	// SNIPPET: myHandler{
	//    Logger: &Logger{},
	// }
	// ]
}
