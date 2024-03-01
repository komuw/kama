module github.com/komuw/kama

go 1.20

require (
	golang.org/x/exp v0.0.0-20240222234643-814bf88cf225
	golang.org/x/tools v0.18.0
)

require (
	github.com/google/go-cmp v0.6.0 // indirect
	golang.org/x/mod v0.15.0 // indirect
)

require (
	github.com/pkg/errors v0.9.1 // test
	go.akshayshah.org/attest v1.0.2 // test
	go.uber.org/goleak v1.3.0 // test
)

retract v0.0.14 // Contains a bug.
