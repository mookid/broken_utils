package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

var (
	cancel context.CancelFunc
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

func setupSignals() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if cancel != nil {
			cancel()
		}
	}()
}

// TODO  -v
func main() {
	if len(os.Args) > 2 {
		usage()
	}
	ctx, c := context.WithCancel(context.Background())
	cancel = c
	setupSignals()
	cmd := exec.CommandContext(ctx, "rg", "--files", "--hidden", "--path-separator=/")
	cmd.Stderr = os.Stderr
	results, err := cmd.StdoutPipe()
	die(err)

	if len(os.Args) == 2 {
		cmd2 := exec.CommandContext(ctx, "rg", "-i", os.Args[1])
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
