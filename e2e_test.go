package kama_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unsafe"

	"github.com/komuw/kama"
	"go.akshayshah.org/attest"
)

const kamaWriteDataForTests = "KAMA_WRITE_DATA_FOR_TESTS"

const longText = `AT last the sleepy atmosphere was stirred—and vigorously: the murder trial came on in the court. It became the absorbing topic of village talk immediately. Tom could not get away from it. Every reference to the murder sent a shudder to his heart, for his troubled conscience and fears almost persuaded him that these remarks were put forth in his hearing as “feelers”; he did not see how he could be suspected of knowing anything about the murder, but still he could not be comfortable in the midst of this gossip. It kept him in a cold shiver all the time. He took Huck to a lonely place to have a talk with him. It would be some relief to unseal his tongue for a little while; to divide his burden of distress with another sufferer. Moreover, he wanted to assure himself that Huck had remained discreet.
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

type secretKey string

func (s secretKey) String() string {
	if len(s) <= 0 {
		return "secretKey(<EMPTY>)"
	}
	return fmt.Sprintf("secretKey(%s<REDACTED>)", string(s[0]))
}

func (s secretKey) GoString() string {
	return s.String()
}

// dealWithTestData asserts that gotContent is equal to data found at path.
//
// If the environment variable [kamaWriteDataForTests] is set, this func
// will write gotContent to path instead.
func dealWithTestData(t *testing.T, path, gotContent string) {
	t.Helper()

	path = strings.ReplaceAll(path, ".go", "")

	p, e := filepath.Abs(path)
	attest.Ok(t, e)

	writeData := os.Getenv(kamaWriteDataForTests) != ""
	if writeData {
		e := os.MkdirAll(filepath.Dir(p), 0o755)
		attest.Ok(t, e)

		err := os.WriteFile(path, []byte(gotContent), 0o644)
		attest.Ok(t, err)
		t.Logf("\n\t written testdata to %s\n", path)
		return
	}

	b, e := os.ReadFile(p)
	attest.Ok(t, e)

	expectedContent := string(b)
	attest.Equal(t, gotContent, expectedContent, attest.Sprintf("path: %s", path))
}

func getDataPath(t *testing.T, testPath, testName string) string { //nolint:unparam
	t.Helper()

	s := strings.ReplaceAll(testName, " ", "_")
	tName := strings.ReplaceAll(s, "/", "_")

	path := filepath.Join("testdata", testPath, tName) + ".txt"

	return path
}

func TestDir(t *testing.T) {
	t.Parallel()

	type someStructWithSlice struct {
		XX []int
	}
	structWithSlice := someStructWithSlice{XX: bigSlice()}

	type someStructWithMap struct {
		XX map[int]string
	}
	structWithMap := someStructWithMap{XX: bigMap()}

	normalString := "heyWorldOkay"

	tt := []struct {
		tName string
		item  interface{}
	}{
		{
			tName: "number int",
			item:  int(44),
		},
		{
			tName: "number int32",
			item:  int32(32),
		},
		{
			tName: "number int64",
			item:  int64(64),
		},
		{
			tName: "number float32",
			item:  float32(32),
		},
		{
			tName: "number float64",
			item:  float64(64),
		},
		{
			tName: "number uintptr",
			item:  uintptr(123),
		},
		{
			tName: "number uint64",
			item:  uint64(88),
		},
		{
			tName: "slice on its own is not compacted",
			item:  bigSlice(),
		},
		{
			tName: "slice in a struct is not compacted",
			item:  structWithSlice,
		},
		{
			tName: "map on its own is not compacted",
			item:  bigMap(),
		},
		{
			tName: "map in a struct is not compacted",
			item:  structWithMap,
		},
		{
			tName: "normal string",
			item:  normalString,
		},
		{
			tName: "custom string type",
			item:  secretKey("ThisIsMyPrivateApiKey"),
		},
	}

	for _, v := range tt {
		v := v

		t.Run(v.tName, func(t *testing.T) {
			t.Parallel()

			res := kama.Dir(v.item)
			path := getDataPath(t, "e2e_test.go", v.tName)
			dealWithTestData(t, path, res)
		})
	}

	t.Run("struct of varying field types", func(t *testing.T) {
		t.Parallel()

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
				_, _ = io.WriteString(rw, "yo")
			},
			IntializedFuncWithReturn: func(http.ResponseWriter) (uint16, error) {
				return uint16(1), nil
			},

			NonZeroPointerStruct: &url.URL{},
			EvenMoreUrl:          &url.URL{Path: "/some/path"},
			SliceOfNonZeroPointerStruct: []*url.URL{
				{Path: "1"},
				{Path: "2"},
				{Path: "3"},
				{Path: "4"},
				{Path: "5"},
				{Path: "6"},
				{Path: "7"},
				{Path: "8"},
			},

			ComplexxySixFour:    complex(float32(5), 7),
			ComplexyOneTwoEight: complex(float64(5), 7),
			NonStructPointer:    &someIntEight,
			SomeUnsafety:        unsafe.Pointer(&someIntEight),
		}

		res := kama.Dir(s)
		path := getDataPath(t, "e2e_test.go", "struct_of_varying_field_types.txt")
		dealWithTestData(t, path, res)
	})

	t.Run("pointer to struct of varying field types", func(t *testing.T) {
		t.Parallel()

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
				_, _ = io.WriteString(rw, "yo")
			},
			IntializedFuncWithReturn: func(http.ResponseWriter) (uint16, error) {
				return uint16(1), nil
			},

			NonZeroPointerStruct: &url.URL{},
			EvenMoreUrl:          &url.URL{Path: "/some/path"},
			SliceOfNonZeroPointerStruct: []*url.URL{
				{Path: "1"},
				{Path: "2"},
				{Path: "3"},
				{Path: "4"},
				{Path: "5"},
				{Path: "6"},
				{Path: "7"},
				{Path: "8"},
			},

			ComplexxySixFour:    complex(float32(5), 7),
			ComplexyOneTwoEight: complex(float64(5), 7),
			NonStructPointer:    &someIntEight,
			SomeUnsafety:        unsafe.Pointer(&someIntEight),
		}

		res := kama.Dir(s)
		path := getDataPath(t, "e2e_test.go", "pointer_to_struct_of_varying_field_types.txt")
		dealWithTestData(t, path, res)
	})

	t.Run("slice of http.Request value structs", func(t *testing.T) {
		t.Parallel()

		sliceOfStruct := func() []http.Request {
			xx := []http.Request{}
			for i := 0; i < 10_000; i++ {
				xx = append(xx, http.Request{Method: fmt.Sprintf("%d", i)})
			}
			return xx
		}

		s := sliceOfStruct()
		res := kama.Dir(s)
		path := getDataPath(t, "e2e_test.go", "slice_of_http_Request_value_structs.txt")
		dealWithTestData(t, path, res)
	})
}

func TestAllAboutInterfaces(t *testing.T) {
	t.Parallel()

	t.Run("interfaces should be well represented", func(t *testing.T) {
		t.Parallel()

		var SomeNilError error = nil
		var SomeConcreteError error = errors.New("unable to read from ftp file")
		var SomeReader io.Reader = strings.NewReader("hello my reader")
		var NilEmptyInterface interface{} = nil
		var NonNilEmptyInterface interface{} = 9

		type SomeStructWithInterfaces struct {
			AAA               io.Reader
			SomeNilError      error
			SomeConcreteError error
		}
		someOne := SomeStructWithInterfaces{
			AAA:               strings.NewReader("hello from strings NewReader"),
			SomeConcreteError: errors.New("houston something bad happened"),
		}
		someTwo := &someOne

		tt := []struct {
			tName string
			item  interface{}
		}{
			{
				tName: "SomeNilError",
				item:  SomeNilError,
			},
			{
				tName: "SomeConcreteError",
				item:  SomeConcreteError,
			},
			{
				tName: "&SomeConcreteError",
				item:  &SomeConcreteError,
			},
			{
				tName: "SomeReader",
				item:  SomeReader,
			},
			{
				tName: "NilEmptyInterface",
				item:  NilEmptyInterface,
			},
			{
				tName: "NonNilEmptyInterface",
				item:  NonNilEmptyInterface,
			},
			{
				tName: "someOne",
				item:  someOne,
			},
			{
				tName: "someTwo",
				item:  someTwo,
			},
		}

		for _, v := range tt {
			v := v

			t.Run(v.tName, func(t *testing.T) {
				t.Parallel()

				res := kama.Dir(v.item)

				path := getDataPath(t, "e2e_test.go", v.tName)
				dealWithTestData(t, path, res)
			})
		}
	})
}
