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
	scanTask := NewScanTaskProcessor(DefaultResQueLen)
	wp := workerpool.New(4, 2, 20, scanTask, workerpool.RR)
	wp.Start()
	return scanTask, wp
}

func PingScreen(hosts []string) (aliveHosts []string) {
	aliveHosts = make([]string, 0, 20)
	var wg sync.WaitGroup
	for _, v := range hosts {	
		h := v
		wg.Add(1)
		go func (h string)  {
			defer wg.Done()
			if b := connection.Ping(h); b {
				
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
		aliveHosts := PingScreen(args.Hosts)
		total := len(aliveHosts) * len(args.Ports)
		count := 0
		for _, h := range aliveHosts {
			for _, p := range args.Ports {
				scanTask := NewScanTask(h, p, "", args.TTL)
				wp.AppendTask(scanTask, 0)
			}
		}
		for r := range processor.GetResults() {
			fmt.Println(r.Buf)
			count++
			if count >= total {
				// TODO: Judge The Result 

				// close result channel
				processor.CloseResC()
			}
		}
	} else {
		// TODO
		panic("todo")
	}

}


