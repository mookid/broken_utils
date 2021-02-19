package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func usage() {
	println("git codereview branch-name")
	os.Exit(2)
}

func run(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func switchTo(branchName string) bool {
	fmt.Printf("Trying to switch to branch %s\n", branchName)
	err := run("git", "switch", branchName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "git switch: %v\n", err)
	}
	return err == nil
}

func fetchAll() bool {
	fmt.Println("Fetching remote branches")
	err := run("git", "fetch", "--all")
	if err != nil {
		fmt.Fprintf(os.Stderr, "git fetch: %v\n", err)
	}
	return err == nil
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		usage()
	}
	branchName := flag.Arg(0)

	if switchTo(branchName) {
		return
	}
	if !fetchAll() {
		os.Exit(2)
	}
	if switchTo(branchName) {
		return
	}
	os.Exit(2)
}
