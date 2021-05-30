package e2e_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"unsafe"

	qt "github.com/frankban/quicktest"
	"github.com/komuw/kama"
)

var longText = `AT last the sleepy atmosphere was stirred—and vigorously: the murder trial came on in the court. It became the absorbing topic of village talk immediately. Tom could not get away from it. Every reference to the murder sent a shudder to his heart, for his troubled conscience and fears almost persuaded him that these remarks were put forth in his hearing as “feelers”; he did not see how he could be suspected of knowing anything about the murder, but still he could not be comfortable in the midst of this gossip. It kept him in a cold shiver all the time. He took Huck to a lonely place to have a talk with him. It would be some relief to unseal his tongue for a little while; to divide his burden of distress with another sufferer. Moreover, he wanted to assure himself that Huck had remained discreet.
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

type Distance uint64

func bigMap() map[int]string {
	y := map[int]string{}
	for i := 0; i < 10_000; i++ {
		y[i] = fmt.Sprintf("%d", i)
	}
	return y
}

func bigChan() chan int {
	z := make(chan int, 10_000)
	for i := 0; i < 122; i++ {
		// TODO: will be fixed by https://github.com/sanity-io/litter/pull/42
		z <- i
	}
	return z
}

func makeDirecteChan() chan<- bool {
	directedChan := make(chan<- bool, 13)
	directedChan <- true
	return directedChan
}

func bigSlice() []int {
	bigSlice := []int{}
	for i := 0; i < 10_000; i++ {
		bigSlice = append(bigSlice, i)
	}
	return bigSlice
}

func makeSliceOfHttpRequests() []http.Request {
	h := []http.Request{}
	for i := 0; i < 100; i++ {
		h = append(h, http.Request{Method: fmt.Sprint(i)})
	}
	return h
}

// We just need a type that will implement the `io.ReadCloser` interface
type CustomReadCloser int64

func (c CustomReadCloser) Read(p []byte) (n int, err error) {
	return 10, nil
}
func (c CustomReadCloser) Close() error {
	return errors.New("CustomReadCloser always fails closing")
}

type FuncWithReturn func(http.ResponseWriter) (uint16, error)

type SomeStruct struct {
	SomeInt            int16
	SomeUintptr        uintptr
	SliceOfHttpRequest []http.Request
	OneHttpRequest     http.Request
	EmptyString        string
	SmallString        string
	LargeString        string
	DistinctType       Distance
	SomeNilError       error
	SomeConcreteError  error
	LargeSlice         []int
	LargeMap           map[int]string
	UndirectedChan     chan int
	DirectedChan       chan<- bool
	SomeBool           bool

	NonIntializedFuncClosure    func() (io.ReadCloser, error)
	NonIntializedFuncFromStdLib http.HandlerFunc
	NonIntializedFuncWithReturn FuncWithReturn
	IntializedFuncClosure       func() (io.ReadCloser, error)
	IntializedFuncFromStdLib    http.HandlerFunc
	IntializedFuncWithReturn    FuncWithReturn

	ZeroPointerStruct           *url.URL // we won't initiliaze this
	NonZeroPointerStruct        *url.URL
	EvenMoreUrl                 *url.URL
	SliceOfNonZeroPointerStruct []*url.URL

	ComplexxySixFour    complex64
	ComplexyOneTwoEight complex128
	NonStructPointer    *int8
	SomeUnsafety        unsafe.Pointer
}

func TestDir(t *testing.T) {
	t.Parallel()

	t.Run("dump numbers", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)

		vals := map[interface{}]string{
			int(44): `
[
NAME: int
KIND: int
SIGNATURE: [int]
FIELDS: []
METHODS: []
SNIPPET: int(44)
]
`,

			int32(32): `
[
NAME: int32
KIND: int32
SIGNATURE: [int32]
FIELDS: []
METHODS: []
SNIPPET: int32(32)
]
`,

			int64(64): `
[
NAME: int64
KIND: int64
SIGNATURE: [int64]
FIELDS: []
METHODS: []
SNIPPET: int64(64)
]
`,

			float32(32): `
[
NAME: float32
KIND: float32
SIGNATURE: [float32]
FIELDS: []
METHODS: []
SNIPPET: float32(32)
]
`,

			float64(64): `
[
NAME: float64
KIND: float64
SIGNATURE: [float64]
FIELDS: []
METHODS: []
SNIPPET: float64(64)
]
`,

			uintptr(123): `
[
NAME: uintptr
KIND: uintptr
SIGNATURE: [uintptr]
FIELDS: []
METHODS: []
SNIPPET: uintptr(123)
]
`,

			uint64(88): `
