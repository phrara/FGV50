package scanner

import (
	"csl/scanner/connection"
	"csl/tools/addr"
	"csl/tools/workerpool"

)


type ScanTask struct {
	Host string
	Port int
	Protocol string
	TTL int
}

func NewScanTask(host string, port int, proto string, ttl int) *ScanTask {
	return &ScanTask{
		Host: host,
		Port: port,
		Protocol: proto,	
		TTL: ttl,
	}
}

func (s *ScanTask) AddrString() string {
	return addr.Conbine(s.Host, s.Port)
}


type Result struct {
	ScanTask
	Buf []byte
	Type string
}

func NewResult(st ScanTask, buf []byte, typ string) Result {
	return Result{
		Buf: buf,
		ScanTask: st,
		Type: typ,
	}
}

const (
	DefaultResQueLen = 100
)


type ScanTaskProcessor struct {
	resC chan Result
	*workerpool.Task
}

func NewScanTaskProcessor(resQueLen int) *ScanTaskProcessor {
	return &ScanTaskProcessor{
		resC: make(chan Result, resQueLen),
	}
}

func (s *ScanTaskProcessor) Process(task any) error {
	t := task.(ScanTask)
	
	b, err := connection.TlsProtocol(t.Host, t.Port, t.TTL)
	s.PutResult(NewResult(t, b, "tls"))



	return err
}

func (s *ScanTaskProcessor) Handle(err error) {
	// TODO
}


func (s *ScanTaskProcessor) GetResults() chan Result {
	if s.resC != nil {
		return s.resC
	} else {
		s.resC = make(chan Result)
		close(s.resC)
		return s.resC
	}
}

func (s *ScanTaskProcessor) PutResult(res Result)  {
	if s.resC != nil {
		s.resC <- res
	} else {
		s.resC = make(chan Result, DefaultResQueLen)
		s.resC <- res
	}
}

func (s *ScanTaskProcessor) CloseResC()  {
	close(s.resC)
}
