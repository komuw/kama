package kama

import (
	"bytes"
	"io"
	"testing"

	"go.akshayshah.org/attest"
)

func Test_getStackTrace(t *testing.T) {
	t.Parallel()

	t.Run("test-get-stacktraces", func(t *testing.T) {
		t.Parallel()

		got := a()

		attest.Subsequence(t, got[0], "github.com/komuw/kama.a")
	})
}

func a() []string {
	return b()
}

func b() []string {
	return c()
}

// c calls social networks to get network graphs
func c() []string {
	return getStackTrace()
}

func Test_stackp(t *testing.T) {
	t.Parallel()

	t.Run("test-stackp", func(t *testing.T) {
		t.Parallel()

		w := &bytes.Buffer{}
		d(w)

		attest.Subsequence(t, w.String(), "kama.Test_stackp")
	})

	t.Run("test-stack", func(t *testing.T) {
		t.Parallel()

		res := Stack()

		attest.Subsequence(t, res, "github.com/komuw/kama.Test_stackp")
	})
}

func d(w io.Writer) {
	stackp(w, false)
}