[
NAME: uint64
KIND: uint64
SIGNATURE: [uint64]
FIELDS: []
METHODS: []
SNIPPET: uint64(88)
]
`,
		}
		for k, v := range vals {
			res := kama.Dir(k)
			c.Assert(res, qt.Equals, v)
		}
	})

	t.Run("slice on its own is not compacted", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)
		expected := `
[
NAME: []int
KIND: slice
SIGNATURE: [[]int]
FIELDS: []
METHODS: []
SNIPPET: []int{
   int(0),
   int(1),
   int(2),
   int(3),
   int(4),
   int(5),
 ...<9994 more redacted>..}
]
`

		mySlice := bigSlice()

		res := kama.Dir(mySlice)
		c.Assert(res, qt.Equals, expected)
	})

	t.Run("slice in a struct is compacted", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)
		expected := `
[
NAME: github.com/komuw/kama/e2e_test.some
KIND: struct
SIGNATURE: [e2e_test.some *e2e_test.some]
FIELDS: [
	XX []int 
	]
METHODS: []
SNIPPET: some{
  XX: []int{int(0),int(1),int(2),int(3),int(4),int(5), ...<9994 more redacted>..},
}
]
`
		type some struct {
			XX []int
		}
		s := some{XX: bigSlice()}

		res := kama.Dir(s)
		c.Assert(res, qt.Equals, expected)
	})

	t.Run("map on its own is not compacted", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)
		expected := `
[
NAME: map[int]string
KIND: map
SIGNATURE: [map[int]string]
FIELDS: []
METHODS: []
SNIPPET: map[int]string{
   int(0): "0", 
   int(1): "1", 
   int(10): "10", 
   int(100): "100", 
   int(1000): "1000", 
   ...<9997 more redacted>..}
]
`

		myMap := bigMap()

		res := kama.Dir(myMap)
		c.Assert(res, qt.Equals, expected)
	})

	t.Run("map in a struct is compacted", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)
		expected := `
[
NAME: github.com/komuw/kama/e2e_test.some
KIND: struct
SIGNATURE: [e2e_test.some *e2e_test.some]
FIELDS: [
	XX map[int]string 
	]
METHODS: []
SNIPPET: some{
  XX: map[int]string{int(0):"0", int(1):"1", int(10):"10", int(100):"100", int(1000):"1000", ...<9997 more redacted>..},
}
]
`

		type some struct {
			XX map[int]string
		}
		s := some{XX: bigMap()}

		res := kama.Dir(s)
		t.Log(res)
		c.Assert(res, qt.Equals, expected)
	})

	// TODO: fix the fields with `NotImplemented`
	t.Run("struct of varying field types", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)
		expected := `
