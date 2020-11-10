package main

import (
	"os"
	"os/exec"
	"strings"
)

func usage() {
	println("rgw [flags...] -- words...")
	os.Exit(2)
}

func exit_like_child(err error) {
	exiterr, ok := err.(*exec.ExitError)
	if ok {
		os.Exit(exiterr.ExitCode())
	}
	if err != nil {
		os.Exit(2)
	}
}

func parseArgs(args []string) (rgflags, words []string) {
	if len(args) == 0 {
		return rgflags, words
	}
	if !strings.HasPrefix(args[0], "-") {
		words = args[:]
		return rgflags, words
	}

	for i, arg := range args {
		if arg == "--" {
			words = args[i+1:]
			break
		}
		rgflags = append(rgflags, arg)
	}
	return rgflags, words
}

func main() {
	if len(os.Args) == 1 {
		usage()
	}

	// TODO handle rg's PATH?
	rgflags, words := parseArgs(os.Args[1:])
	rgflags = append(rgflags, strings.Join(words, ".*"))
	cmd := exec.Command("rg", rgflags...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	exit_like_child(cmd.Run())
}
