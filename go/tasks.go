package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func die(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	regex := "(public|private|protected).*Task.* (\\w+)(<\\w+>)?\\("
	r := regexp.MustCompile(regex)
	cmd := exec.Command("rg", "-n", regex, "-g", "*.cs", "-g", "!src/tests/*", "-M0")
	out, err := cmd.Output()
	die(err)
	for _, line := range strings.Split(string(out), "\n") {
		if line == "" {
			continue
		}
		matches := r.FindStringSubmatch(line)
		if len(matches) < 2 {
			fmt.Printf("WARN: matches: %s (%d)\n", line, len(matches))
			continue
		}
		methodName := matches[2]
		if strings.HasSuffix(methodName, "Async") ||
			methodName == "Main" ||
			methodName == "TimeoutAfter" {
			continue
		}
		fmt.Printf("missing Async suffix: %s", line)
	}
}
