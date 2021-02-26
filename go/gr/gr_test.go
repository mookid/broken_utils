package main

import "testing"

func testHdrLength1(t *testing.T, input string, exp int) {
	actual := hdrLength(input)
	if actual != exp {
		t.Errorf("hdrLength(%s): exp: %d actual: %d", input, exp, actual)
	}
}

func Test_hdrLength(t *testing.T) {
	testHdrLength1(t, "123: foo", 3)
	testHdrLength1(t, "123 foo", 0)
	testHdrLength1(t, "foo", 0)
}
