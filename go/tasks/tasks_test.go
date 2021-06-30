package main

import (
	"testing"
)

func Test_regex(t *testing.T) {
	tt := []struct {
		name       string
		line       string
		doMatch    bool
		methodName string
	}{
		{name: "public", line: "public async Task FooAsync(", doMatch: true, methodName: "FooAsync"},
		{name: "private", line: "private async Task FooAsync(", doMatch: true, methodName: "FooAsync"},
		{name: "protected", line: "protected async Task FooAsync(", doMatch: true, methodName: "FooAsync"},
		{name: "no async", line: "protected async Task Foo(", doMatch: true, methodName: "Foo"},
		{name: "no visibility modifier", line: "async Task FooAsync(", doMatch: false},
		{name: "no visibility modifier 2", line: "provate async Task FooAsync(", doMatch: false},
		{name: "no parenthesis", line: `public const string TaskName = "Foo"`, doMatch: false},
		{name: "generics", line: `public static async Task<F<T>> Foo<T>(this Task<T> t)`, doMatch: true, methodName: "Foo"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			matches := r.FindStringSubmatch(tc.line)
			if (len(matches) >= 2) != tc.doMatch {
				m := "<2"
				if tc.doMatch {
					m = ">=2"
				}
				t.Errorf("line: %s exp: %s actual: %v", tc.line, m, len(matches))
			}

			if tc.doMatch && tc.methodName != matches[2] {
				t.Errorf("line: %s exp: %s actual: %v", tc.line, tc.methodName, matches[2])
			}
		})
	}

}
