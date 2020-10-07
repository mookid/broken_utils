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

func doKill(pid string) {
	cmd := exec.Command("taskkill.exe", "/pid", pid, "/f")
	cmd.Stdin, cmd.Stdout = os.Stdin, os.Stdout
	die(cmd.Run())
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

func chooseFrom(pids []string) {
	signal.Reset(os.Interrupt, syscall.SIGTERM)
	if len(pids) == 0 {
		fmt.Println("process not found")
		os.Exit(2)
	}
	fmt.Println(strings.Repeat("=", 72))
	fmt.Printf("kill process: ")
}

type proc = struct{ name, pid string }

func unquote(str string) string {
	str = strings.TrimRight(str, "\"")
	str = strings.TrimLeft(str, "\"")
	return str
}

func parse(str string) (proc, error) {
	// TODO proper parsing of "fi,eld2","field2"
	fields := strings.Split(str, ",")
	return proc{unquote(fields[0]), unquote(fields[1])}, nil
}

func main() {
	withFilter := false
	ctx, c := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "tasklist.exe", "/nh", "/fo", "csv")
	results, err := cmd.StdoutPipe()
	die(err)

	if len(os.Args) > 1 {
		withFilter = true
		args := append([]string{"-i"}, os.Args[1:]...)
		cmd2 := exec.CommandContext(ctx, "rg", args...)
		cmd2.Stdin = results
		results, err = cmd2.StdoutPipe()
		die(err)
		die(cmd2.Start())
	}
	cancel = c
	out := bufio.NewReader(results)
	die(cmd.Start())

	setupSignals()

	var pids []string

	go func() {
		for {
			procString, err := readline(out)
			if err == io.EOF {
				chooseFrom(pids)
				return
			}
			proc, err := parse(procString)
			if err != nil {
				fmt.Fprintln(os.Stderr, "parsing error: '%s'", procString)
				continue
			}
			fmt.Printf("%-8d %s\n", len(pids), proc.name)
			pids = append(pids, proc.pid)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := readline(reader)
		if err != nil {
			continue
		}
		text = strings.TrimSpace(text)
		if text == "*" && withFilter {
			for _, pid := range pids {
				doKill(pid)
			}
			os.Exit(0)
		}
		selected, err := strconv.Atoi(text)
		if err == nil && 0 <= selected && selected < len(pids) {
			doKill(pids[selected])
			os.Exit(0)
		}
	}
}
