package scanner

import "csl/tools/workerpool"

func initWorkPool() (*ScanTask, *workerpool.WorkerPool) {
	scanTask := NewScanTask(DefaultResQueLen)
	wp := workerpool.New(4, 2, 20, scanTask, workerpool.RR)
	wp.Start()
	return scanTask, wp
}

func Run() {
	// scanTask, wp := initWorkPool()
	// TODO
}