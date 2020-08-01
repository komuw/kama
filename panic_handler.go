package main

import "fmt"

// Usage:
//
// func main() {
// 	defer panicHandler()
// }
//
// Do note that this would not catch all panics.
// eg; `panic(nil)` will not be caught
//
func panicHandler() {
	// keep an eye on the accepeted proposal: issues/37023
	// https://github.com/golang/go/issues/37023

	r := recover()
	if r != nil {
		fmt.Println(`
		\n
		panicHandler recovered.
		This is where you call your logging service.
		And or your metric service to tke error counts.
		\n\n
	`)
		panic(r)
	}
}
