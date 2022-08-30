package scanner

import (
	"csl/flag"
	"csl/scanner/connection"
	"csl/tools/workerpool"
	"fmt"
	"sync"
)

var (
	ism sync.Mutex
)

func initWorkPool() (*ScanTaskProcessor, *workerpool.WorkerPool) {
	scanTaskProc := NewScanTaskProcessor(DefaultResQueLen)
	wp := workerpool.New(4, 2, 20, scanTaskProc, workerpool.RR)
	wp.Start()
	return scanTaskProc, wp
}

func PingScreen(hosts []string, ttl int) (aliveHosts []string) {
	aliveHosts = make([]string, 0, 20)
	var wg sync.WaitGroup
	for _, v := range hosts {	
		h := v
		wg.Add(1)
		go func (h string)  {
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

func RunCli(args *flag.Args) {
	if args.Url == nil {
		// initiate the workpool
		processor, wp := initWorkPool()
		aliveHosts := PingScreen(args.Hosts, args.TTL)
		if len(aliveHosts) == 0 {
			return
		}
		total := len(aliveHosts) * len(args.Ports)
		fmt.Println("total of targets are", total)
		count := 0
		for _, h := range aliveHosts {
			for _, p := range args.Ports {
				scanTask := NewScanTask(h, p, "", args.TTL)
				wp.AppendTask(scanTask, 0)
			}
		}
		for r := range processor.GetResults() {
			fmt.Println(string(r.Buf))
			count++
			if count >= total * 3 {
				// TODO: Judge The Result 

				// close result channel
				processor.CloseResC()
			}
		}
		wp.Shut()
	} else {
		// TODO
		panic("todo")
	}

}


