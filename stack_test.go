package kama

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func Test_getStackTrace(t *testing.T) {
	t.Run("test-get-stacktraces", func(t *testing.T) {
		c := qt.New(t)

		got := a()

		c.Assert(got[1], qt.Contains, "func a() []string {")
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
	d()
	return getStackTrace()
}

func d() {
	stackp()
}

// func Test_stackp(t *testing.T) {
// 	t.Run("test-stackp", func(t *testing.T) {
// 		c := qt.New(t)

// 		got := d()

// 		c.Assert(got[1], qt.Contains, "func a() []string {")
// 	})
// }
