module github.com/komuw/kama

go 1.16

require (
	github.com/frankban/quicktest v1.13.0 // test
	github.com/pkg/errors v0.9.1 // test
	github.com/sanity-io/litter v0.0.0-00010101000000-000000000000
	go.uber.org/goleak v1.1.10 // test
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/sys v0.0.0-20210514084401-e8d321eab015 // indirect
	golang.org/x/tools v0.1.1
)

// undo after https://github.com/sanity-io/litter/pull/42
replace github.com/sanity-io/litter => github.com/komuw/litter v1.5.2-0.20210519173802-98968e92f504
