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
“'Bout what?”
“You know what.”
“Oh—'course I haven't.”
“Never a word?”
“Never a solitary word, so help me. What makes you ask?”
“Well, I was afeard.”
“Why, Tom Sawyer, we wouldn't be alive two days if that got found out. You know that.”
Tom felt more comfortable. After a pause:
“Huck, they couldn't anybody get you to tell, could they?”
“Get me to tell? Why, if I wanted that halfbreed devil to drownd me they could get me to tell. They ain't no different way.”
“Well, that's all right, then. I reckon we're safe as long as we keep mum. But let's swear again, anyway. It's more surer.”
“I'm agreed.”
So they swore again with dread solemnities.
“What is the talk around, Huck? I've heard a power of it.”
“Talk? Well, it's just Muff Potter, Muff Potter, Muff Potter all the time. It keeps me in a sweat, constant, so's I want to hide som'ers.”
“That's just the same way they go on round me. I reckon he's a goner. Don't you feel sorry for him, sometimes?”
“Most always—most always. He ain't no account; but then he hain't ever done anything to hurt anybody. Just fishes a little, to get money to get drunk on—and loafs around considerable; but lord, we all do that—leastways most of us—preachers and such like. But he's kind of good—he give me half a fish, once, when there warn't enough for two; and lots of times he's kind of stood by me when I was out of luck.”
“Well, he's mended kites for me, Huck, and knitted hooks on to my line. I wish we could get him out of there.”
“My! we couldn't get him out, Tom. And besides, 'twouldn't do any good; they'd ketch him again.”
“Yes—so they would. But I hate to hear 'em abuse him so like the dickens when he never done—that.”
“I do too, Tom. Lord, I hear 'em say he's the bloodiest looking villain in this country, and they wonder he wasn't ever hung before.”
“Yes, they talk like that, all the time. I've heard 'em say that if he was to get free they'd lynch him.”
“And they'd do it, too.”
The boys had a long talk, but it brought them little comfort. As the twilight drew on, they found themselves hanging about the neighborhood of the little isolated jail, perhaps with an undefined hope that something would happen that might clear away their difficulties. But nothing happened; there seemed to be no angels or fairies interested in this luckless captive.
The boys did as they had often done before—went to the cell grating and gave Potter some tobacco and matches. He was on the ground floor and there were no guards.`

type SomeStructWIthSlice struct {
	Name      string
	MyAwesome []int
}

func TestBasicVariables(t *testing.T) {
	t.Parallel()

	tt := []struct {
		tName    string
		variable interface{}
		expected vari
	}{
		{
			tName:    "value struct",
			variable: Person{Name: "John"},
			expected: vari{
				Name:      "github.com/komuw/kama.Person",
				Kind:      reflect.Struct,
				Signature: []string{"kama.Person", "*kama.Person"},
				Fields:    []string{"Name string", "Age int", "Height float32"},
				Methods:   []string{"ValueMethodOne func(kama.Person)", "ValueMethodTwo func(kama.Person)", "PtrMethodOne func(*kama.Person)", "PtrMethodTwo func(*kama.Person) float32"},
				Val: `Person{
  Name: "John",
  Age: int(0),
  Height: float32(0),
}`,
			},
		},
		{
			tName:    "pointer struct",
			variable: &Person{Name: "Jane"},
			expected: vari{
				Name:      "github.com/komuw/kama.Person",
				Kind:      reflect.Struct,
				Signature: []string{"*kama.Person", "kama.Person"},
				Fields:    []string{"Name string", "Age int", "Height float32"},
				Methods:   []string{"ValueMethodOne func(kama.Person)", "ValueMethodTwo func(kama.Person)", "PtrMethodOne func(*kama.Person)", "PtrMethodTwo func(*kama.Person) float32"},
				Val: `&Person{
  Name: "Jane",
  Age: int(0),
  Height: float32(0),
}`,
			},
		},
		{
			tName:    "function",
			variable: ThisFunction,
			expected: vari{
				Name:      "github.com/komuw/kama.ThisFunction",
				Kind:      reflect.Func,
				Signature: []string{"func(string, int) (string, error)"},
				Fields:    []string{},
				Methods:   []string{},
				Val:       "func(string, int) (string, error)",
			},
		},
		{
			tName:    "function variable",
			variable: thisFunctionVar,
			expected: vari{
				Name:      "github.com/komuw/kama.ThisFunction",
				Kind:      reflect.Func,
				Signature: []string{"func(string, int) (string, error)"},
				Fields:    []string{},
				Methods:   []string{},
				Val:       "func(string, int) (string, error)",
			},
		},
		{
			tName:    "distinct type",
			variable: customerID(9),
			expected: vari{
				Name:      "github.com/komuw/kama.customerID",
				Kind:      reflect.Uint16,
				Signature: []string{"kama.customerID"},
				Fields:    []string{},
				Methods:   []string{"ID func(kama.customerID) uint16"},
				Val:       "kama.customerID(9)",
			},
		},
		{
			tName:    "big slice",
			variable: MyBigSlice,
			expected: vari{
				Name:      "[]int",
				Kind:      reflect.Slice,
				Signature: []string{"[]int"},
				Fields:    []string{},
				Methods:   []string{},
				Val: `[]int{
   int(0),
   int(1),
   int(2),
   int(3),
   int(4),
   int(5),
   int(6),
   int(7),
   int(8),
   int(9),
   int(10),
   int(11),
   int(12),
   int(13),
   int(14),
   int(15),
   int(16),
   int(17),
   int(18),
   int(19),
 ...<9980 more redacted>..}`,
			},
		},
		{
			tName:    "big map",
			variable: bigMap(),
			expected: vari{
				Name:      "map[int]string",
				Kind:      reflect.Map,
				Signature: []string{"map[int]string"},
				Fields:    []string{},
				Methods:   []string{},
				Val: `map[int]string{
   int(0): "0", 
   int(1): "1", 
   int(10): "10", 
   int(100): "100", 
   int(1000): "1000", 
   int(1001): "1001", 
   int(1002): "1002", 
   int(1003): "1003", 
   int(1004): "1004", 
   int(1005): "1005", 
   int(1006): "1006", 
   int(1007): "1007", 
   int(1008): "1008", 
   int(1009): "1009", 
   int(101): "101", 
   int(1010): "1010", 
   int(1011): "1011", 
   int(1012): "1012", 
   int(1013): "1013", 
   int(1014): "1014", 
   int(1015): "1015", 
   int(1016): "1016", 
   ...<9980 more redacted>..}`,
			},
		},
		{
			tName:    "big chan",
			variable: bigChan(),
			expected: vari{
				Name:      "chan int",
				Kind:      reflect.Chan,
				Signature: []string{"chan int"},
				Fields:    []string{},
				Methods:   []string{},
				Val:       "chan int (len=10000, cap=10000)",
			},
		},
		{
			tName:    "big array",
			variable: bigArray(),
			expected: vari{
				Name:      "[10000]int",
				Kind:      reflect.Array,
				Signature: []string{"[10000]int"},
				Fields:    []string{},
				Methods:   []string{},
				Val: `[10000]int{
   int(0),
   int(1),
   int(2),
   int(3),
   int(4),
   int(5),
   int(6),
   int(7),
   int(8),
   int(9),
   int(10),
   int(11),
   int(12),
   int(13),
   int(14),
   int(15),
   int(16),
   int(17),
   int(18),
   int(19),
 ...<9980 more redacted>..}`,
			},
		},
		{
			tName:    "big string",
			variable: BigString,
			expected: vari{
				Name:      "string",
				Kind:      reflect.String,
				Signature: []string{"string"},
				Fields:    []string{},
				Methods:   []string{},
				Val:       `"AT last the sleepy atmosphere was stirred—and vigorously: the murder trial came on in the court. It ...<3332 more redacted>..`,
			},
		},
		{
			tName:    "value struct with slice field",
			variable: SomeStructWIthSlice{Name: "Hello", MyAwesome: bigSlice()},
			expected: vari{
				Name:      "github.com/komuw/kama.SomeStructWIthSlice",
				Kind:      reflect.Struct,
				Signature: []string{"kama.SomeStructWIthSlice", "*kama.SomeStructWIthSlice"},
				Fields:    []string{"Name string", "MyAwesome []int"},
				Methods:   []string{},
				Val: `SomeStructWIthSlice{
  Name: "Hello",
  MyAwesome: []int{
   int(0),
   int(1),
   int(2),
   int(3),
   int(4),
   int(5),
   int(6),
   int(7),
   int(8),
   int(9),
   int(10),
   int(11),
   int(12),
   int(13),
   int(14),
   int(15),
   int(16),
   int(17),
   int(18),
   int(19),
 ...<9980 more redacted>..},
}`,
			},
		},
		{
			tName:    "pointer struct with slice field",
			variable: &SomeStructWIthSlice{Name: "HelloPointery", MyAwesome: bigSlice()},
			expected: vari{
				Name:      "github.com/komuw/kama.SomeStructWIthSlice",
				Kind:      reflect.Struct,
				Signature: []string{"*kama.SomeStructWIthSlice", "kama.SomeStructWIthSlice"},
				Fields:    []string{"Name string", "MyAwesome []int"},
				Methods:   []string{},
				Val: `&SomeStructWIthSlice{
  Name: "HelloPointery",
  MyAwesome: []int{
   int(0),
   int(1),
   int(2),
   int(3),
   int(4),
   int(5),
   int(6),
   int(7),
   int(8),
   int(9),
   int(10),
   int(11),
   int(12),
   int(13),
   int(14),
   int(15),
   int(16),
   int(17),
   int(18),
   int(19),
 ...<9980 more redacted>..},
}`,
			},
		},
	}

	for _, v := range tt {
		v := v
		t.Run(v.tName, func(t *testing.T) {
			t.Parallel()

			c := qt.New(t)
			res := newVari(v.variable)
			c.Assert(res, qt.DeepEquals, v.expected)
		})
	}
}

