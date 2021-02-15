package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	groups map[string][]string = make(map[string][]string)
)

func readFile(filename string) {
	var (
		f   io.Reader
		err error
	)
	if filename == "-" {
		f = os.Stdin
	} else {
		f, err = os.Open(filename)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read file %s: %v\n", filename, err)
	}

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "failed to read from stdin: %v\n", err)
			}
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "invalid line: %v\n", parts)
			continue
		}

		groups[parts[0]] = append(groups[parts[0]], parts[1])
	}
}

func hdrLength(fileMatch string) int {
	hdrLength := 0
	inHdr := true
	for _, ch := range fileMatch {
		switch {
		case ch == ':':
			inHdr = false
		case '0' <= ch && ch <= '9':
			hdrLength++
		default:
			hdrLength = 0
			inHdr = false
		}
		if !inHdr {
			break
		}
	}
	return hdrLength
}

func main() {
	fromStdin := true
	for _, arg := range os.Args[1:] {
		fromStdin = false
		readFile(arg)
	}
	if fromStdin {
		readFile("-")
	}
	b := color.New(color.FgCyan)
	g := color.New(color.FgGreen)
	for filename, results := range groups {
		b.Println(filename)
		for _, result := range results {
			n := hdrLength(result)
			g.Print(result[:n])
			fmt.Print(result[n:])
		}
	}
}
