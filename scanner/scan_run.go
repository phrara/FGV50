package scanner

import (
	"encoding/json"
	"fgv50/err"
	"fgv50/flag"
	"fgv50/scanner/connection"
	"fgv50/scanner/jg"
	"fgv50/tools"
	"fgv50/tools/workerpool"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var (
	ism sync.Mutex
	pyScriptPath string
)

func init() {
	pyScriptPath, _ = os.Getwd()
	pyScriptPath = filepath.Join(pyScriptPath, "/python/main.py")
}

func initWorkPool(taskQueLen int) (*ScanTaskProcessor, *workerpool.WorkerPool) {
	scanTaskProc := NewScanTaskProcessor(taskQueLen*2 + 1)
	wp := workerpool.New(3, 2, taskQueLen+1, scanTaskProc, workerpool.SRC)
	wp.Start()
	return scanTaskProc, wp
}

func PingScreen(hosts []string, ttl int) (aliveHosts []string) {
	aliveHosts = make([]string, 0, 20)
	var wg sync.WaitGroup
	for _, v := range hosts {
		h := v
		wg.Add(1)
		go func(h string) {
			defer wg.Done()
			if b := connection.Ping(h, ttl); b {

				ism.Lock()
				aliveHosts = append(aliveHosts, h)
				ism.Unlock()
			}
		}(h)
	}
	wg.Wait()
	return aliveHosts
}

func RunCli(args *flag.Args) []byte {
	if args.Url == nil {
		// Get alive host by PingScreening
		aliveHosts := PingScreen(args.Hosts, args.TTL)
		if len(aliveHosts) == 0 {
			return nil
		}
		total := len(aliveHosts) * len(args.Ports)
		fmt.Println("total of targets are", total)
		// initiate the workpool
		processor, wp := initWorkPool(total)

		resfile := tools.Open()

		count := 0
		for _, h := range aliveHosts {
			for _, p := range args.Ports {
				scanTask := NewScanTask(h, p, "", args.TTL)
				wp.AppendTask(scanTask, 0)
				wp.AppendTask(scanTask, 1)
				//wp.AppendTask(scanTask, 2)
			}
		}
		var wg sync.WaitGroup
		resList := make([]tools.Result, 0, total*2)
		for r := range processor.GetResults() {
			res := tools.NewRes(r.Host, r.Protocol, r.Type, r.IdString, r.BString, r.Port, r.TTL, r.Buf, r.IdBool)
			fmt.Printf("%s:%d has been scanned\n", res.Host, res.Port)
			res.BString = tools.ByteToStringParse1(res.Buf)

			wg.Add(1)
			switch res.Type {
			case "tcp":
				go func() {
					defer wg.Done()

					ok := jg.IdentifyTcp(res)
					if ok && res.Protocol != "" {

						res.Time = time.Now().Format("2006-01-02 15:04:05")
						//fmt.Println(res)
						resList = append(resList, *res)
					}
				}()

			case "tls":
				go func() {
					defer wg.Done()

					ok := jg.IdentifyTls(res)
					if ok && res.Protocol != "" {

						res.Time = time.Now().Format("2006-01-02 15:04:05")
						//fmt.Println(res)
						resList = append(resList, *res)
					}
				}()
			case "udp":
				go func() {
					defer wg.Done()

					ok := jg.IdentifyUdp(res)
					if ok && res.Protocol != "" {

						res.Time = time.Now().Format("2006-01-02 15:04:05")
						//fmt.Println(res)
						resList = append(resList, *res)
					}
				}()
			}

			count++

			if count >= total*2 {
				// close result channel
				processor.CloseResC()
			}
		}
		wg.Wait()
		wp.Shut()

		resJson, _ := json.Marshal(resList)

		// write resJson into res.json
		tools.Write(resJson, resfile)
		tools.Close(resfile)

		// start up the spider to relate the info to vuls
		// write vulJson to ali_cve.json
		py := exec.Command("python3", pyScriptPath)
		err1 := py.Run()
		if err1 != nil {
			fmt.Println(fmt.Errorf("%s: %s", err.ErrRunPython, err1))
		}

		// store to histDB
		// get a timestamp as historical index
		kTime := []byte(time.Now().Format("2006-01-02 15:04:05"))
		var kTime1 []byte = []byte(string(kTime))  
		// read vulJson
		vj := tools.ReadVulJson()
		if args.HistDB != nil {
			args.HistDB.PutVulRecord(kTime, vj)
			args.HistDB.PutResRecord(kTime1, resJson)
		}

		return kTime

	} else {
		
		r := &tools.Result{}
		resfile := tools.Open()
		resList := make([]tools.Result, 0, 1)
		switch args.Url.Scheme {
		case "http":
			p := 0
			if args.Url.Port() == "" {
				p = 80
			} else {
				p, _ = strconv.Atoi(args.Url.Port())
			}
			h := args.Url.Hostname()
			b, _ := connection.TcpProtocol(h, p, args.TTL)
			r = tools.NewRes(h, "http", "tcp", "", "", p, args.TTL, b, false)
			jg.TcpHTTP(r, args.Url.String())
		case "https":
			p := 0
			if args.Url.Port() == "" {
				p = 443
			} else {
				p, _ = strconv.Atoi(args.Url.Port())
			}
			h := args.Url.Hostname()
			b, _ := connection.TlsProtocol(h, p, args.TTL)
			r = tools.NewRes(h, "https", "tls", "", "", p, args.TTL, b, false)
			jg.TcpHTTP(r, args.Url.String())
		}
		
		
		resList = append(resList, *r)

		resJson, _ := json.Marshal(resList)

		// write resJson into res.json
		tools.Write(resJson, resfile)
		tools.Close(resfile)

		// start up the spider to relate the info to vuls
		// write vulJson to ali_cve.json
		py := exec.Command("python3", pyScriptPath)
		err1 := py.Run()
		if err1 != nil {
			fmt.Println(fmt.Errorf("%s: %s", err.ErrRunPython, err1))
		}

		kTime := []byte(time.Now().Format("2006-01-02 15:04:05"))
		var kTime1 []byte = []byte(string(kTime))  
		// read vulJson
		vj := tools.ReadVulJson()
		if args.HistDB != nil {
			args.HistDB.PutVulRecord(kTime, vj)
			args.HistDB.PutResRecord(kTime1, resJson)
		}
		return kTime
	}

}
