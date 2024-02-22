// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kama

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestSmallDiff(t *testing.T) {
	t.Parallel()

	tests := []struct {
		text1 string
		text2 string
		diff  string
	}{
		{"a b c", "a b d e f", "a b -c +d +e +f"},
		{"", "a b c", "+a +b +c"},
		{"a b c", "", "-a -b -c"},
		{"a b c", "d e f", "-a -b -c +d +e +f"},
		{"a b c d e f", "a b d e f", "a b -c d e f"},
		{"a b c e f", "a b c d e f", "a b c +d e f"},
	}

	for _, tt := range tests {
		// Turn spaces into \n.
		text1 := strings.ReplaceAll(tt.text1, " ", "\n")
		if text1 != "" {
			text1 += "\n"
		}
		text2 := strings.ReplaceAll(tt.text2, " ", "\n")
		if text2 != "" {
			text2 += "\n"
		}
		out := diff(text1, text2)
		// Cut final \n, cut spaces, turn remaining \n into spaces.
		out = strings.ReplaceAll(strings.ReplaceAll(strings.TrimSuffix(out, "\n"), " ", ""), "\n", " ")
		if out != tt.diff {
			t.Errorf("diff(%q, %q) = %q, want %q", text1, text2, out, tt.diff)
		}
	}
}

func TestOya(t *testing.T) {
	a := Dir(http.Request{Method: "GET"})
	b := Dir(http.Request{Method: "POST"})

	x := diff(a, b)
	fmt.Println("x: ", x)
}
