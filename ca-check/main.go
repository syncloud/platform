package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	dir := flag.String("certs", "", "directory of shipped root certificates")
	years := flag.Int("years", 3, "minimum remaining validity of the long-lived part of the chain")
	flag.Parse()

	checker, err := NewChecker(*dir, time.Now().AddDate(*years, 0, 0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	failed := false
	for _, host := range flag.Args() {
		err = checker.Check(host)
		if err != nil {
			fmt.Println(err)
			failed = true
		}
	}
	if failed || flag.NArg() == 0 {
		os.Exit(1)
	}
}
