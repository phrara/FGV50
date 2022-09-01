package main

import (
	"fgv50/flag"
	"fgv50/scanner"
	"fgv50/tools/storage"
	"fgv50/web"
	"fmt"
)

func main() {
	args, w, hdb := flag.FlagParse()
	if w {
		// open history db
		hd, err := storage.NewHistDB()
		if err != nil {
			fmt.Println(err)
			return
		}
		// TODO: start web server
		web.StartUpWeb(hd)

	} else {
		if hdb {
			// open history db
			hd, err := storage.NewHistDB()
			if err != nil {
				fmt.Println(err)
				return
			}
			args.HistDB = hd
		}
		// start up scanning	
		b := scanner.RunCli(args)
		fmt.Println("execute successfully:", string(b))
		return		
	}

}