[
NAME: github.com/komuw/kama/e2e_test.SomeStruct
KIND: struct
SIGNATURE: [e2e_test.SomeStruct *e2e_test.SomeStruct]
FIELDS: [
	SomeInt int16 
	SomeUintptr uintptr 
	SliceOfHttpRequest []http.Request 
	OneHttpRequest http.Request 
	EmptyString string 
	SmallString string 
	LargeString string 
	DistinctType e2e_test.Distance 
	SomeNilError error 
	SomeConcreteError error 
	LargeSlice []int 
	LargeMap map[int]string 
	UndirectedChan chan int 
	DirectedChan chan<- bool 
	SomeBool bool 
	NonIntializedFuncClosure func() (io.ReadCloser, error) 
	NonIntializedFuncFromStdLib http.HandlerFunc 
	NonIntializedFuncWithReturn e2e_test.FuncWithReturn 
	IntializedFuncClosure func() (io.ReadCloser, error) 
	IntializedFuncFromStdLib http.HandlerFunc 
	IntializedFuncWithReturn e2e_test.FuncWithReturn 
	ZeroPointerStruct *url.URL 
	NonZeroPointerStruct *url.URL 
	EvenMoreUrl *url.URL 
	SliceOfNonZeroPointerStruct []*url.URL 
	ComplexxySixFour complex64 
	ComplexyOneTwoEight complex128 
	NonStructPointer *int8 
	SomeUnsafety unsafe.Pointer 
	]
METHODS: []
SNIPPET: SomeStruct{
  SomeInt: int16(13),
  SomeUintptr: uintptr(64902),
  SliceOfHttpRequest: []http.Request{Request{Method: "0",},Request{Method: "1",},Request{Method: "2",},Request{Method: "3",},Request{Method: "4",},Request{Method: "5",}, ...<94 more redacted>..},
  OneHttpRequest: Request{Method: "Hello",},
  EmptyString: "",
  SmallString: "What up?",
  LargeString: "AT last the sleepy atmosphere was stirred—and vig ...<3454 more redacted>..,
  DistinctType: e2e_test.Distance(9131),
  SomeNilError: interface NotImplemented,
  SomeConcreteError: interface NotImplemented,
  LargeSlice: []int{int(0),int(1),int(2),int(3),int(4),int(5), ...<9994 more redacted>..},
  LargeMap: map[int]string{int(0):"0", int(1):"1", int(10):"10", int(100):"100", int(1000):"1000", ...<9997 more redacted>..},
  UndirectedChan: chan int (len=122, cap=10000),
  DirectedChan: chan<- bool (len=1, cap=13),
  SomeBool: true,
  NonIntializedFuncClosure: func() (io.ReadCloser, error),
  NonIntializedFuncFromStdLib: http.HandlerFunc(http.ResponseWriter, *http.Request),
  NonIntializedFuncWithReturn: e2e_test.FuncWithReturn(http.ResponseWriter) (uint16, error),
  IntializedFuncClosure: func() (io.ReadCloser, error),
  IntializedFuncFromStdLib: http.HandlerFunc(http.ResponseWriter, *http.Request),
  IntializedFuncWithReturn: e2e_test.FuncWithReturn(http.ResponseWriter) (uint16, error),
  ZeroPointerStruct: *url.URL(nil),
  NonZeroPointerStruct: &URL{},
  EvenMoreUrl: &URL{Path: "/some/path",},
  SliceOfNonZeroPointerStruct: []*url.URL{&URL{Path: "1",},&URL{Path: "2",},&URL{Path: "3",},&URL{Path: "4",},&URL{Path: "5",},&URL{Path: "6",}, ...<2 more redacted>..},
  ComplexxySixFour: complex64(5+7i),
  ComplexyOneTwoEight: complex128(5+7i),
  NonStructPointer: &int8(14),
  SomeUnsafety: unsafe.Pointer,
}
]
`

		someIntEight := int8(14)
		s := SomeStruct{
			SomeInt:            13,
			SomeUintptr:        uintptr(64_902),
			SliceOfHttpRequest: makeSliceOfHttpRequests(),
			OneHttpRequest:     http.Request{Method: "Hello"},
			EmptyString:        "",
			SmallString:        "What up?",
			LargeString:        longText,
			DistinctType:       Distance(9131),
			SomeNilError:       nil,
			SomeConcreteError:  errors.New("Houston something bad happened"),
			LargeSlice:         bigSlice(),
			LargeMap:           bigMap(),
			UndirectedChan:     bigChan(),
			DirectedChan:       makeDirecteChan(),
			SomeBool:           true,

			IntializedFuncClosure: func() (io.ReadCloser, error) {
				return CustomReadCloser(900), nil
			},
			IntializedFuncFromStdLib: func(rw http.ResponseWriter, r *http.Request) {
				_ = r.Close
				rw.Write([]byte("yo"))
			},
			IntializedFuncWithReturn: func(http.ResponseWriter) (uint16, error) {
				return uint16(1), nil
			},

			NonZeroPointerStruct: &url.URL{},
			EvenMoreUrl:          &url.URL{Path: "/some/path"},
			SliceOfNonZeroPointerStruct: []*url.URL{
				&url.URL{Path: "1"},
				&url.URL{Path: "2"},
				&url.URL{Path: "3"},
				&url.URL{Path: "4"},
				&url.URL{Path: "5"},
				&url.URL{Path: "6"},
				&url.URL{Path: "7"},
				&url.URL{Path: "8"},
			},

			ComplexxySixFour:    complex(float32(5), 7),
			ComplexyOneTwoEight: complex(float64(5), 7),
			NonStructPointer:    &someIntEight,
			SomeUnsafety:        unsafe.Pointer(&someIntEight),
		}

		res := kama.Dir(s)
		c.Assert(res, qt.Equals, expected)
	})

	// TODO: fix the fields with `NotImplemented`
	t.Run("pointer to struct of varying field types", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)
		expected := `
