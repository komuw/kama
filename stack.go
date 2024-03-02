package kama

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const (
	runtimeColor    = "blue"
	thirdPartyColor = "yellow"
	yourColor       = "red"
)

func stackp(w io.Writer, noColor bool) {
	goModCache := os.Getenv("GOMODCACHE")
	re := regexp.MustCompile(`\d:`) // this pattern is the one created in `readLastLine()`

	traces := getStackTrace()
	if len(traces) > 0 && (!noColor) {
		printWithColor(
			w,
			noColor,
			fmt.Sprintf("LEGEND:\n compiler: %s\n thirdParty: %s\n yours: %s\n", runtimeColor, thirdPartyColor, yourColor),
			"DEFAULT",
			true,
		)
	}

	for _, v := range traces {
		if strings.Contains(v, "go/src/") {
			// compiler
			printWithColor(w, noColor, v, runtimeColor, false)
		} else if goModCache != "" && strings.Contains(v, goModCache) {
			// third party
			printWithColor(w, noColor, v, thirdPartyColor, false)
		} else if re.MatchString(v) {
			// this is code snippets
			printWithColor(w, noColor, v, yourColor, false)
		} else {
			// your code
			printWithColor(w, noColor, v, yourColor, true)
		}
	}
}

// frm is like a runtime.Frame
type frm struct {
	file     string
	line     int
	function string
}

func getStackTrace() []string {
	const maxStackLength = 50

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

	n := 0
	txtLast := ""
	txtLastButOne := ""
	for k := range frms {
		if strings.Contains(frms[k].file, "komuw/kama") || strings.Contains(frms[k].file, "/kama/") || strings.Contains(frms[k].file, "go/src/") {
			// Do not display expanded source code for this library or Go runtime.
			n = n + 1
		} else {
			break
		}
	}

	if len(frms) > n {
		txtLast = readLastLine(frms[n].file, int64(frms[n].line))
	}
	next := n + 1
	if len(frms) > next {
		txtLastButOne = readLastLine(frms[next].file, int64(frms[next].line))
	}

	traces := []string{}
	for k, v := range frms {
		traces = append(traces, fmt.Sprintf("\t%s:%d %s", v.file, v.line, v.function))

		if k == n && txtLast != "" {
			traces = append(traces, txtLast)
		}

		if k == next && txtLastButOne != "" {
			traces = append(traces, txtLastButOne)
		}
	}

	return traces
}

func readLastLine(file string, line int64) string {
	txt := ""

	f, err := os.Open(filepath.Clean(file))
	if err != nil {
		return txt
	}
	defer func() {
		_ = f.Close()
	}()

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

func reset(w io.Writer) {
	const escape = "\x1b"
	const r = 0
	_, _ = fmt.Fprintf(w, "%s[%dm", escape, r)
}

func setColor(w io.Writer, code int, bold bool) {
	const escape = "\x1b"
	if bold {
		_, _ = fmt.Fprintf(w, "%s[1%dm", escape, code)
	} else {
		_, _ = fmt.Fprintf(w, "%s[%dm", escape, code)
	}
}

func printWithColor(w io.Writer, noColor bool, s, color string, bold bool) {
	// TODO: should be iota
	colors := map[string]int{
		// Go compiler == compilerColor
		// Third party == thirdPartyColor
		// Your code == yourColor

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

	if hasNoColor() || noColor {
		_, _ = fmt.Fprintln(w, s)
	} else {
		defer reset(w)
		color = strings.ToLower(color)
		setColor(w, colors[color], bold)
		_, _ = fmt.Fprintln(w, s)
	}
}

// hasNoColor indicates whether the terminal in question supports color.
// see: https://github.com/fatih/color/blob/v1.13.0/color.go#L22-L23
func hasNoColor() bool {
	_, exists := os.LookupEnv("NO_COLOR")
	if exists {
		return true
	}

	if strings.ToLower(os.Getenv("TERM")) == "dumb" {
		return true
	}

	return false
}
