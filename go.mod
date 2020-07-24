module github.com/komuw/dir

go 1.14

require (
	github.com/bradfitz/iter v0.0.0-20191230175014-e8f45d346db8
	github.com/containerd/containerd v1.3.6 // indirect
	github.com/komuw/meli v0.2.2
	github.com/pkg/errors v0.8.1
	github.com/sanity-io/litter v1.2.0
)

replace github.com/docker/docker => github.com/docker/engine v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
