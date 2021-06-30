package main

import "testing"

func Test_hdrLength(t *testing.T) {
	tt := []struct {
		line string
		exp  int
	}{
		{"123: foo", 3},
		{"132 foo", 0},
		{"foo", 0},
	}

	for _, tc := range tt {
		t.Run(tc.line, func(t *testing.T) {
			actual := hdrLength(tc.line)
			if actual != tc.exp {
				t.Errorf("hdrLength(%s): exp: %d actual: %d", tc.line, tc.exp, actual)
			}
		})
	}
}
