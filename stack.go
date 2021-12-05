package kama

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
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
	fmt.Fprintf(os.Stderr, "%s[%dm", escape, r)
}

func setColor(code int, bold bool) {
	const escape = "\x1b"
	if bold {
		fmt.Fprintf(os.Stderr, "%s[1%dm", escape, code)
	} else {
		fmt.Fprintf(os.Stderr, "%s[%dm", escape, code)
	}
}

func printWithColor(s string, color string, bold bool) {
	defer reset()
	color = strings.ToLower(color)
	setColor(colors[color], bold)
	fmt.Fprintln(os.Stderr, s)
}

func stackp() {
	// x := debug.Stack()

	// lines := [][]byte{}
	// curLine := []byte{}
	// for _, v := range x {
	// 	if v == 10 {
	// 		lines = append(lines, curLine)
	// 		curLine = []byte{}
	// 	} else {
	// 		curLine = append(curLine, v)
	// 	}
	// }

	// tab := "\t"
	// for _, v := range lines {
	// 	it := string(v)
	// 	if strings.Contains(it, "runtime/") {
	// 		// compiler
	// 		printWithColor(it, "cyan", false)
	// 	} else if strings.Contains(it, "[running]") {
	// 		// treat it like compiler
	// 		printWithColor(it, "cyan", false)
	// 	} else if strings.Contains(it, "github.com/komuw/") {
	// 		// third party
	// 		printWithColor(it, "magenta", false)
	// 	} else if strings.Contains(it, "/home/komuw") {
	// 		// your code
	// 		printWithColor(tab+it, "green", true)
	// 	} else {
	// 		printWithColor(tab+it, "green", false)
	// 	}
	// }

	// curLine = nil
	// lines = nil

	fmt.Println(getStackTrace())
}

const maxStackLength = 50

// frm is like a runtime.Frame
type frm struct {
	file     string
	line     int
	function string
}

func getStackTrace() string {
	stackBuf := make([]uintptr, maxStackLength)
	length := runtime.Callers(4, stackBuf[:])
	stack := stackBuf[:length]

	frames := runtime.CallersFrames(stack)

	var frms []frm
	for {
		frame, more := frames.Next()
		frms = append(frms, frm{file: frame.File, line: frame.Line, function: frame.Function})
		if !more {
			break
		}
	}

	txtLast := readLastLine(frms[0].file, int64(frms[0].line))

	trace := ""
	for k, v := range frms {
		trace = trace + fmt.Sprintf("\n\t%s:%d %s", v.file, v.line, v.function)
		if k == 0 && txtLast != "" {
			trace = trace + "\n" + txtLast
		}
	}

	return trace
}

func readLastLine(file string, line int64) string {
	txt := ""

	f, err := os.Open(file)
	if err != nil {
		return txt
	}
	defer f.Close()

	// NB: scanner will error if file is larger than scanner.MaxScanTokenSize
	// which is about 65_000 lines. We should do something about that in the future.
	scanner := bufio.NewScanner(f)

	curLine := int64(0)
	for scanner.Scan() {
		curLine = curLine + 1

		// show the code at that line, including the previous two lines and the next one line.
		if curLine == line-2 || curLine == line-1 || curLine == line+1 {
			txt = txt + fmt.Sprintf("%d: %s", curLine, "\t\t") + scanner.Text() + "\n"
		} else if curLine == line {
			txt = txt + fmt.Sprintf("%d: %s", curLine, "\t\t") + "--> " + scanner.Text() + "\n"
		}
	}

	err = scanner.Err()
	if err != nil {
		return txt
	}

	return strings.TrimSuffix(txt, "\n")
}
