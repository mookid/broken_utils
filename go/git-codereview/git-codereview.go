package main

import (
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

func switch_to(branch_name string) {
	fmt.Printf("Trying to switch to branch %s\n", branch_name)
	if err := run("git", "switch", branch_name); err == nil {
		os.Exit(0)
	}
}

func fetch_all() {
	fmt.Println("Fetching remote branches")
	if err := run("git", "fetch", "--all"); err != nil {
		os.Exit(2)
	}
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	branch_name := os.Args[1]

	switch_to(branch_name)
	fetch_all()
	switch_to(branch_name)
	os.Exit(2)
}
