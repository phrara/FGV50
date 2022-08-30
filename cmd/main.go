package main

import (
	"csl/flag"
	"csl/scanner"
	"fmt"
)

func main() {
	if args, b := flag.FlagParse(); b {
		fmt.Println("web ui is prepared")
	} else {
		scanner.RunCli(args)
		fmt.Println("execute successfully")
		return
	}

}