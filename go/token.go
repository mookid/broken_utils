package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/atotto/clipboard"
)

func die(err error) {
	if err != nil {
		panic(err)
	}
}

func process(txt string, v bool) string {
	re := regexp.MustCompile("access_token=(.*)&token_type=")
	matches := re.FindStringSubmatch(txt)
	if len(matches) == 0 {
		println("not found")
		os.Exit(2)
	}
	txt = matches[1]
	if v {
		println(txt)
	}
	return txt
}

func main() {
	v := flag.Bool("v", false, "verbose")
	flag.Parse()
	txt, err := clipboard.ReadAll()
	die(err)
	txt = process(txt, *v)
	clipboard.WriteAll(txt)
	fmt.Println("token saved to clipboard")
}
