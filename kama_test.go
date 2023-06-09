package kama

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"go.akshayshah.org/attest"
	"go.uber.org/goleak"

	pkgErrors "github.com/pkg/errors"
)

/*
  go test -timeout 1m -race -cover -v ./...

  This tests are from/inspired by;
  https://github.com/sanity-io/litter/blob/b3546bd0a12c8738436e565b9e016bcd1876403d/dump_test.go

  That repo as of that commit has an MIT license;
  https://github.com/sanity-io/litter/blob/b3546bd0a12c8738436e565b9e016bcd1876403d/LICENSE
*/

const (
	acceptableCodeCoverage = 0.8 // 80%
	kamaWriteDataForTests  = "KAMA_WRITE_DATA_FOR_TESTS"
)

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags

	exitCode := m.Run()
	if exitCode == 0 && testing.CoverMode() != "" {
		coverage := testing.Coverage()
		// note: for some reason the value of `coverage` is always less
		// than the one reported on the terminal by go test
		if coverage < acceptableCodeCoverage {
			fmt.Printf("\n\tThe test code coverage has fallen below the acceptable value of %v. The current value is %v. \n", acceptableCodeCoverage, coverage)
			exitCode = -1
		}
	}

	writeData := os.Getenv(kamaWriteDataForTests) != ""
	if writeData {
		fmt.Printf("\n\t env var %s is set.\n\n", kamaWriteDataForTests)
		os.Exit(77)
	}

	exitCode = leakDetector(exitCode)
	os.Exit(exitCode)
}

// see:
// https://github.com/uber-go/goleak/blob/v1.1.10/testmain.go#L40-L52
func leakDetector(exitCode int) int {
	if exitCode == 0 {
		if err := goleak.Find(); err != nil {
			fmt.Fprintf(os.Stderr, "goleak: Errors on successful test run: %v\n", err)
			exitCode = 1
		}
	}
	return exitCode
}

// dealWithTestData asserts that gotContent is equal to data found at path.
//
// If the environment variable [kamaWriteDataForTests] is set, this func
// will write gotContent to path instead.
func dealWithTestData(t *testing.T, path, gotContent string) {
	t.Helper()

	p, e := filepath.Abs(path)
	attest.Ok(t, e)

	writeData := os.Getenv(kamaWriteDataForTests) != ""
	if writeData {
		attest.Ok(t,
			os.WriteFile(path, []byte(gotContent), 0o644),
		)
		t.Logf("\n\t written testdata to %s\n", path)
		return
	}

	b, e := os.ReadFile(p)
	attest.Ok(t, e)

	expectedContent := string(b)
	attest.Equal(t, gotContent, expectedContent, attest.Sprintf("path: %s", path))
}

func getDataPath(t *testing.T, testPath, testName string) string {
	s := strings.ReplaceAll(testName, " ", "_")
	tName := strings.ReplaceAll(s, "/", "_")

	path := filepath.Join("testdata", testPath, tName) + ".txt"

	return path
}

type (
	BlankStruct struct{}
	BasicStruct struct {
		Public  int
		private int
	}
)
type IntAlias int

func SomeFunction(arg1 string, arg2 int) (string, error) {
	return "", nil
}

func TestPrimitives(t *testing.T) {
	t.Parallel()

	tt := []interface{}{
		false,
		true,
		7,
		12.3535,
		int8(10),
		int16(10),
		int32(10),
		int64(10),
		uint8(10),
		uint16(10),
		uint32(10),
		uint64(10),
		uint(10),
		float32(12.3),
		float64(12.3),
		complex64(12 + 10.5i),
		complex128(-1.2 - 0.1i),
		(func(v int) *int { return &v })(10),
		"string with \"quote\"",
		[]int{1, 2, 3},
		interface{}("hello from interface"),
		BlankStruct{},
		&BlankStruct{},
		BasicStruct{1, 2},
		&BasicStruct{Public: 6_913, private: 90_350},
		IntAlias(10),
		(func(v IntAlias) *IntAlias { return &v })(10),
		SomeFunction,

		func(arg string) (bool, error) { return false, nil },
		nil,
		interface{}(nil),
		make(chan int, 10_000),
		map[int]string{},
		[10_000]int{},
		[]uint16{},
	}

	for _, v := range tt {
		v := v
		Dir(v)
	}
}

type myHandler struct{}

func (h myHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
}

func TestStdlibTypes(t *testing.T) {
	t.Parallel()

	var req http.Request
	var reqPtr *http.Request

	tt := []interface{}{
		errors.New,
		reflect.Value{},
		http.Handle,
		http.HandleFunc,
		http.Handler(myHandler{}),
		req,
		&req,
		reqPtr,
	}

	for _, v := range tt {
		v := v
		Dir(v)
	}
}

func TestThirdPartyTypes(t *testing.T) {
	t.Parallel()

	tt := []interface{}{
		pkgErrors.New,
	}

	for _, v := range tt {
		v := v
		Dir(v)
	}
}
