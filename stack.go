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
	traces := getStackTrace()
	for _, v := range traces {
		if strings.Contains(v, "go/src/") {
			// compiler
			printWithColor(v, "cyan", false)
		} else if strings.Contains(v, "github.com/komuw/") {
			// third party
			printWithColor(v, "magenta", false)
		} else if strings.Contains(v, "/home/komuw") {
			// your code
			printWithColor(v, "green", true)
		} else {
			printWithColor(v, "green", false)
		}
	}
}

const maxStackLength = 50

// frm is like a runtime.Frame
type frm struct {
	file     string
	line     int
	function string
}

func getStackTrace() []string {
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
	txtLastButOne := readLastLine(frms[1].file, int64(frms[1].line))

	traces := []string{}
	for k, v := range frms {
		traces = append(traces, fmt.Sprintf("\t%s:%d %s", v.file, v.line, v.function))

		if k == 0 && txtLast != "" {
			traces = append(traces, txtLast)
		}

		if k == 1 && txtLastButOne != "" {
			traces = append(traces, txtLastButOne)
		}

	}

	return traces
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
