package workerpool

import "sync"

type Strategy int

const (
	// RR select que by using Round Robin
	RR Strategy = 0
	// SRC select que by using outer source
	SRC Strategy = 1
)

type WorkerPool struct {
	MaxQueNum int
	AvgWorker int
	MaxQueLen int
	taskQueue []chan any
	tasker    Tasker
	qid       int
	Mod       Strategy
	mux sync.Mutex
}

func New(maxQueNum, avgWorker, maxQueLen int, tasker Tasker, mod Strategy) *WorkerPool {
	return &WorkerPool{
		MaxQueNum: maxQueNum,
		AvgWorker: avgWorker,
		MaxQueLen: maxQueLen,
		taskQueue: make([]chan any, maxQueNum),
		tasker:    tasker,
		qid:       0,
		Mod:       mod,
		mux: sync.Mutex{},
	}
}

func (w *WorkerPool) work(qid int) {
	for task := range w.taskQueue[qid] {
		err := w.tasker.Process(task, qid)
		if len(err) != 0 {
			w.tasker.Handle(err)
		}
	}
}

func (w *WorkerPool) Start() {
	for i := 0; i < w.MaxQueNum; i++ {
		w.taskQueue[i] = make(chan any, w.MaxQueLen)
		for j := 0; j < w.AvgWorker; j++ {
			go w.work(i)
		}
	}
}

// AppendTask used Round Robin or Source
func (w *WorkerPool) AppendTask(task any, src int) {
	switch w.Mod {
	case RR:
		w.mux.Lock()
		w.qid = w.qid % w.MaxQueNum
		w.taskQueue[w.qid] <- task
		//fmt.Printf("task que %d recv a task\n", w.qid)
		w.qid++
		w.mux.Unlock()
	case SRC:
		qid := src % w.MaxQueNum
		w.taskQueue[qid] <- task
		//fmt.Printf("task que %d recv a task\n", w.qid)
	}

}

func (w *WorkerPool) Shut() {
	for i := 0; i < w.MaxQueNum; i++ {
		close(w.taskQueue[i])
	}
}