[
NAME: github.com/komuw/kama/e2e_test.SomeStruct
KIND: struct
SIGNATURE: [*e2e_test.SomeStruct e2e_test.SomeStruct]
FIELDS: [
	SomeInt int16 
	SomeUintptr uintptr 
	SliceOfHttpRequest []http.Request 
	OneHttpRequest http.Request 
	EmptyString string 
	SmallString string 
	LargeString string 
	DistinctType e2e_test.Distance 
	SomeNilError error 
	SomeConcreteError error 
	LargeSlice []int 
	LargeMap map[int]string 
	UndirectedChan chan int 
	DirectedChan chan<- bool 
	SomeBool bool 
	NonIntializedFuncClosure func() (io.ReadCloser, error) 
	NonIntializedFuncFromStdLib http.HandlerFunc 
	NonIntializedFuncWithReturn e2e_test.FuncWithReturn 
	IntializedFuncClosure func() (io.ReadCloser, error) 
	IntializedFuncFromStdLib http.HandlerFunc 
	IntializedFuncWithReturn e2e_test.FuncWithReturn 
	ZeroPointerStruct *url.URL 
	NonZeroPointerStruct *url.URL 
	EvenMoreUrl *url.URL 
	SliceOfNonZeroPointerStruct []*url.URL 
	ComplexxySixFour complex64 
	ComplexyOneTwoEight complex128 
	NonStructPointer *int8 
	SomeUnsafety unsafe.Pointer 
	]
METHODS: []
SNIPPET: &SomeStruct{
  SomeInt: int16(13),
  SomeUintptr: uintptr(64902),
  SliceOfHttpRequest: []http.Request{Request{Method: "0",},Request{Method: "1",},Request{Method: "2",},Request{Method: "3",},Request{Method: "4",},Request{Method: "5",}, ...<94 more redacted>..},
  OneHttpRequest: Request{Method: "Hello",},
  EmptyString: "",
  SmallString: "What up?",
  LargeString: "AT last the sleepy atmosphere was stirred—and vig ...<3454 more redacted>..,
  DistinctType: e2e_test.Distance(9131),
  SomeNilError: interface NotImplemented,
  SomeConcreteError: interface NotImplemented,
  LargeSlice: []int{int(0),int(1),int(2),int(3),int(4),int(5), ...<9994 more redacted>..},
  LargeMap: map[int]string{int(0):"0", int(1):"1", int(10):"10", int(100):"100", int(1000):"1000", ...<9997 more redacted>..},
  UndirectedChan: chan int (len=122, cap=10000),
  DirectedChan: chan<- bool (len=1, cap=13),
  SomeBool: true,
  NonIntializedFuncClosure: func() (io.ReadCloser, error),
  NonIntializedFuncFromStdLib: http.HandlerFunc(http.ResponseWriter, *http.Request),
  NonIntializedFuncWithReturn: e2e_test.FuncWithReturn(http.ResponseWriter) (uint16, error),
  IntializedFuncClosure: func() (io.ReadCloser, error),
  IntializedFuncFromStdLib: http.HandlerFunc(http.ResponseWriter, *http.Request),
  IntializedFuncWithReturn: e2e_test.FuncWithReturn(http.ResponseWriter) (uint16, error),
  ZeroPointerStruct: *url.URL(nil),
  NonZeroPointerStruct: &URL{},
  EvenMoreUrl: &URL{Path: "/some/path",},
  SliceOfNonZeroPointerStruct: []*url.URL{&URL{Path: "1",},&URL{Path: "2",},&URL{Path: "3",},&URL{Path: "4",},&URL{Path: "5",},&URL{Path: "6",}, ...<2 more redacted>..},
  ComplexxySixFour: complex64(5+7i),
  ComplexyOneTwoEight: complex128(5+7i),
  NonStructPointer: &int8(14),
  SomeUnsafety: unsafe.Pointer,
}
]
`

		someIntEight := int8(14)
		s := &SomeStruct{
			SomeInt:            13,
			SomeUintptr:        uintptr(64_902),
			SliceOfHttpRequest: makeSliceOfHttpRequests(),
			OneHttpRequest:     http.Request{Method: "Hello"},
			EmptyString:        "",
			SmallString:        "What up?",
			LargeString:        longText,
			DistinctType:       Distance(9131),
			SomeNilError:       nil,
			SomeConcreteError:  errors.New("Houston something bad happened"),
			LargeSlice:         bigSlice(),
			LargeMap:           bigMap(),
			UndirectedChan:     bigChan(),
			DirectedChan:       makeDirecteChan(),
			SomeBool:           true,

			IntializedFuncClosure: func() (io.ReadCloser, error) {
				return CustomReadCloser(900), nil
			},
			IntializedFuncFromStdLib: func(rw http.ResponseWriter, r *http.Request) {
				_ = r.Close
				rw.Write([]byte("yo"))
			},
			IntializedFuncWithReturn: func(http.ResponseWriter) (uint16, error) {
				return uint16(1), nil
			},

			NonZeroPointerStruct: &url.URL{},
			EvenMoreUrl:          &url.URL{Path: "/some/path"},
			SliceOfNonZeroPointerStruct: []*url.URL{
				&url.URL{Path: "1"},
				&url.URL{Path: "2"},
				&url.URL{Path: "3"},
				&url.URL{Path: "4"},
				&url.URL{Path: "5"},
				&url.URL{Path: "6"},
				&url.URL{Path: "7"},
				&url.URL{Path: "8"},
			},

			ComplexxySixFour:    complex(float32(5), 7),
			ComplexyOneTwoEight: complex(float64(5), 7),
			NonStructPointer:    &someIntEight,
			SomeUnsafety:        unsafe.Pointer(&someIntEight),
		}

		res := kama.Dir(s)
		c.Assert(res, qt.Equals, expected)
	})
}
