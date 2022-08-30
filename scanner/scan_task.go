package scanner

import (
	"csl/scanner/connection"
	"csl/tools/addr"
	"csl/tools/workerpool"

	"github.com/zhzyker/dismap/pkg/logger"
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
	*ScanTask
	Buf []byte
	Type string
}

func NewResult(st *ScanTask, buf []byte, typ string) Result {
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

func (s *ScanTaskProcessor) Process(task any) []error {
	t := task.(*ScanTask)
	errs := make([]error, 0, 3)

	b1, err1 := connection.TlsProtocol(t.Host, t.Port, t.TTL)
	errs = append(errs, err1)
	s.PutResult(NewResult(t, b1, "tls"))

	b2, err2 := connection.TcpProtocol(t.Host, t.Port, t.TTL)
	errs = append(errs, err2)
	s.PutResult(NewResult(t, b2, "tcp"))

	b3, err3 := connection.UdpProtocol(t.Host, t.Port, t.TTL)
	errs = append(errs, err3)
	s.PutResult(NewResult(t, b3, "udp"))

	return errs
}

func (s *ScanTaskProcessor) Handle(errs []error) {
	for _, e := range errs {
		if e != nil {
			logger.DebugError(e)
		}
	}
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
