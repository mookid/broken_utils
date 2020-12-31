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

func usage() {
	println("l [filename]")
	os.Exit(2)
}

func main() {
	if len(os.Args) > 2 {
		usage()
	}

	flags := map[string][]string{
		"ls":   {"-lH"},
		"less": {"-RFMi"},
	}
	prog := "ls"
	filename := ""

	if len(os.Args) == 2 {
		filename = os.Args[1]
		file, err := os.Open(filename)
		if err != nil {
			println(err.Error())
			os.Exit(2)
		}
		fileinfo, err := file.Stat()
		die(err)
		if !fileinfo.IsDir() {
			prog = "less"
		}
	}
	args := flags[prog]
	if filename != "" {
		args = append(args, filename)
	}
	cmd := exec.Command(prog, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	die(cmd.Run())
}
