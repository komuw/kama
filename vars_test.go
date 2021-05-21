package kama

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	qt "github.com/frankban/quicktest"
)

type Person struct {
	Name   string
	Age    int
	Height float32
	//lint:ignore U1001,U1000 used for tests
	somePrivateField string
}

//lint:ignore U1001,U1000 used for tests
func (p Person) somePrivateMethodOne() {}

//lint:ignore U1001,U1000 used for tests
func (p *Person) somePrivateMethodTwo() {}

//lint:ignore U1001 used for tests
func (p Person) ValueMethodOne() {}

//lint:ignore U1001 used for tests
func (p *Person) PtrMethodOne() {}

//lint:ignore U1001 used for tests
func (p Person) ValueMethodTwo() {}

//lint:ignore U1001 used for tests
func (p *Person) PtrMethodTwo() float32 { return p.Height }

func ThisFunction(arg1 string, arg2 int) (string, error) {
	return "", nil
}

var thisFunctionVar = ThisFunction

type customerID uint16

//lint:ignore U1001 used for tests
func (c customerID) ID() uint16 { return uint16(c) }

func getArray() [10_000]int {
	a := [10_000]int{}
	for i := 0; i < 10_000; i++ {
		a[i] = i
	}
	return a
}

func getChan() chan int {
	z := make(chan int, 10_000)
	for i := 0; i < 10_000; i++ {
		z <- i
	}
	return z
}

func bigSlice() []int {
	x := []int{}
	for i := 0; i < 10_000; i++ {
		x = append(x, i)
	}
	return x
}

var MyBigSlice = bigSlice()

func sliceOfStruct() []http.Request {
	xx := []http.Request{}
	for i := 0; i < 10_000; i++ {
		xx = append(xx, http.Request{Method: fmt.Sprintf("%d", i)})
	}
	return xx
}

func TestBasicVariables(t *testing.T) {
	tt := []struct {
		variable interface{}
		expected vari
	}{

		{
			getChan(), vari{
				Name:      "chan",
				Kind:      reflect.Chan,
				Signature: []string{"chan int"},
				Fields:    []string{},
				Methods:   []string{},
				Val:       "chan int"},
		},

		{
			getArray(), vari{
				Name:      "array",
				Kind:      reflect.Array,
				Signature: []string{"[10000]int"},
				Fields:    []string{},
				Methods:   []string{},
				Val:       "[10000]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,1 ...<snipped>.."},
		},

		{
			Person{Name: "John"}, vari{
				Name:      "github.com/komuw/kama.Person",
				Kind:      reflect.Struct,
				Signature: []string{"kama.Person", "*kama.Person"},
				Fields:    []string{"Name string", "Age int", "Height float32"},
				Methods:   []string{"ValueMethodOne func(kama.Person)", "ValueMethodTwo func(kama.Person)", "PtrMethodOne func(*kama.Person)", "PtrMethodTwo func(*kama.Person) float32"},
				Val: `Person{
  Name: "John",
  Age: 0,
  Height: 0,
}`,
			},
		},
		{

			&Person{Name: "Jane"}, vari{
				Name:      "github.com/komuw/kama.Person",
				Kind:      reflect.Struct,
				Signature: []string{"*kama.Person", "kama.Person"},
				Fields:    []string{"Name string", "Age int", "Height float32"},
				Methods:   []string{"ValueMethodOne func(kama.Person)", "ValueMethodTwo func(kama.Person)", "PtrMethodOne func(*kama.Person)", "PtrMethodTwo func(*kama.Person) float32"},
				Val: `&Person{
  Name: "Jane",
  Age: 0,
  Height: 0,
}`,
			},
		},
		{
			ThisFunction, vari{
				Name:      "github.com/komuw/kama.ThisFunction",
				Kind:      reflect.Func,
				Signature: []string{"func(string, int) (string, error)"},
				Fields:    []string{},
				Methods:   []string{},
				Val:       "ThisFunction",
			},
		},
		{
			thisFunctionVar, vari{
				Name:      "github.com/komuw/kama.ThisFunction",
				Kind:      reflect.Func,
				Signature: []string{"func(string, int) (string, error)"},
				Fields:    []string{},
				Methods:   []string{},
				Val:       "ThisFunction",
			},
		},
		{
			customerID(9), vari{
				Name:      "github.com/komuw/kama.customerID",
				Kind:      reflect.Uint16,
				Signature: []string{"kama.customerID"},
				Fields:    []string{},
				Methods:   []string{"ID func(kama.customerID) uint16"},
				Val:       "9",
			},
		},
		{
			MyBigSlice, vari{
				Name:      "slice",
				Kind:      reflect.Slice,
				Signature: []string{"[]int"},
				Fields:    []string{},
				Methods:   []string{},
				Val:       "[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17, ...<snipped>..",
			},
		},
		{
			sliceOfStruct(), vari{
				// TODO: fix this name
				Name:      "slice",
				Kind:      reflect.Slice,
				Signature: []string{"[]http.Request"},
				Fields:    []string{},
				Methods:   []string{},
				Val:       `[]Request{Request{Method:"0",URL:nil,Proto:"",Prot ...<snipped>..`,
			},
		},
	}

	for _, v := range tt {
		v := v
		t.Run(fmt.Sprintf("runing test for: %s", v.expected.Name), func(t *testing.T) {
			c := qt.New(t)
			res := newVari(v.variable)
			c.Assert(res, qt.DeepEquals, v.expected)
		})
	}
}

