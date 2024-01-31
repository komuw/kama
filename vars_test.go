package kama

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"sync"
	"testing"
	"time"
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

var thisFunctionVar = ThisFunction //nolint:gochecknoglobals

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

var MyBigSlice = bigSlice() //nolint:gochecknoglobals

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

const BigString = `AT last the sleepy atmosphere was stirred—and vigorously: the murder trial came on in the court. It became the absorbing topic of village talk immediately. Tom could not get away from it. Every reference to the murder sent a shudder to his heart, for his troubled conscience and fears almost persuaded him that these remarks were put forth in his hearing as “feelers”; he did not see how he could be suspected of knowing anything about the murder, but still he could not be comfortable in the midst of this gossip. It kept him in a cold shiver all the time. He took Huck to a lonely place to have a talk with him. It would be some relief to unseal his tongue for a little while; to divide his burden of distress with another sufferer. Moreover, he wanted to assure himself that Huck had remained discreet.
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

var zeroValuePointer *http.Request //nolint:gochecknoglobals

type StructWithTags struct {
	Allowed bool   `json:"enabled"`
	Name    string `json:"their_name"`
}

type Hey struct {
	Another struct {
		Allowed bool   `json:"enabled"`
		Name    string `json:"their_name"`
	}
}

func TestBasicVariables(t *testing.T) {
	t.Parallel()

	tt := []struct {
		tName    string
		variable interface{}
	}{
		{
			tName:    "value struct",
			variable: Person{Name: "John"},
		},
		{
			tName:    "pointer struct",
			variable: &Person{Name: "Jane"},
		},
		{
			tName:    "zero value pointer",
			variable: zeroValuePointer,
		},
		{
			tName:    "function",
			variable: ThisFunction,
		},
		{
			tName:    "function variable",
			variable: thisFunctionVar,
		},
		{
			tName:    "distinct type",
			variable: customerID(9),
		},
		{
			tName:    "big slice",
			variable: MyBigSlice,
		},
		{
			tName:    "big map",
			variable: bigMap(),
		},
		{
			tName:    "big chan",
			variable: bigChan(),
		},
		{
			tName:    "big array",
			variable: bigArray(),
		},
		{
			tName:    "big string",
			variable: BigString,
		},
		{
			tName:    "value struct with slice field",
			variable: SomeStructWIthSlice{Name: "Hello", MyAwesome: bigSlice()},
		},
		{
			tName:    "pointer struct with slice field",
			variable: &SomeStructWIthSlice{Name: "HelloPointery", MyAwesome: bigSlice()},
		},
		{
			tName:    "struct with tags",
			variable: StructWithTags{},
		},
		{
			tName: "embedded struct with tags",
			variable: Hey{Another: struct {
				Allowed bool   `json:"enabled"`
				Name    string `json:"their_name"`
			}{
				Allowed: true,
				Name:    "Jane",
			}},
		},
	}

	for _, v := range tt {
		v := v
		t.Run(v.tName, func(t *testing.T) {
			t.Parallel()

			res := newVari(v.variable)

			path := getDataPath(t, "vars_test.go", v.tName)
			dealWithTestData(t, path, res.String())
		})
	}
}

func TestStdlibVariables(t *testing.T) {
	t.Parallel()

	tt := []struct {
		tName    string
		variable interface{}
	}{
		{
			tName:    "value struct of http.Request",
			variable: http.Request{},
		},

		{
			tName:    "pointer struct of http.Request",
			variable: &http.Request{},
		},
	}

	for _, v := range tt {
		v := v
		t.Run(v.tName, func(t *testing.T) {
			t.Parallel()

			res := newVari(v.variable)

			path := getDataPath(t, "vars_test.go", v.tName)
			dealWithTestData(t, path, res.String())
		})
	}
}

func TestSliceMap(t *testing.T) {
	t.Parallel()

	var nilSlice []string = nil
	var nilMap map[string]int = nil

	tt := []struct {
		tName    string
		variable interface{}
	}{
		{
			tName:    "nil slice",
			variable: nilSlice,
		},
		{
			tName:    "no element slice",
			variable: []string{},
		},
		{
			tName:    "one element slice",
			variable: []string{"hello"},
		},
		{
			tName:    "two element slice",
			variable: []string{"one", "two"},
		},
		{
			tName:    "nil map",
			variable: nilMap,
		},
		{
			tName:    "no element map",
			variable: map[string]int{},
		},
		{
			tName:    "one element map",
			variable: map[string]int{"o": 1},
		},
		{
			tName:    "two element map",
			variable: map[string]int{"o": 1, "two": 2},
		},
	}

	for _, v := range tt {
		v := v

		t.Run(v.tName, func(t *testing.T) {
			t.Parallel()

			res := newVari(v.variable)

			path := getDataPath(t, "vars_test.go", v.tName)
			dealWithTestData(t, path, res.String())
		})
	}
}

type customContext struct{ parent context.Context }

func (d customContext) Deadline() (time.Time, bool)       { return time.Time{}, false }
func (d customContext) Done() <-chan struct{}             { return nil }
func (d customContext) Err() error                        { return nil }
func (d customContext) Value(key interface{}) interface{} { return d.parent.Value(key) }

func TestContexts(t *testing.T) {
	t.Parallel()

	type myContextKeyType string

	{
		// const shortForm = "2006-Jan-02"
		// when, err := time.ParseInLocation(shortForm, "2013-Feb-03", time.UTC)
		// attest.Ok(t, err)

		// // This WithDeadline does not work because the printed value of WithDeadline is not stable
		// // https://github.com/golang/go/blob/39effbc105f5c54117a6011af3c48e3c8f14eca9/src/context/context.go#L654-L657
		// ctxWithDeadline, cancel := context.WithDeadline(context.Background(), when)
		// t.Cleanup(func() {
		// 	cancel()
		// })
	}

	type StructWithContext struct {
		Name   string
		Age    int64
		OurCtx context.Context
	}
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
	})

	ctxWithValue := context.WithValue(context.TODO(), myContextKeyType("ctxWithValueType"), "OKAYY") //nolint:gocritic

	encapsulatedStdlibCtx := context.WithValue(ctxWithCancel, myContextKeyType("myContextKeyType"), "ThisIsSomeContextValue")

	tt := []struct {
		tName    string
		variable interface{}
	}{
		{
			tName:    "stdlib context TODO",
			variable: context.TODO(), //nolint:gocritic
		},
		{
			tName:    "stdlib context Background",
			variable: context.Background(),
		},
		{
			tName:    "stdlib context WithCancel",
			variable: ctxWithCancel,
		},
		// This WithDeadline does not work because the printed value of WithDeadline is not stable
		// https://github.com/golang/go/blob/39effbc105f5c54117a6011af3c48e3c8f14eca9/src/context/context.go#L654-L657
		// {
		// 	tName:    "stdlib context WithDeadline",
		// 	variable: ctxWithDeadline,
		// },
		{
			tName:    "stdlib context WithValue",
			variable: ctxWithValue,
		},
		{
			tName:    "stdlib context encapsulated",
			variable: encapsulatedStdlibCtx,
		},
		{
			tName:    "custom context",
			variable: customContext{context.Background()},
		},
		{
			tName:    "context inside struct",
			variable: StructWithContext{Name: "John", Age: 763, OurCtx: encapsulatedStdlibCtx},
		},
	}

	for _, v := range tt {
		v := v
		t.Run(v.tName, func(t *testing.T) {
			t.Parallel()

			res := newVari(v.variable)

			path := getDataPath(t, "vars_test.go", v.tName)
			dealWithTestData(t, path, res.String())
		})
	}
}

func TestLong(t *testing.T) {
	// t.Parallel() // This cannot be ran in Parallel since it mutates a global var.

	oldCfg := cfg

	type Hey struct {
		BigSlice  []int
		BigArray  [10_000]int
		BigMap    map[int]string
		BigString string
	}
	h := Hey{
		BigSlice:  bigSlice(),
		BigArray:  bigArray(),
		BigMap:    bigMap(),
		BigString: BigString,
	}

	tt := []struct {
		tName    string
		variable interface{}
		c        Config
	}{
		{
			tName:    "no_config",
			variable: h,
		},
		{
			tName:    "default_config",
			variable: h,
			c:        oldCfg,
		},
		{
			tName:    "maxLength_config",
			variable: h,
			c:        Config{MaxLength: math.MaxInt},
		},
		{
			tName:    "maxLength_big-string_config",
			variable: BigString,
			c:        Config{MaxLength: math.MaxInt},
		},
		{
			tName:    "maxLength_empty-string_config",
			variable: "",
			c:        Config{MaxLength: math.MaxInt},
		},
	}

	for _, v := range tt {
		v := v

		t.Run(v.tName, func(t *testing.T) {
			// t.Parallel() // This cannot be ran in Parallel since it mutates a global var.

			{ // Set the new config and schedule to return old config.
				onceCfg = &sync.Once{}
				_ = Dir("", v.c)
				t.Cleanup(func() {
					cfg = oldCfg
				})
			}

			res := newVari(v.variable)

			path := getDataPath(t, "vars_test.go", v.tName)
			dealWithTestData(t, path, res.String())
		})
	}
}
