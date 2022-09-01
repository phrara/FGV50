package main

import (
	"fgv50/flag"
	"fgv50/scanner"
	"fgv50/tools/storage"
	"fgv50/web"
	"fmt"
)

const LOGO = `
'        ___           ___           ___       ___       ___     
'       /\__\         /\  \         /\__\     /\__\     /\  \    
'      /:/  /        /::\  \       /:/  /    /:/  /    /::\  \   
'     /:/__/        /:/\:\  \     /:/  /    /:/  /    /:/\:\  \  
'    /::\  \ ___   /::\~\:\  \   /:/  /    /:/  /    /:/  \:\  \ 
'   /:/\:\  /\__\ /:/\:\ \:\__\ /:/__/    /:/__/    /:/__/ \:\__\
'   \/__\:\/:/  / \:\~\:\ \/__/ \:\  \    \:\  \    \:\  \ /:/  /
'        \::/  /   \:\ \:\__\    \:\  \    \:\  \    \:\  /:/  / 
'        /:/  /     \:\ \/__/     \:\  \    \:\  \    \:\/:/  /  
'       /:/  /       \:\__\        \:\__\    \:\__\    \::/  /   
'       \/__/         \/__/         \/__/     \/__/     \/__/    `


func main() {
	fmt.Printf("\u001B[1;35m%s\u001B[0m\n", LOGO)
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
		if args.HistDB != nil {
			args.HistDB.Close()
		}
		fmt.Println("execute successfully:", string(b))
		return		
	}

}