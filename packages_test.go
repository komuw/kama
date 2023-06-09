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
	}{
		{
			tName:      "third party github.com/pkg/errors",
			importPath: "github.com/pkg/errors",
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

func TestError(t *testing.T) {
	t.Parallel()

	_, err := newPak("github.com/pkg/NoSuchModule")
	if err == nil {
		t.Errorf("got no error, yet expected an error")
	}
	attest.Subsequence(t, err.Error(), "no required module provides package")
}
