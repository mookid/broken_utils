package main

import "testing"

func Test_regex(t *testing.T) {
	tt := []struct {
		line       string
		doMatch    bool
		methodName string
	}{
		{line: "public async Task FooAsync(", doMatch: true, methodName: "FooAsync"},
		{line: "private async Task FooAsync(", doMatch: true, methodName: "FooAsync"},
		{line: "protected async Task FooAsync(", doMatch: true, methodName: "FooAsync"},
		{line: "protected async Task Foo(", doMatch: true, methodName: "Foo"},
		{line: "async Task FooAsync(", doMatch: false},
		{line: "provate async Task FooAsync(", doMatch: false},
		{line: "provate async Task FooAsync(", doMatch: false},
		{line: `public const string TaskName = "Foo"`, doMatch: false},
	}

	for _, tc := range tt {
		matches := r.FindStringSubmatch(tc.line)
		if (len(matches) >= 2) != tc.doMatch {
			m := "<2"
			if tc.doMatch {
				m = ">=2"
			}
			t.Errorf("line: %s exp: %s actual: %v", tc.line, m, len(matches))

			if tc.doMatch && tc.methodName != matches[2] {
				t.Errorf("line: %s exp: %s actual: %v", tc.line, tc.methodName, matches[2])
			}
		}
	}

}
