package main

import "testing"

func Test_parse(t *testing.T) {
	tt := []struct {
		name  string
		input string
		exp   proc
	}{
		{"normal2", `"field1","field2"`, proc{"field1", "field2"}},
		{"normal3", `"field1","field2","field3"`, proc{"field1", "field2"}},
		{"with comma in string", `"fiel,d1","field2""`, proc{"fiel,d1", "field2"}},
		{"with escaping", `"field1","fi\"eld2""`, proc{"field1", `fi"eld2`}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := parse(tc.input)
			if err != nil {
				t.Error("unexpected error")
			}
			if actual.name != tc.exp.name {
				t.Errorf("name mismatch: '%s' '%s'", actual.name, tc.exp.name)
			}
			if actual.pid != tc.exp.pid {
				t.Errorf("pid mismatch: '%s' '%s'", actual.pid, tc.exp.pid)
			}
		})
	}
}
