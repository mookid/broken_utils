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

func main() {
	cmd := exec.Command("rg", "--files")
	out, err := cmd.Output()
	die(err)
	for _, line := range strings.Split(string(out), "\n") {
		if line == "" {
			continue
		}
		if strings.HasSuffix(line, "orig") ||
			strings.Contains(line, "_BACKUP_") ||
			strings.Contains(line, "_BASE_") ||
			strings.Contains(line, "_LOCAL_") ||
			strings.Contains(line, "_REMOTE_") {
			fmt.Printf("removing file %s...\n", line)
			os.Remove(line)
		}
	}
}
