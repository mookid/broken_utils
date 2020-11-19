package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var (
	regex string         = "(public|private|protected).*Task[=]* (\\w+)(<\\w+>)?\\("
	r     *regexp.Regexp = regexp.MustCompile(regex)
)

func die(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	out, err := exec.Command("rg", "-n", regex, "-g", "*.cs", "-g", "!src/tests/*", "-M0").Output()
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
