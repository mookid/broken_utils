package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var code int

func sha256sumDir(dirname string) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		failure(err)
		return
	}
	for _, file := range files {
		sha256sum(filepath.Join(dirname, file.Name()))
	}
}

func failure(err error) {
	fmt.Fprintf(os.Stderr, "%v", err)
	code = 2
}

func sha256sum(filename string) {
	var (
		f    io.Reader
		data []byte
		err  error
	)
	if filename == "-" {
		f = os.Stdin
	} else {
		fileInfo, err := os.Stat(filename)
		if err == nil {
			if fileInfo.IsDir() {
				sha256sumDir(filename)
				return
			}
			f, err = os.Open(filename)
		}
	}
	if err == nil {
		data, err = ioutil.ReadAll(f)
	}
	if err != nil {
		failure(err)
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
