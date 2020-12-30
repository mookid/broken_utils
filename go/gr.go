package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	groups := make(map[string][]string)
	b := color.New(color.FgCyan).SprintFunc()
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
	for filename, results := range groups {
		fmt.Println(b(filename))
		for _, result := range results {
			fmt.Print(result)
		}
	}
}
