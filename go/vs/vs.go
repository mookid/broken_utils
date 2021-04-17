package main

import (
	"bufio"
	"context"
	"flag"
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

func usage() {
	println(`usage: vs [options...] patterns...
Options:
  -k       kill existing devenv instances
`)
	os.Exit(2)
}

func doKillExisting() (err error) {
	cmd := exec.Command("re-kill", "-a", "devenv.exe")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	cmd.Wait()
	return nil
}

func main() {
	killExisting := flag.Bool("k", false, "")
	flag.Usage = usage
	flag.Parse()
	patterns := flag.Args()

	if *killExisting {
		die(doKillExisting())
	}

	ctx, c := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "rg", "--files", "--no-ignore", "--hidden")
	results, err := cmd.StdoutPipe()

	if len(patterns) != 0 {
		args := append([]string{"-i", "-M0"}, patterns...)
		cmd2 := exec.CommandContext(ctx, "rg", args...)
		cmd2.Stderr = os.Stderr
		cmd2.Stdin = results
		results, err = cmd2.StdoutPipe()
		die(cmd2.Start())
	}
	cancel = c
	die(err)
	out := bufio.NewReader(results)
	die(cmd.Start())

	setupSignals()

	var slns []string

	go func() {
		firstLine := ""
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
			line := fmt.Sprintf("%-8d %s", len(slns), sln)
			if len(slns) == 0 {
				firstLine = line
			} else if len(slns) == 1 {
				fmt.Println(firstLine)
				fmt.Println(line)
			} else {
				fmt.Println(line)
			}
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