func TestStdlibVariables(t *testing.T) {
	t.Parallel()

	tt := []struct {
		tName    string
		variable interface{}
		expected vari
	}{
		{
			tName:    "value struct of http.Request",
			variable: http.Request{},
			expected: vari{
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
  URL: *url.URL(nil),
  Proto: "",
  ProtoMajor: int(0),
  ProtoMinor: int(0),
  Header: http.Header{(nil)},
  Body: io.ReadCloser nil,
  GetBody: func() (io.ReadCloser, error),
  ContentLength: int64(0),
  TransferEncoding: []string{(nil)},
  Close: false,
  Host: "",
  Form: url.Values{(nil)},
  PostForm: url.Values{(nil)},
  MultipartForm: *multipart.Form(nil),
  Trailer: http.Header{(nil)},
  RemoteAddr: "",
  RequestURI: "",
  TLS: *tls.ConnectionState(nil),
  Cancel: <-chan struct {} (len=0, cap=0),
  Response: *http.Response(nil),
}`,
			},
		},

		{
			tName:    "pointer struct of http.Request",
			variable: &http.Request{},
			expected: vari{
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
  URL: *url.URL(nil),
  Proto: "",
  ProtoMajor: int(0),
  ProtoMinor: int(0),
  Header: http.Header{(nil)},
  Body: io.ReadCloser nil,
  GetBody: func() (io.ReadCloser, error),
  ContentLength: int64(0),
  TransferEncoding: []string{(nil)},
  Close: false,
  Host: "",
  Form: url.Values{(nil)},
  PostForm: url.Values{(nil)},
  MultipartForm: *multipart.Form(nil),
  Trailer: http.Header{(nil)},
  RemoteAddr: "",
  RequestURI: "",
  TLS: *tls.ConnectionState(nil),
  Cancel: <-chan struct {} (len=0, cap=0),
  Response: *http.Response(nil),
}`,
			},
		},
	}

	for _, v := range tt {
		v := v
		t.Run(v.tName, func(t *testing.T) {
			t.Parallel()

			c := qt.New(t)
			res := newVari(v.variable)
			c.Assert(res, qt.DeepEquals, v.expected)
		})
	}
}

