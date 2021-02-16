package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type kvp struct {
	key, value string
}

func runNetsh(args ...string) ([]kvp, error) {
	command := strings.Fields("netsh wlan show profile")
	args = append(command[1:], args...)
	output, err := exec.Command(command[0], args...).Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(output), "\n")
	var results []kvp
	for _, line := range lines {
		if !strings.Contains(line, ":") {
			continue
		}
		fields := strings.SplitN(line, ":", 2)
		if len(fields) != 2 {
			continue
		}

		results = append(results, kvp{key: strings.TrimSpace(fields[0]), value: strings.TrimSpace(fields[1])})
	}

	return results, nil
}

func netshListNets() []string {
	kvps, err := runNetsh()
	if err != nil {
		fmt.Fprintf(os.Stderr, "exec %v\n", err)
		os.Exit(2)
	}

	var nets []string
	for _, kvp := range kvps {
		if kvp.key == "All User Profile" {
			nets = append(nets, kvp.value)
		}
	}
	return nets
}

func netshGetPasswd(net string) string {
	kvps, err := runNetsh("name="+net, "key=clear")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while reading password of net %s: %v\n", net, err)
		return ""
	}
	for _, kvp := range kvps {
		if kvp.key == "Key Content" {
			return kvp.value
		}
	}
	return ""
}

func main() {
	var wg sync.WaitGroup
	results := make(chan string)
	nets := netshListNets()

	for _, net := range nets {
		wg.Add(1)
		go func(net string) {
			defer wg.Done()
			passwd := netshGetPasswd(net)
			if passwd != "" {
				results <- fmt.Sprintf("%q\t%q\n", net, passwd)
			}
		}(net)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Print(result)
	}
}
