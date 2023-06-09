package kama

import (
	"testing"

	"go.akshayshah.org/attest"
)

func TestStdlibPackages(t *testing.T) {
	t.Parallel()

	tt := []struct {
		importPath string
		expected   pak
	}{
		{
			"errors", pak{
				Name: "errors",

				Constants: []string{},
				Variables: []string{},
				Functions: []string{
					"As(err error, target any) bool",
					"Is(err error, target error) bool",
					"Join(errs ...error) error",
					"New(text string) error",
					"Unwrap(err error) error",
				},
				Types: map[string][]string{},
			},
		},

		{
			"archive/tar", pak{
				Name:      "archive/tar",
				Constants: []string{"FormatGNU Format", "FormatPAX Format", "FormatUSTAR Format", "FormatUnknown Format", "TypeBlock untyped rune", "TypeChar untyped rune", "TypeCont untyped rune", "TypeDir untyped rune", "TypeFifo untyped rune", "TypeGNULongLink untyped rune", "TypeGNULongName untyped rune", "TypeGNUSparse untyped rune", "TypeLink untyped rune", "TypeReg untyped rune", "TypeRegA untyped rune", "TypeSymlink untyped rune", "TypeXGlobalHeader untyped rune", "TypeXHeader untyped rune"},
				Variables: []string{"ErrFieldTooLong error", "ErrHeader error", "ErrInsecurePath error", "ErrWriteAfterClose error", "ErrWriteTooLong error"},
				Functions: []string{"FileInfoHeader(fi io/fs.FileInfo, link string) (*Header, error)", "NewReader(r io.Reader) *Reader", "NewWriter(w io.Writer) *Writer"},
				Types: map[string][]string{
					"Format int": {
						"(Format) String() string",
					},
					"Header struct": {
						"(*Header) FileInfo() io/fs.FileInfo",
					},
					"Reader struct": {
						"(*Reader) Next() (*Header, error)",
						"(*Reader) Read(b []byte) (int, error)",
					},
					"Writer struct": {
						"(*Writer) Close() error",
						"(*Writer) Flush() error",
						"(*Writer) Write(b []byte) (int, error)",
						"(*Writer) WriteHeader(hdr *Header) error",
					},
				},
			},
		},
	}

	for _, v := range tt {
		v := v

		p, err := newPak(v.importPath)
		if err != nil {
			t.Errorf("\ngot \n\t%#+v \nwanted \n\t%#+v", err, v.expected)
		}

		attest.Equal(t, p, v.expected)
	}
}

func TestThirdPartyPackages(t *testing.T) {
	t.Parallel()

	tt := []struct {
		importPath string
		expected   pak
	}{
		{
			"github.com/pkg/errors", pak{
				Name:      "github.com/pkg/errors",
				Constants: []string{},
				Variables: []string{},
				Functions: []string{
					"As(err error, target interface{}) bool",
					"Cause(err error) error",
					"Errorf(format string, args ...interface{}) error",
					"Is(err error, target error) bool",
					"New(message string) error",
					"Unwrap(err error) error",
					"WithMessage(err error, message string) error",
					"WithMessagef(err error, format string, args ...interface{}) error",
					"WithStack(err error) error",
					"Wrap(err error, message string) error",
					"Wrapf(err error, format string, args ...interface{}) error",
				},
				Types: map[string][]string{
					"Frame uintptr": {
						"(Frame) Format(s fmt.State, verb rune)",
						"(Frame) MarshalText() ([]byte, error)",
					},
					"StackTrace []Frame": {
						"(StackTrace) Format(s fmt.State, verb rune)",
					},
				},
			},
		},
	}

	for _, v := range tt {
		v := v

		p, err := newPak(v.importPath)
		if err != nil {
			t.Errorf("\ngot \n\t%#+v \nwanted \n\t%#+v", err, v.expected)
		}

		attest.Equal(t, p, v.expected)
	}
}

func TestError(t *testing.T) {
	t.Parallel()

	_, err := newPak("github.com/pkg/NoSuchModule")
	if err == nil {
		t.Errorf("got no error, yet expected an error")
	}
	attest.Subsequence(t, err.Error(), "no required module provides package")
}