func TestSliceMap(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	var nilSlice []string = nil
	var nilMap map[string]int = nil

	tt := []struct {
		tName    string
		variable interface{}
		expected vari
	}{
		{
			tName:    "nil slice",
			variable: nilSlice,
			expected: vari{
				Name:      "[]string",
				Kind:      reflect.Slice,
				Signature: []string{"[]string"},
				Fields:    []string{},
				Methods:   []string{},
				Val:       "[]string{(nil)}",
			},
		},
		{
			tName:    "no element slice",
			variable: []string{},
			expected: vari{
				Name:      "[]string",
				Kind:      reflect.Slice,
				Signature: []string{"[]string"},
				Fields:    []string{},
				Methods:   []string{},
				Val:       "[]string{}",
			},
		},
		{
			tName:    "one element slice",
			variable: []string{"hello"},
			expected: vari{
				Name:      "[]string",
				Kind:      reflect.Slice,
				Signature: []string{"[]string"},
				Fields:    []string{},
				Methods:   []string{},
				Val: `[]string{
   "hello",
}`,
			},
		},
		{
			tName:    "two element slice",
			variable: []string{"one", "two"},
			expected: vari{
				Name:      "[]string",
				Kind:      reflect.Slice,
				Signature: []string{"[]string"},
				Fields:    []string{},
				Methods:   []string{},
				Val: `[]string{
   "one",
   "two",
}`,
			},
		},
		{
			tName:    "nil map",
			variable: nilMap,
			expected: vari{
				Name:      "map[string]int",
				Kind:      reflect.Map,
				Signature: []string{"map[string]int"},
				Fields:    []string{},
				Methods:   []string{},
				Val:       "map[string]int{(nil)}",
			},
		},
		{
			tName:    "no element map",
			variable: map[string]int{},
			expected: vari{
				Name:      "map[string]int",
				Kind:      reflect.Map,
				Signature: []string{"map[string]int"},
				Fields:    []string{},
				Methods:   []string{},
				Val:       "map[string]int{}",
			},
		},
		{
			tName:    "one element map",
			variable: map[string]int{"o": 1},
			expected: vari{
				Name:      "map[string]int",
				Kind:      reflect.Map,
				Signature: []string{"map[string]int"},
				Fields:    []string{},
				Methods:   []string{},
				Val: `map[string]int{
   "o": int(1), }`,
			},
		},
		{
			tName:    "two element map",
			variable: map[string]int{"o": 1, "two": 2},
			expected: vari{
				Name:      "map[string]int",
				Kind:      reflect.Map,
				Signature: []string{"map[string]int"},
				Fields:    []string{},
				Methods:   []string{},
				Val: `map[string]int{
   "o": int(1), 
   "two": int(2), }`,
			},
		},
	}

	for _, v := range tt {
		v := v

		t.Run(v.tName, func(t *testing.T) {
			t.Parallel()

			res := newVari(v.variable)
			c.Assert(res, qt.DeepEquals, v.expected)
		})
	}
}
