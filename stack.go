package kama

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// TODO: should be iota
var colors = map[string]int{
	// Go compiler == cyan
	// Third party == magenta
	// Your code == green

	"black":   30,
	"red":     31,
	"green":   32,
	"yellow":  33,
	"blue":    34,
	"magenta": 35,
	"cyan":    36,
	"white":   37,

	"DEFAULT": 39,
	"RESET":   0,
}

func reset() {
	const escape = "\x1b"
	const r = 0
	fmt.Printf("%s[%dm", escape, r)
}

func setColor(code int, bold bool) {
	const escape = "\x1b"
	if bold {
		fmt.Printf("%s[1%dm", escape, code)
	} else {
		fmt.Printf("%s[%dm", escape, code)
	}
}

func printWithColor(s string, color string, bold bool) {
	defer reset()
	color = strings.ToLower(color)
	setColor(colors[color], bold)
	fmt.Println(s)
}

// TODO: move the exported funcs into kama.go
func Stack() string {
	x := debug.Stack()

	lines := [][]byte{}
	curLine := []byte{}
	for _, v := range x {
		if v == 10 {
			lines = append(lines, curLine)
			curLine = []byte{}
		} else {
			curLine = append(curLine, v)
		}
	}

	for _, v := range lines {
		it := string(v)
		if strings.Contains(it, "runtime/") {
			// compiler
			printWithColor(it, "cyan", false)
		} else if strings.Contains(it, "[running]") {
			// treat it like compiler
			printWithColor(it, "cyan", false)
		} else if strings.Contains(it, "github.com/komuw/") {
			// third party
			printWithColor(it, "magenta", false)
		} else if strings.Contains(it, "/home/komuw") {
			// your code
			printWithColor(it, "green", true)
		} else {
			printWithColor(it, "green", false)
		}
	}

	curLine = nil
	lines = nil
	return "" //string(x)
}

func Stackp() {
	fmt.Println(Stack())
}
