package kama_test

// import (
// 	"log"
// 	"net/http"
// )

// type myHandler struct{ Logger *log.Logger }

// func (h myHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
// }

// This is blocked on:
// 1. https://github.com/golang/go/issues/41980
// 2. https://github.com/golang/go/issues/5128#issuecomment-708940093
// func ExampleDirp() {
// 	h := myHandler{Logger: log.New(os.Stderr, "", 0)}
// 	kama.Dirp(h)

// 	// Output:
// 	// [
// 	// NAME: github.com/komuw/kama_test.myHandler
// 	// KIND: struct
// 	// SIGNATURE: [kama_test.myHandler *kama_test.myHandler]
// 	// FIELDS: [
// 	// 	Logger
// 	// 	]
// 	// METHODS: [
// 	// 	ServeHTTP func(*kama_test.myHandler, http.ResponseWriter, *http.Request)
// 	// 	ServeHTTP func(kama_test.myHandler, http.ResponseWriter, *http.Request)
// 	// 	]
// 	// SNIPPET: myHandler{
// 	//    Logger: &Logger{},
// 	// }
// 	// ]
// }
