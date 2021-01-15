package main

import (
	"os"
	"os/exec"
	"strings"
)

func die(err error) {
	if err != nil {
		panic(err)
	}
}

func usage() {
	println("usage: clone repo-url")
	os.Exit(2)
}

func repoName(url string) string {
	pieces := strings.Split(url, "/")

	for i := len(pieces) - 1; i >= 0; i-- {
		if pieces[i] != "" {
			return strings.TrimSuffix(pieces[i], ".git")
		}
	}

	usage()
	return "" // unreachable
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	url := os.Args[1]
	rn := repoName(url)
	cmd := exec.Command("git", "clone", url, "D:/src/"+rn)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	die(cmd.Run())

}
