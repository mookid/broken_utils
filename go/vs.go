package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
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

func readline(reader *bufio.Reader) (string, error) {
	text, err := reader.ReadString('\n')
	if len(text) == 0 {
		return "", io.EOF
	}
	text = strings.TrimSpace(text)
	return text, err
}

func openSln(sln string) {
	cmd := exec.Command("cmd.exe", "/c", "start /b "+sln)
	cmd.Start()
	os.Exit(0)
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

func chooseFrom(slns []string) {
	signal.Reset(os.Interrupt, syscall.SIGTERM)
	if len(slns) == 0 {
		fmt.Println("*.sln: not found")
		os.Exit(2)
	}
	if len(slns) == 1 {
		openSln(slns[0])
	}
	fmt.Println(strings.Repeat("=", 72))
	fmt.Printf("open solution: ")
}

func main() {
	ctx, c := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "rg", "--files")
	results, err := cmd.StdoutPipe()
	cancel = c
	die(err)
	out := bufio.NewReader(results)
	die(cmd.Start())

	setupSignals()

	var slns []string

	go func() {
		for {
			sln, err := readline(out)
			if err == io.EOF {
				chooseFrom(slns)
				return
			}
			die(err)
			if !strings.HasSuffix(sln, "sln") {
				continue
			}
			fmt.Printf("%-8d %s\n", len(slns), sln)
			slns = append(slns, sln)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := readline(reader)
		if err != nil {
			continue
		}
		selected, err := strconv.Atoi(strings.TrimSpace(text))
		if err == nil && 0 <= selected && selected < len(slns) {
			openSln(slns[selected])
		}
	}
}
