package kama

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type Person struct {
	Name   string
	Age    int
	Height float32
	//lint:ignore U1001 used for tests
	somePrivateField string
}

//lint:ignore U1001 used for tests
func (p Person) somePrivateMethodOne() {}

//lint:ignore U1001 used for tests
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

func TestBasicVariables(t *testing.T) {
	tt := []struct {
		variable interface{}
		expected vari
	}{
		{
			Person{Name: "John"}, vari{
				Name:      "github.com/komuw/kama.Person",
				Kind:      reflect.Struct,
				Signature: []string{"kama.Person", "*kama.Person"},
				Fields:    []string{"Name", "Age", "Height"},
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
				Fields:    []string{"Name", "Age", "Height"},
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
	}

	for _, v := range tt {
		v := v
		t.Run(fmt.Sprintf("runing test for: %s", v.expected.Name), func(t *testing.T) {
			res := newVari(v.variable)

			if !cmp.Equal(res, v.expected) {
				diff := cmp.Diff(v.expected, res)
				t.Errorf("\ngot: \n\t%#+v \nwanted: \n\t%#+v \ndiff: \n======================\n%s\n======================\n", res, v.expected, diff)
			}
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
					"Method",
					"URL",
					"Proto",
					"ProtoMajor",
					"ProtoMinor",
					"Header",
					"Body",
					"GetBody",
					"ContentLength",
					"TransferEncoding",
					"Close",
					"Host",
					"Form",
					"PostForm",
					"MultipartForm",
					"Trailer",
					"RemoteAddr",
					"RequestURI",
					"TLS",
					"Cancel",
					"Response",
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
			},
		},

		{
			&http.Request{}, vari{
				Name:      "net/http.Request",
				Kind:      reflect.Struct,
				Signature: []string{"*http.Request", "http.Request"},
				Fields: []string{
					"Method",
					"URL",
					"Proto",
					"ProtoMajor",
					"ProtoMinor",
					"Header",
					"Body",
					"GetBody",
					"ContentLength",
					"TransferEncoding",
					"Close",
					"Host",
					"Form",
					"PostForm",
					"MultipartForm",
					"Trailer",
					"RemoteAddr",
					"RequestURI",
					"TLS",
					"Cancel",
					"Response",
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
			},
		},
	}

	for _, v := range tt {
		v := v
		t.Run(fmt.Sprintf("runing test for: %s", v.expected.Name), func(t *testing.T) {
			res := newVari(v.variable)
			if !cmp.Equal(res, v.expected) {
				diff := cmp.Diff(v.expected, res)
				t.Errorf("\ngot: \n\t%#+v \nwanted: \n\t%#+v \ndiff: \n======================\n%s\n======================\n", res, v.expected, diff)
			}
		})
	}
}
