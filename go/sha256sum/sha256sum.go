package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var code int

func sha256sum(filename string) {
	var (
		f    io.Reader
		data []byte
		err  error
	)
	if filename == "-" {
		f = os.Stdin
	} else {
		f, err = os.Open(filename)
	}
	if err == nil {
		data, err = ioutil.ReadAll(f)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		code = 2
		return
	}
	sum := sha256.Sum256(data)
	fmt.Printf("%x %s\n", sum, filename)
}

func main() {
	// TODO help?
	// TODO -c?
	// TODO -s?
	fromStdin := true
	for _, filename := range os.Args[1:] {
		fromStdin = false
		sha256sum(filename)
	}
	if fromStdin {
		sha256sum("-")
	}
	os.Exit(code)
}
