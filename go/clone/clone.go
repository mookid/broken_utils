package main

import (
	"fmt"
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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	url := os.Args[1]
	rn := repoName(url)
	dst := "D:/src/" + rn
	if fileExists(dst) {
		fmt.Fprintf(os.Stderr, "destination file %s already exists", dst)
		os.Exit(2)
	}
	cmd := exec.Command("git", "clone", url, dst)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	die(cmd.Run())

}
