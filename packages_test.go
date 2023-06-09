package kama

import (
	"testing"

	"go.akshayshah.org/attest"
)

func TestStdlibPackages(t *testing.T) {
	t.Parallel()

	tt := []struct {
		tName      string
		importPath string
	}{
		{
			tName:      "stdlib errors pkg",
			importPath: "errors",
		},

		{
			tName:      "stdlib archive/tar pkg",
			importPath: "archive/tar",
		},
	}

	for _, v := range tt {
		v := v
		t.Run(v.tName, func(t *testing.T) {
			t.Parallel()

			p, err := newPak(v.importPath)
			attest.Ok(t, err)

			path := getDataPath(t, "packages_test.go", v.tName)
			dealWithTestData(t, path, p.String())
		})
	}
}

func TestThirdPartyPackages(t *testing.T) {
	t.Parallel()

	tt := []struct {
		tName      string
		importPath string
		expected   pak
	}{
		{
			"third party github.com/pkg/errors",
			"github.com/pkg/errors",
			pak{
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
		t.Run(v.tName, func(t *testing.T) {
			t.Parallel()

			p, err := newPak(v.importPath)
			if err != nil {
				t.Errorf("\ngot \n\t%#+v \nwanted \n\t%#+v", err, v.expected)
			}

			attest.Equal(t, p, v.expected)
		})
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
