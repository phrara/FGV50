package scanner

import "csl/tools/workerpool"

const (
	DefaultResQueLen = 100
)

type ScanTask struct {
	resC chan []byte
	*workerpool.Task
}

func NewScanTask(resQueLen int) *ScanTask {
	return &ScanTask{
		resC: make(chan []byte, resQueLen),
	}
}

func (s *ScanTask) Process(task any) error {
	// TODO
	return nil
}

func (s *ScanTask) Handle(err error) {
	// TODO
}


func (s *ScanTask) GetResults() chan []byte {
	if s.resC != nil {
		return s.resC
	} else {
		s.resC = make(chan []byte)
		close(s.resC)
		return s.resC
	}
}

func (s *ScanTask) PutResult(res []byte)  {
	if s.resC != nil {
		s.resC <- res
	} else {
		s.resC = make(chan []byte, DefaultResQueLen)
		s.resC <- res
	}
}

func (s *ScanTask) CloseResC()  {
	close(s.resC)
}