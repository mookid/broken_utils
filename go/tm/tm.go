package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	now := time.Now()
	for _, locstring := range []string{
		"America/Vancouver",
		"America/New_York",
		"Asia/Hong_Kong",
		"Europe/Paris",
	} {
		loc, err := time.LoadLocation(locstring)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while loading location %s: %v\n", loc, err)
		}
		fmt.Printf("%20v %v\n", loc, now.In(loc).Format("15:04"))
	}
}
