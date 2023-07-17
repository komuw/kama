## kama          

[![Go Reference](https://pkg.go.dev/badge/github.com/komuw/kama.svg)](https://pkg.go.dev/github.com/komuw/kama)
[![ci](https://github.com/komuw/kama/workflows/kama%20ci/badge.svg)](https://github.com/komuw/kama/actions)
[![codecov](https://codecov.io/gh/komuw/kama/branch/main/graph/badge.svg)](https://codecov.io/gh/komuw/kama)


`kama` prints exported information of types, variables, packages, modules, imports etc     
It also pretty prints data structures.    
It can be used to aid debugging and testing.        
If you have heard of [kr/pretty](https://github.com/kr/pretty), [sanity-io/litter](https://github.com/sanity-io/litter), [davecgh/go-spew](https://github.com/davecgh/go-spew) etc; then `kama` is like those except that it;   
(a) prints the exported API of types, modules etc     
and     
(b) pretty prints data structures.         

It is heavily inspired by Python's [`dir`](https://docs.python.org/3/library/functions.html#dir) builtin function.       

It's name is derived from Kenyan hip hop artiste, `Kama`(One third of the hiphop group `Kalamashaka`).                               


## Installation

```shell
go get -u github.com/komuw/kama
```   

## Usage:    

#### (a) print exported api of modules
```go
import "github.com/komuw/kama"

kama.Dirp("compress/flate")
kama.Dirp("github.com/pkg/errors")
```
that will print:
```bash
[
NAME: compress/flate
CONSTANTS: [
	BestCompression untyped int 
	BestSpeed untyped int 
	DefaultCompression untyped int 
	HuffmanOnly untyped int 
	NoCompression untyped int 
	]
VARIABLES: []
FUNCTIONS: [
	NewReader(r io.Reader) io.ReadCloser 
	NewReaderDict(r io.Reader, dict []byte) io.ReadCloser 
	NewWriter(w io.Writer, level int) (*Writer, error) 
	NewWriterDict(w io.Writer, level int, dict []byte) (*Writer, error) 
	]
TYPES: [
	CorruptInputError int64
		(CorruptInputError) Error() string 
	InternalError string
		(InternalError) Error() string 
	ReadError struct
		(*ReadError) Error() string 
	Reader interface
		(Reader) Read(p []byte) (n int, err error)
		(Reader) ReadByte() (byte, error) 
	Resetter interface
		(Resetter) Reset(r io.Reader, dict []byte) error 
	WriteError struct
		(*WriteError) Error() string 
	Writer struct
		(*Writer) Close() error
		(*Writer) Flush() error
		(*Writer) Reset(dst io.Writer)
		(*Writer) Write(data []byte) (n int, err error)]
]
```
```bash
[
NAME: github.com/pkg/errors
CONSTANTS: []
VARIABLES: []
FUNCTIONS: [
	As(err error, target interface{}) bool 
	Cause(err error) error 
	Errorf(format string, args ...interface{}) error 
	Is(err error, target error) bool 
	New(message string) error 
	Unwrap(err error) error 
	WithMessage(err error, message string) error 
	WithMessagef(err error, format string, args ...interface{}) error 
	WithStack(err error) error 
	Wrap(err error, message string) error 
	Wrapf(err error, format string, args ...interface{}) error 
	]
TYPES: [
	Frame uintptr
		(Frame) Format(s fmt.State, verb rune)
		(Frame) MarshalText() ([]byte, error) 
	StackTrace []Frame
		(StackTrace) Format(s fmt.State, verb rune)]
]
```

#### (b) pretty data structures:
```go
req, _ := http.NewRequest("GET", "https://example.com", nil)
req.Header.Set("Content-Type", "application/octet-stream")
req.AddCookie(&http.Cookie{Name: "hello", Value: "world"})

kama.Dirp(req)
```
that will print:
```bash
[
NAME: net/http.Request
KIND: struct
SIGNATURE: [*http.Request http.Request]
FIELDS: [
	Method string 
	URL *url.URL 
	Proto string 
	ProtoMajor int 
	ProtoMinor int 
	Header http.Header 
	Body io.ReadCloser 
	GetBody func() (io.ReadCloser, error) 
	ContentLength int64 
	TransferEncoding []string 
	Close bool 
	Host string 
	Form url.Values 
	PostForm url.Values 
	MultipartForm *multipart.Form 
	Trailer http.Header 
	RemoteAddr string 
	RequestURI string 
	TLS *tls.ConnectionState 
	Cancel <-chan struct {} 
	Response *http.Response 
	]
METHODS: [
	AddCookie func(*http.Request, *http.Cookie) 
	BasicAuth func(*http.Request) (string, string, bool) 
	Clone func(*http.Request, context.Context) *http.Request 
	Context func(*http.Request) context.Context 
	Cookie func(*http.Request, string) (*http.Cookie, error) 
	Cookies func(*http.Request) []*http.Cookie 
	FormFile func(*http.Request, string) (multipart.File, *multipart.FileHeader, error) 
	FormValue func(*http.Request, string) string 
	MultipartReader func(*http.Request) (*multipart.Reader, error) 
	ParseForm func(*http.Request) error 
	ParseMultipartForm func(*http.Request, int64) error 
	PostFormValue func(*http.Request, string) string 
	ProtoAtLeast func(*http.Request, int, int) bool 
	Referer func(*http.Request) string 
	SetBasicAuth func(*http.Request, string, string) 
	UserAgent func(*http.Request) string 
	WithContext func(*http.Request, context.Context) *http.Request 
	Write func(*http.Request, io.Writer) error 
	WriteProxy func(*http.Request, io.Writer) error 
	]
SNIPPET: &Request{
  Method: "GET",
  URL: &URL{
    Scheme: "https",
    Host: "example.com",
  },
  Proto: "HTTP/1.1",
  ProtoMajor: int(1),
  ProtoMinor: int(1),
  Header: http.Header{
   "Content-Type": []string{
   "application/octet-stream",
      }, 
   "Cookie": []string{
   "hello=world",
      }, 
    },
  Body: io.ReadCloser nil,
  GetBody: func() (io.ReadCloser, error),
  ContentLength: int64(0),
  TransferEncoding: []string{(nil)},
  Close: false,
  Host: "example.com",
  Form: url.Values{(nil)},
  PostForm: url.Values{(nil)},
  MultipartForm: *multipart.Form(nil),
  Trailer: http.Header{(nil)},
  RemoteAddr: "",
  RequestURI: "",
  TLS: *tls.ConnectionState(nil),
  Cancel: <-chan struct {} (len=0, cap=0),
  Response: *http.Response(nil),
}
]
```
See [testdata](testdata) directory for more examples.    

## Testing
```shell
# run tests:
export KAMA_WRITE_DATA_FOR_TESTS=YES
unset KAMA_WRITE_DATA_FOR_TESTS
go test -race ./... -count=1
```

## Inspiration
1. Python's [`dir`](https://docs.python.org/3/library/functions.html#dir) builtin function.    
2. [`godex`](https://pkg.go.dev/golang.org/x/tools/cmd/godex).   
3. [`sanity-io/litter`](https://github.com/sanity-io/litter).

## Prior art
1. https://github.com/kr/pretty
2. https://github.com/sanity-io/litter
3. https://github.com/davecgh/go-spew
4. https://github.com/hexops/valast
5. https://github.com/alecthomas/repr
6. https://github.com/k0kubun/pp
7. https://github.com/jba/printsrc
