package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/komuw/kama"
)

var sawyer = `AT last the sleepy atmosphere was stirred—and vigorously: the murder trial came on in the court. It became the absorbing topic of village talk immediately. Tom could not get away from it. Every reference to the murder sent a shudder to his heart, for his troubled conscience and fears almost persuaded him that these remarks were put forth in his hearing as “feelers”; he did not see how he could be suspected of knowing anything about the murder, but still he could not be comfortable in the midst of this gossip. It kept him in a cold shiver all the time. He took Huck to a lonely place to have a talk with him. It would be some relief to unseal his tongue for a little while; to divide his burden of distress with another sufferer. Moreover, he wanted to assure himself that Huck had remained discreet.
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

type myHandler struct{ Logger *log.Logger }

func (h myHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
}

type House struct {
	Name   string
	Age    int16
	Alas   uintptr
	Chairs []int
	HTTP   []http.Request
	Hello  http.Request
	One    string
	Two    string
}

type Hello struct{ Age uint64 }

func main() {
	x := []int{}
	for i := 0; i < 10_000; i++ {
		x = append(x, i)
	}

	h := []http.Request{}
	for i := 0; i < 10; i++ {
		h = append(h, http.Request{Method: fmt.Sprint(i)})
	}

	house := House{
		Name:   sawyer,
		Age:    64,
		Alas:   uintptr(90),
		Chairs: x,
		HTTP:   h,
		Hello:  http.Request{Method: "HEllo"},
		One:    "",
		Two:    "Twooot",
	}
	kama.Dirp(house)

	// kama.Dirp(&Hello{78})

	// kama.Dirp(&house)
}

// TODO: clean up

// TODO: fuzz test

// TODO: add documentation for `kama`

// TODO: add a command line api.
//   eg; `kama http.Request` or `kama http`
// have a look at `golang.org/x/tools/cmd/godex`
