package scanner

import (
	"fgv50/tools"
	"fgv50/flag"
	"fgv50/scanner/connection"
	"fgv50/tools/workerpool"
	"fmt"
	"sync"
)

var (
	ism sync.Mutex
)

func initWorkPool(taskQueLen int) (*ScanTaskProcessor, *workerpool.WorkerPool) {
	scanTaskProc := NewScanTaskProcessor(taskQueLen * 2 + 1)
	wp := workerpool.New(3, 2, taskQueLen + 1, scanTaskProc, workerpool.SRC)
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
		// Get alive host by PingScreening
		aliveHosts := PingScreen(args.Hosts, args.TTL)
		if len(aliveHosts) == 0 {
			return
		}
		total := len(aliveHosts) * len(args.Ports)
		fmt.Println("total of targets are", total)
		// initiate the workpool
		processor, wp := initWorkPool(total)
		
		
		count := 0
		for _, h := range aliveHosts {
			for _, p := range args.Ports {
				scanTask := NewScanTask(h, p, "", args.TTL)
				wp.AppendTask(scanTask, 0)
				wp.AppendTask(scanTask, 1)
				//wp.AppendTask(scanTask, 2)
			}
		}
		for r := range processor.GetResults() {
			res := tools.NewRes(r.Host, r.Protocol, r.Type, r.IdString, r.BString, r.Port, r.TTL, r.Buf, r.IdBool)
			fmt.Printf("%s:%d has been scanned\n", res.Host, res.Port)
			fmt.Println(string(res.Buf))
			// TODO: Judge The Result
			
			
			
			
			count++

			if count >= total * 2 {
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


