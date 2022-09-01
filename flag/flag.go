package flag

import (
	"fgv50/err"
	"flag"
	"fmt"
)
const (
	help = `
WITH FLAG:
-H int
	-H 1 (default 0) 
	when -H 0: Do not store this scan result, otherwise, store it.
-i string
	-i 123.123.123.123
-ns string
	-ns 123.123.123.1~10 or -ns 123.123.123.1+9
-p int
	-p 3306 (default -1)
-u string
	-u http://123.123.123.123:80
-t int
	-t 3 (default 3)
WITHOUT FLAG:
	history
		View historical scan records
	webui
		Start the front-end UI
	`
)

func init() {
	flag.Usage = func() {
		fmt.Println(help)
	}
}



func FlagParse() (args *Args, web bool, hdb bool) {
	var networkSegment string
	var url string
	var histDB int
	var ip string
	var port int
	var TTL int
	flag.StringVar(&networkSegment, "ns", "", "-ns 123.123.123.1~10 or -ns 123.123.123.1+9")
	flag.StringVar(&url, "u", "", "-u http://123.123.123.123:80")
	flag.IntVar(&histDB, "H", 0, "-H 1")
	flag.StringVar(&ip, "i", "", "-i 123.123.123.123")
	flag.IntVar(&port, "p", -1, "-p 3306")
	flag.IntVar(&TTL, "t", 3, "-t 3")
	flag.Parse() 
	
	if flag.NFlag() == 0 && flag.NArg() == 0 {
		fmt.Println(err.ErrArgsLack)
		return nil, false, false
	}
 	if flag.NArg() == 0 {
		if (networkSegment != "" && url != "") || (url != "" && ip != "") || (url != "" && port != -1) {
			fmt.Println(err.ErrArgsConflict)
			return nil, false, false
		}
		if argus, err := NewArgs(url, networkSegment, ip, port, TTL); err == nil {
			if histDB == 0 {
				return argus, false, false
			} else {
				return argus, false, true
			}
			
		} else {
			fmt.Println(err)
			return nil, false, false
		}

	}
	arg := flag.Args()[0]
	switch arg {
	case "history":
		return nil, false, false
	case "webui":
		return nil, true, true
	default:
		fmt.Println(err.ErrUnknownArgs)
		return nil, false, false
	}

}