func TestStdlibVariables(t *testing.T) {
	tt := []struct {
		variable interface{}
		expected vari
	}{
		{
			http.Request{}, vari{
				Name:      "net/http.Request",
				Kind:      reflect.Struct,
				Signature: []string{"http.Request", "*http.Request"},
				Fields: []string{
					"Method string",
					"URL *url.URL",
					"Proto string",
					"ProtoMajor int",
					"ProtoMinor int",
					"Header http.Header",
					"Body io.ReadCloser",
					"GetBody func() (io.ReadCloser, error)",
					"ContentLength int64",
					"TransferEncoding []string",
					"Close bool",
					"Host string",
					"Form url.Values",
					"PostForm url.Values",
					"MultipartForm *multipart.Form",
					"Trailer http.Header",
					"RemoteAddr string",
					"RequestURI string",
					"TLS *tls.ConnectionState",
					"Cancel <-chan struct {}",
					"Response *http.Response",
				},
				Methods: []string{
					"AddCookie func(*http.Request, *http.Cookie)",
					"BasicAuth func(*http.Request) (string, string, bool)",
					"Clone func(*http.Request, context.Context) *http.Request",
					"Context func(*http.Request) context.Context",
					"Cookie func(*http.Request, string) (*http.Cookie, error)",
					"Cookies func(*http.Request) []*http.Cookie",
					"FormFile func(*http.Request, string) (multipart.File, *multipart.FileHeader, error)",
					"FormValue func(*http.Request, string) string",
					"MultipartReader func(*http.Request) (*multipart.Reader, error)",
					"ParseForm func(*http.Request) error",
					"ParseMultipartForm func(*http.Request, int64) error",
					"PostFormValue func(*http.Request, string) string",
					"ProtoAtLeast func(*http.Request, int, int) bool",
					"Referer func(*http.Request) string",
					"SetBasicAuth func(*http.Request, string, string)",
					"UserAgent func(*http.Request) string",
					"WithContext func(*http.Request, context.Context) *http.Request",
					"Write func(*http.Request, io.Writer) error",
					"WriteProxy func(*http.Request, io.Writer) error",
				},
				Val: `Request{
  Method: "",
  URL: nil,
  Proto: "",
  ProtoMajor: 0,
  ProtoMinor: 0,
  Header: Header(nil),
  Body: nil,
  GetBody: ,
  ContentLength: 0,
  TransferEncoding: nil,
  Close: false,
  Host: "",
  Form: Values(nil),
  PostForm: Values(nil),
  MultipartForm: nil,
  Trailer: Header(nil),
  RemoteAddr: "",
  RequestURI: "",
  TLS: nil,
  Cancel: <-chan struct {},
  Response: nil,
}`,
			},
		},

		{
			&http.Request{}, vari{
				Name:      "net/http.Request",
				Kind:      reflect.Struct,
				Signature: []string{"*http.Request", "http.Request"},
				Fields: []string{
					"Method string",
					"URL *url.URL",
					"Proto string",
					"ProtoMajor int",
					"ProtoMinor int",
					"Header http.Header",
					"Body io.ReadCloser",
					"GetBody func() (io.ReadCloser, error)",
					"ContentLength int64",
					"TransferEncoding []string",
					"Close bool",
					"Host string",
					"Form url.Values",
					"PostForm url.Values",
					"MultipartForm *multipart.Form",
					"Trailer http.Header",
					"RemoteAddr string",
					"RequestURI string",
					"TLS *tls.ConnectionState",
					"Cancel <-chan struct {}",
					"Response *http.Response",
				},
				Methods: []string{
					"AddCookie func(*http.Request, *http.Cookie)",
					"BasicAuth func(*http.Request) (string, string, bool)",
					"Clone func(*http.Request, context.Context) *http.Request",
					"Context func(*http.Request) context.Context",
					"Cookie func(*http.Request, string) (*http.Cookie, error)",
					"Cookies func(*http.Request) []*http.Cookie",
					"FormFile func(*http.Request, string) (multipart.File, *multipart.FileHeader, error)",
					"FormValue func(*http.Request, string) string",
					"MultipartReader func(*http.Request) (*multipart.Reader, error)",
					"ParseForm func(*http.Request) error",
					"ParseMultipartForm func(*http.Request, int64) error",
					"PostFormValue func(*http.Request, string) string",
					"ProtoAtLeast func(*http.Request, int, int) bool",
					"Referer func(*http.Request) string",
					"SetBasicAuth func(*http.Request, string, string)",
					"UserAgent func(*http.Request) string",
					"WithContext func(*http.Request, context.Context) *http.Request",
					"Write func(*http.Request, io.Writer) error",
					"WriteProxy func(*http.Request, io.Writer) error",
				},
				Val: `&Request{
  Method: "",
  URL: nil,
  Proto: "",
  ProtoMajor: 0,
  ProtoMinor: 0,
  Header: Header(nil),
  Body: nil,
  GetBody: ,
  ContentLength: 0,
  TransferEncoding: nil,
  Close: false,
  Host: "",
  Form: Values(nil),
  PostForm: Values(nil),
  MultipartForm: nil,
  Trailer: Header(nil),
  RemoteAddr: "",
  RequestURI: "",
  TLS: nil,
  Cancel: <-chan struct {},
  Response: nil,
}`,
			},
		},
	}

	for _, v := range tt {
		v := v
		t.Run(fmt.Sprintf("runing test for: %s", v.expected.Name), func(t *testing.T) {
			c := qt.New(t)

			res := newVari(v.variable)
			c.Assert(res, qt.DeepEquals, v.expected)
		})
	}
}
