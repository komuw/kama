package kama

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const (
	compilerColor   = "blue"
	thirdPartyColor = "yellow"
	yourColor       = "red"
)

func stackp() {
	goModCache := os.Getenv("GOMODCACHE")
	re := regexp.MustCompile(`\d:`) // this pattern is the one created in `readLastLine()`

	traces := getStackTrace()
	if len(traces) > 0 {
		printWithColor(
			fmt.Sprintf("LEGEND:\n compiler: %s\n thirdParty: %s\n yours: %s\n", compilerColor, thirdPartyColor, yourColor),
			"DEFAULT",
			true,
		)
	}

	for _, v := range traces {
		if strings.Contains(v, "go/src/") {
			// compiler
			printWithColor(v, compilerColor, false)
		} else if goModCache != "" && strings.Contains(v, goModCache) {
			// third party
			printWithColor(v, thirdPartyColor, false)
		} else if re.MatchString(v) {
			// this is code snippets
			printWithColor(v, yourColor, false)
		} else {
			// your code
			printWithColor(v, yourColor, true)
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

	txtLast := ""
	txtLastButOne := ""
	if len(frms) > 0 {
		txtLast = readLastLine(frms[0].file, int64(frms[0].line))
	}
	if len(frms) > 1 {
		txtLastButOne = readLastLine(frms[1].file, int64(frms[1].line))
	}

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

func printWithColor(s, color string, bold bool) {
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

	if noColor() {
		fmt.Fprintln(os.Stderr, s)
	} else {
		defer reset()
		color = strings.ToLower(color)
		setColor(colors[color], bold)
		fmt.Fprintln(os.Stderr, s)
	}
}

// noColor indicates whether the terminal in question supports color.
// see: https://github.com/fatih/color/blob/v1.13.0/color.go#L22-L23
func noColor() bool {
	_, exists := os.LookupEnv("NO_COLOR")
	if exists {
		return true
	}

	if strings.ToLower(os.Getenv("TERM")) == "dumb" {
		return true
	}

	return false
}
