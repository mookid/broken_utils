package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
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

	if *nflag {
		bf := bufio.NewReader(f)
		for {
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
			fmt.Printf("%6d\t", n)
			os.Stdout.Write(line)
			n++
		}
	} else {
		io.Copy(os.Stdout, f)
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
