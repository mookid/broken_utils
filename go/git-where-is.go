package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var (
	logformat string
)

func usage() {
	fmt.Println(`git-where-is [words...]
List branches that contains a commit matching a certain pattern.
`)
	os.Exit(2)
}

func die(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
	os.Exit(2)
}

func warn(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "WARNING: "+format+"\n", args...)
}

type branch struct {
	Objectname string
	Refname    string
}

func LogFormat() string {
	if logformat != "" {
		return logformat
	}
	logformatBytes, err := json.Marshal(branch{
		Refname:    "%(refname)",
		Objectname: "%(objectname)",
	})
	if err != nil {
		panic(err)
	}
	logformat = string(logformatBytes)
	return logformat
}

func listBranches() (branches []*branch, err error) {
	data, err := exec.Command("git", "for-each-ref", "--format="+LogFormat()).Output()
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		branch := new(branch)
		if err := json.Unmarshal([]byte(line), branch); err != nil {
			warn("error parsing branch %s: %v", line, err)
			continue
		}
		branches = append(branches, branch)
	}
	return branches, nil
}

func checkIfContains(wg *sync.WaitGroup, pattern string, branch *branch) {
	defer wg.Done()
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "git", "log", branch.Objectname, "--oneline", "--grep="+pattern)
	out, err := cmd.StdoutPipe()
	if err != nil {
		warn("failed to run git log for branch %s", branch.Refname)
		return
	}
	defer out.Close()
	if err := cmd.Start(); err != nil {
		warn("failed to run git log for branch %s", branch.Refname)
		return
	}
	defer cancel()
	bufout := bufio.NewReader(out)
	data, err := bufout.ReadString('\n')
	if err != nil {
		if err != io.EOF {
			warn("error while reading log for branch %s: %v", branch.Refname, err)
		}
		return
	}
	if data = strings.TrimSpace(data); len(data) != 0 {
		fmt.Println(branch.Refname, data)
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()
	branches, err := listBranches()
	if err != nil {
		die("error while listing branches", err)
	}
	pattern := strings.Join(flag.Args(), ".*")
	var wg sync.WaitGroup
	for _, branch := range branches {
		wg.Add(1)
		go checkIfContains(&wg, pattern, branch)
	}
	wg.Wait()
}
