package main

import "testing"

var rgflags, words []string

func assertLengths(t *testing.T, nrgflags, nwords int) {
	if len(rgflags) != nrgflags || len(words) != nwords {
		t.Errorf("len error: %d (%d) %d (%d)", len(rgflags), nrgflags, len(words), nwords)
	}
}

func assertContent(t *testing.T, exp, actual []string) {
	error := false
	m := len(exp)
	if len(exp) != len(actual) {
		t.Errorf("len error: %d (%d)", len(actual), len(exp))
		error = true
		if len(actual) < m {
			m = len(actual)
		}
	}
	for i := 0; i < m; i++ {
		if exp[i] != actual[i] {
			error = true
			break
		}
	}
	if error {
		t.Errorf("content error: %v (%v)", actual, exp)
	}
}

func assertWords(t *testing.T, exp ...string) {
	assertContent(t, exp, words)
}

func assertFlags(t *testing.T, exp ...string) {
	assertContent(t, exp, rgflags)
}

func run(t *testing.T, args ...string) {
	rgflags, words = parseArgs(args)
	t.Logf("%v => (%v, %v)\n", args, rgflags, words)
}

func Test_parseArgs(t *testing.T) {
	run(t)
	assertLengths(t, 0, 0)

	run(t, "bar", "baz")
	assertLengths(t, 0, 2)
	assertWords(t, "bar", "baz")

	run(t, "bar", "--", "baz")
	assertLengths(t, 0, 3)
	assertWords(t, "bar", "--", "baz")

	run(t, "-bar", "--", "baz")
	assertFlags(t, "-bar")
	assertWords(t, "baz")
}
