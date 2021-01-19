package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func die(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	cmd := exec.Command("git", "diff", "--name-only", "--diff-filter=U")
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	die(err)
	home, found := os.LookupEnv("HOME")
	if !found {
		panic("HOME variable not set")
	}
	commandArgs := []string{"runemacs.exe", "-Q", "-l", filepath.Join(home, ".emacs.d", ".emacs.git")}

	conflicts := false
	for _, line := range strings.Split(string(out), "\n") {
		if line == "" {
			continue
		}
		conflicts = true
		commandArgs = append(commandArgs, line)
	}
	if conflicts {
		command := strings.Join(commandArgs, " ")
		cmd = exec.Command("cmd.exe", "/c", command)
		cmd.Start()
	}
}
