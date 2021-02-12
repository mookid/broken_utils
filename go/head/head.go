package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	code int
)

func usage() {
	fmt.Println(`head [-n] file`)
	os.Exit(2)
}

func head(nlines int, filename string) {
	var (
		f   io.Reader
		err error
	)
	if filename == "-" {
		f = os.Stdin
		filename = "<stdin>"
	} else {
		f, err = os.Open(filename)
	}
	if err != nil {
		fmt.Printf("WARNING: error while opening %s\n", filename)
		return
	}

	bf := bufio.NewReader(f)
	for i := 0; i < nlines; i++ {
		line, err := bf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Fprintf(os.Stderr, "read %s: %v\n", filename, err)
			code = 2
			return
		}
		if len(line) == 0 {
			return
		}
		os.Stdout.Write(line)
	}
}

func main() {
	nlines := flag.Int("n", 15, "")
	flag.Usage = usage
	flag.Parse()

	any := false
	for _, arg := range flag.Args() {
		any = true
		head(*nlines, arg)
	}
	if !any {
		head(*nlines, "-")
	}
	os.Exit(code)
}
