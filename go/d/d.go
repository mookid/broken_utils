package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func exitLikeChild(err error) {
	exiterr, ok := err.(*exec.ExitError)
	if ok {
		os.Exit(exiterr.ExitCode())
	}
	if err != nil {
		os.Exit(2)
	}
}

func usage() {
	fmt.Println(`d [-c] [arg...]
Forward to git diff.
OPTIONS:
  -c                  View staged changes (--cached)
  -m                  View diff from origin/master
  args                Other args are forwarded to git diff`)
	os.Exit(2)
}

func main() {
	args := []string{"diff", "--patch-with-stat", "--stat-width=1000"}
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-c":
			arg = "--cached"
		case "-m":
			arg = "origin/master..."
		case "-h", "--help":
			usage()
		}
		if strings.HasPrefix(arg, "-") {
			i, err := strconv.Atoi(arg[1:])
			if err == nil && i > 0 {
				arg = fmt.Sprintf("HEAD~%d..", i)
			}
		}
		args = append(args, arg)
	}
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	exitLikeChild(cmd.Run())
}
