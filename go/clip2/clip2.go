package main

import (
	"io/ioutil"
	"os"

	"github.com/atotto/clipboard"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	if err := clipboard.WriteAll(string(data)); err != nil {
		panic(err)
	}
}
