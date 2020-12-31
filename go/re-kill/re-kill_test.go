package main

import "testing"

func parse_ok(t *testing.T, actual, exp *proc, err error) {
	if err != nil {
		t.Error("unexpected error")
	}
	if actual.name != exp.name {
		t.Errorf("name mismatch: '%s' '%s'", actual.name, exp.name)
	}
	if actual.pid != exp.pid {
		t.Errorf("pid mismatch: '%s' '%s'", actual.pid, exp.pid)
	}
}

func Test_parse(t *testing.T) {
	res, err := parse(`"field1","field2"`)
	parse_ok(t, res, &proc{"field1", "field2"}, err)

	res, err = parse(`"field1","field2","field3"`)
	parse_ok(t, res, &proc{"field1", "field2"}, err)

	res, err = parse(` "fiel,d1","field2"`)
	parse_ok(t, res, &proc{"fiel,d1", "field2"}, err)
}

func Test_parseWithEscaping(t *testing.T) {
	res, err := parse(`"field1",  "fi\"eld2"`)
	parse_ok(t, res, &proc{"field1", `fi"eld2`}, err)
}
