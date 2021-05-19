package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
)

var (
	nflag = flag.Bool("n", false, "")
	n     = 1
	code  = 0
)

func usage() {
	fmt.Println(`cat [-n] file...`)
	os.Exit(2)
}

// hack to detect ctrl+D on windows
func isEOFHack(line []byte) bool {
	return runtime.GOOS == "windows" && bytes.Equal(line, []byte{4, 13, 10})
}

func cat(filename string) {
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
	for {
		line, err := bf.ReadBytes('\n')
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "read %s: %v\n", filename, err)
			code = 2
			return
		}
		if len(line) == 0 || isEOFHack(line) {
			return
		}
		if *nflag {
			fmt.Printf("%6d\t", n)
			n++
		}
		os.Stdout.Write(line)
		if len(line) == 0 || line[len(line)-1] != '\n' {
			os.Stdout.Write([]byte("\n"))
		}
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()

	any := false
	for _, arg := range flag.Args() {
		any = true
		cat(arg)
	}
	if !any {
		cat("-")
	}
	os.Exit(code)
}
