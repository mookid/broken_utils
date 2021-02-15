package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func die(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	cmd := exec.Command("rg", "^^namespace", "-g", "*.cs")
	out, err := cmd.Output()
	die(err)
	r := color.New(color.FgRed)
	g := color.New(color.FgGreen)
	for _, line := range strings.Split(string(out), "\n") {
		if line == "" {
			continue
		}
		if strings.Index(line, ":") < 0 {
			fmt.Printf("ERROR: %s does not contain :", line)
			continue
		}
		expected := filepath.Dir(line[:strings.Index(line, ":")])
		expected = strings.ReplaceAll(expected, "\\", "/")
		expected = strings.ReplaceAll(expected, "/", ".")
		expected = strings.TrimPrefix(expected, "src.tests.")
		expected = strings.TrimPrefix(expected, "src.")
		expected = strings.TrimSuffix(expected, ".cs")
		expected = fmt.Sprintf("Microsoft.%s", expected)
		expected = strings.TrimSpace(expected)
		actual := line[strings.Index(line, "namespace")+len("namespace"):]
		actual = strings.TrimSpace(actual)

		if actual != expected {
			fmt.Printf("WARN: inconsistent namespace: %s\n", line)
			r.Printf("\t%-20s\t%s\n", "actual", actual)
			g.Printf("\t%-20s\t%s\n", "expected", expected)
		}
	}
}
