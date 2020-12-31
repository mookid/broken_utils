package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func die(err error) {
	if err != nil {
		panic(err)
	}
}

func usage() {
	println("git ff [<pattern>]")
	os.Exit(2)
}

func readline(reader *bufio.Reader) (string, error) {
	text, err := reader.ReadString('\n')
	if len(text) == 0 {
		return "", io.EOF
	}
	text = strings.TrimSpace(text)
	return text, err
}

func main() {
	if len(os.Args) > 2 {
		usage()
	}
	cmd := exec.Command("rg", "--files", "--hidden")
	cmd.Stderr = os.Stderr
	results, err := cmd.StdoutPipe()
	die(err)

	if len(os.Args) == 2 {
		cmd2 := exec.Command("rg", "-i", os.Args[1])
		cmd2.Stderr = os.Stderr
		cmd2.Stdin = results
		results, err = cmd2.StdoutPipe()
		die(cmd2.Start())
	}
	out := bufio.NewReader(results)
	die(cmd.Start())

	for {
		line, err := readline(out)
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		if line == "" {
			continue
		}
		fmt.Println(line)
	}
}
