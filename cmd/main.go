package main

import (
	"csl/flag"
	"fmt"
)

func main() {
	if args, b := flag.FlagParse(); b {
		fmt.Println("web ui is prepared")
	} else {
		if args != nil {
			//fmt.Printf("args are %v", *args)
			args.PrintAllAddrs()
		}
		fmt.Println("execute successfully")
		return
	}

}