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

func bigMap() map[int]string {
	y := map[int]string{}
	for i := 0; i < 10_000; i++ {
		y[i] = fmt.Sprintf("%d", i)
	}
	return y
}

func bigChan() chan int {
	z := make(chan int, 10_000)
	for i := 0; i < 10_000; i++ {
		// TODO: will be fixed by https://github.com/sanity-io/litter/pull/42
		z <- i
	}
	return z
}

func bigArray() [10_000]int {
	a := [10_000]int{}
	for i := 0; i < 10_000; i++ {
		a[i] = i
	}
	return a
}

var BigString = `AT last the sleepy atmosphere was stirred—and vigorously: the murder trial came on in the court. It became the absorbing topic of village talk immediately. Tom could not get away from it. Every reference to the murder sent a shudder to his heart, for his troubled conscience and fears almost persuaded him that these remarks were put forth in his hearing as “feelers”; he did not see how he could be suspected of knowing anything about the murder, but still he could not be comfortable in the midst of this gossip. It kept him in a cold shiver all the time. He took Huck to a lonely place to have a talk with him. It would be some relief to unseal his tongue for a little while; to divide his burden of distress with another sufferer. Moreover, he wanted to assure himself that Huck had remained discreet.
“Huck, have you ever told anybody about—that?”
“’Bout what?”
“You know what.”
“Oh—’course I haven’t.”
“Never a word?”
“Never a solitary word, so help me. What makes you ask?”
“Well, I was afeard.”
“Why, Tom Sawyer, we wouldn’t be alive two days if that got found out. You know that.”
Tom felt more comfortable. After a pause:
“Huck, they couldn’t anybody get you to tell, could they?”
“Get me to tell? Why, if I wanted that halfbreed devil to drownd me they could get me to tell. They ain’t no different way.”
“Well, that’s all right, then. I reckon we’re safe as long as we keep mum. But let’s swear again, anyway. It’s more surer.”
“I’m agreed.”
So they swore again with dread solemnities.
“What is the talk around, Huck? I’ve heard a power of it.”
“Talk? Well, it’s just Muff Potter, Muff Potter, Muff Potter all the time. It keeps me in a sweat, constant, so’s I want to hide som’ers.”
“That’s just the same way they go on round me. I reckon he’s a goner. Don’t you feel sorry for him, sometimes?”
“Most always—most always. He ain’t no account; but then he hain’t ever done anything to hurt anybody. Just fishes a little, to get money to get drunk on—and loafs around considerable; but lord, we all do that—leastways most of us—preachers and such like. But he’s kind of good—he give me half a fish, once, when there warn’t enough for two; and lots of times he’s kind of stood by me when I was out of luck.”
“Well, he’s mended kites for me, Huck, and knitted hooks on to my line. I wish we could get him out of there.”
“My! we couldn’t get him out, Tom. And besides, ’twouldn’t do any good; they’d ketch him again.”
“Yes—so they would. But I hate to hear ’em abuse him so like the dickens when he never done—that.”
“I do too, Tom. Lord, I hear ’em say he’s the bloodiest looking villain in this country, and they wonder he wasn’t ever hung before.”
“Yes, they talk like that, all the time. I’ve heard ’em say that if he was to get free they’d lynch him.”
“And they’d do it, too.”
The boys had a long talk, but it brought them little comfort. As the twilight drew on, they found themselves hanging about the neighborhood of the little isolated jail, perhaps with an undefined hope that something would happen that might clear away their difficulties. But nothing happened; there seemed to be no angels or fairies interested in this luckless captive.
The boys did as they had often done before—went to the cell grating and gave Potter some tobacco and matches. He was on the ground floor and there were no guards.`

type SomeStructWIthSlice struct {
	Name      string
	MyAwesome []int
}

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
				Name:      "[]int",
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
				Name:      "[]http.Request",
				Kind:      reflect.Slice,
				Signature: []string{"[]http.Request"},
				Fields:    []string{},
				Methods:   []string{},
				Val:       `[]Request{Request{Method:"0",URL:nil,Proto:"",Prot ...<snipped>..`,
			},
		},
		{
			bigMap(), vari{
				Name:      "map[int]string",
				Kind:      reflect.Map,
				Signature: []string{"map[int]string"},
				Fields:    []string{},
				Methods:   []string{},
				Val:       `map[int]string{0:"0",1:"1",10:"10",100:"100",1000: ...<snipped>..`,
			},
		},
		{
			bigChan(), vari{
				Name:      "chan int",
				Kind:      reflect.Chan,
				Signature: []string{"chan int"},
				Fields:    []string{},
				Methods:   []string{},
				Val:       "chan int",
			},
		},
		{
			bigArray(), vari{
				Name:      "[10000]int",
				Kind:      reflect.Array,
				Signature: []string{"[10000]int"},
				Fields:    []string{},
				Methods:   []string{},
				Val:       "[10000]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,1 ...<snipped>.."},
		},
		{
			BigString, vari{
				Name:      "string",
				Kind:      reflect.String,
				Signature: []string{"string"},
				Fields:    []string{},
				Methods:   []string{},
				Val:       `"AT last the sleepy atmosphere was stirred—and v ...<snipped>..`},
		},
		{

			SomeStructWIthSlice{Name: "Hello", MyAwesome: bigSlice()}, vari{
				Name:      "github.com/komuw/kama.SomeStructWIthSlice",
				Kind:      reflect.Struct,
				Signature: []string{"kama.SomeStructWIthSlice", "*kama.SomeStructWIthSlice"},
				Fields:    []string{"Name string", "MyAwesome []int"},
				Methods:   []string{},
				Val: `SomeStructWIthSlice{
  Name: "Hello",
  MyAwesome: []int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17, ...<snipped>..,
}`,
			},
		},
		{

			&SomeStructWIthSlice{Name: "HelloPointery", MyAwesome: bigSlice()}, vari{
				Name:      "github.com/komuw/kama.SomeStructWIthSlice",
				Kind:      reflect.Struct,
				Signature: []string{"*kama.SomeStructWIthSlice", "kama.SomeStructWIthSlice"},
				Fields:    []string{"Name string", "MyAwesome []int"},
				Methods:   []string{},
				Val: `&SomeStructWIthSlice{
  Name: "HelloPointery",
  MyAwesome: []int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17, ...<snipped>..,
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
