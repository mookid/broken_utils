package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func usage() {
	fmt.Printf(
		`git zip [-o archive-name]

Creates a snapshot of the checked-in files in the current subdirectory
of a git repository, in a zipped archive.

Options:
    -o  archive-name    the name of the output file (default: repo.zip)
`)
	os.Exit(2)
}

func copyfile(w *zip.Writer, relpath string) error {
	file, err := os.Open(relpath)
	if err != nil {
		return err
	}
	defer file.Close()

	f, err := w.Create(relpath)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, file)
	return err
}

func main() {
	out := flag.String("o", "repo.zip", "")
	flag.Usage = usage
	flag.Parse()

	files, err := exec.Command("git", "ls-files", "-z").Output()
	if err != nil {
		fmt.Printf("not in a git repository\n")
		os.Exit(2)
	}

	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	hi := 0
	for lo := 0; lo < len(files); lo = hi + 1 {
		for hi = lo; hi < len(files) && files[hi] != 0; hi++ {
		}
		if err := copyfile(w, string(files[lo:hi])); err != nil {
			fmt.Printf("error creating the archive: %s\n", err.Error())
			os.Exit(2)
		}
	}
	w.Close()

	if err := ioutil.WriteFile(*out, buf.Bytes(), 0644); err != nil {
		fmt.Printf("error while writing the archive: %s %s\n", *out, err.Error())
		os.Exit(2)
	}
}
