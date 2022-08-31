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
	-H 123 (default -1)
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



func FlagParse() (*Args, bool) {
	var networkSegment string
	var url string
	var hisID int
	var ip string
	var port int
	var TTL int
	flag.StringVar(&networkSegment, "ns", "", "-ns 123.123.123.1~10 or -ns 123.123.123.1+9")
	flag.StringVar(&url, "u", "", "-u http://123.123.123.123:80")
	flag.IntVar(&hisID, "H", -1, "-H 123")
	flag.StringVar(&ip, "i", "", "-i 123.123.123.123")
	flag.IntVar(&port, "p", -1, "-p 3306")
	flag.IntVar(&TTL, "t", 3, "-t 3")
	flag.Parse() 
	
	if flag.NFlag() == 0 {
		fmt.Println(err.ErrArgsLack)
		return nil, false
	}
 	if flag.NArg() == 0 {
		if (networkSegment != "" && url != "") || (url != "" && ip != "") || (url != "" && port != -1) {
			fmt.Println(err.ErrArgsConflict)
			return nil, false
		}
		if args, err := NewArgs(url, networkSegment, ip, port, TTL); err == nil {
			return args, false
		} else {
			fmt.Println(err)
			return nil, false
		}

	}
	arg := flag.Args()[0]
	switch arg {
	case "history":
		return nil, false
	case "webui":
		return nil, true
	default:
		fmt.Println(err.ErrUnknownArgs)
		return nil, false
	}

}