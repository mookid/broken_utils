package main

import (
	"flag"
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
	println(`usage: clone repo-url
OPTIONS:
-d                  change destination directory name`)
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
	repoDiskName := flag.String("d", "", "")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
	}
	url := flag.Arg(0)
	rn := repoName(url)
	if *repoDiskName == "" {
		*repoDiskName = rn
	}
	dst := "D:/src/" + *repoDiskName
	if fileExists(dst) {
		fmt.Fprintf(os.Stderr, "destination file %s already exists", dst)
		os.Exit(2)
	}
	cmd := exec.Command("git", "clone", url, dst)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	die(cmd.Run())

}
