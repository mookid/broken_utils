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

type r struct {
	data    []byte
	relpath string
}

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

func readFile(relpath string) (result *r, err error) {
	file, err := os.Open(relpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &r{data, relpath}, nil
}

func readPaths() (relpaths []string) {
	files, err := exec.Command("git", "ls-files", "-z").Output()
	if err != nil {
		fmt.Printf("not in a git repository\n")
		os.Exit(2)
	}

	hi := 0
	for lo := 0; lo < len(files); lo = hi + 1 {
		for hi = lo; hi < len(files) && files[hi] != 0; hi++ {
		}
		relpaths = append(relpaths, string(files[lo:hi]))
	}
	return relpaths
}

func readFiles(relpaths []string) (results []r) {
	ok := true
	results = make([]r, len(relpaths))
	for i, relpath := range relpaths {
		result, err := readFile(relpath)
		if err != nil {
			ok = false
			fmt.Printf("%v\n", err)
		}
		results[i] = *result
	}
	if !ok {
		os.Exit(2)
	}
	return results
}

func createArchive(results []r) []byte {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	for _, result := range results {
		f, err := w.Create(result.relpath)
		if err == nil {
			_, err = io.Copy(f, bytes.NewReader(result.data))
		}
		if err != nil {
			fmt.Printf("error creating the archive: %v\n", err)
		}
	}
	w.Close()
	return buf.Bytes()
}

func writeArchive(data []byte, out string) {
	if err := ioutil.WriteFile(out, data, 0644); err != nil {
		fmt.Printf("error while writing the archive: %s %v\n", out, err)
		os.Exit(2)
	}
}

func main() {
	out := flag.String("o", "repo.zip", "")
	flag.Usage = usage
	flag.Parse()
	relpaths := readPaths()
	results := readFiles(relpaths)
	data := createArchive(results)
	writeArchive(data, *out)
}
