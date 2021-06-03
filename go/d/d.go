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
	fmt.Println(`d [-c] [-m] [-o] [-v DIR] [arg...]
Forward to git diff.
OPTIONS:
  -c                  View staged changes (--cached)
  -m                  View diff from origin/master
  -o                  List changed files only (--name-only)
  -v DIR              Exclude the directory from the diff
  args                Other args are forwarded to git diff`)
	os.Exit(2)
}

var mainBranchName string

func inferMainBranchName() string {
	if mainBranchName == "" {
		for _, name := range []string{"origin/master", "origin/main"} {
			if exec.Command("git", "show-ref", "--verify", "--quiet", "refs/remotes/"+name).Run() == nil {
				mainBranchName = name
				break
			}
		}
	}
	return mainBranchName
}

func main() {
	argsLists := [][]string{{"diff", "--patch-with-stat", "--stat-width=1000"}, {}}
	for i := 1; i < len(os.Args); i++ {
		arg0 := os.Args[i]
		arg1 := ""
		switch arg0 {
		case "-c":
			arg0 = "--cached"
		case "-m":
			arg0 = inferMainBranchName() + "..."
		case "-o":
			arg0 = "--name-only"
		case "-v":
			i++
			arg0 = ""
			arg1 = fmt.Sprintf(":!%s", os.Args[i])
		case "-h", "--help":
			usage()
		}
		if strings.HasPrefix(arg0, "-") {
			i, err := strconv.Atoi(arg0[1:])
			if err == nil && i > 0 {
				arg0 = fmt.Sprintf("HEAD~%d..", i)
			}
		}
		argsLists[0] = append(argsLists[0], arg0)
		argsLists[1] = append(argsLists[1], arg1)
	}
	var args []string
	for _, argsList := range argsLists {
		for _, arg := range argsList {
			if arg != "" {
				args = append(args, arg)
			}
		}
	}
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	exitLikeChild(cmd.Run())
}
