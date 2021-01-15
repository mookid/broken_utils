package main

import (
	"fmt"
	"time"
)

func die(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	now := time.Now()
	for _, locstring := range []string{
		"America/Vancouver",
		"America/New_York",
		"Asia/Hong_Kong",
		"Europe/Paris",
	} {
		loc, err := time.LoadLocation(locstring)
		die(err)
		fmt.Printf("%20v %v\n", loc, now.In(loc).Format("15:04"))
	}
}
