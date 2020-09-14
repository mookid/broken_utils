package main

import (
	"os"
	"os/exec"
)

func die(err error) {
	if err != nil {
		panic(err)
	}
}

func run(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	if run("git", "status") != nil {
		die(run("ls", "-lH"))
	}
}
