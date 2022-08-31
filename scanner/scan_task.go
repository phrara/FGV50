package scanner

import (
	"fgv50/scanner/connection"
	"fgv50/tools/addr"
	"fgv50/tools/workerpool"
	"sync"

)

func init() {

	
}
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

	IdBool   bool
	IdString  string
	BString  string

}

func NewResult(st *ScanTask, buf []byte, typ string) Result {
	return Result{
		Buf: buf,
		ScanTask: st,
		Type: typ,
	}
}

const (
	DefaultResQueLen = 300
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

var (
	tcpM sync.Mutex
	tlsM sync.Mutex
	udpM sync.Mutex
)
func (s *ScanTaskProcessor) Process(task any, qid int) []error {
	t := task.(*ScanTask)
	errs := make([]error, 0, 3)
	
	
	switch qid {
	case 0:
		tcpM.Lock()
		//fmt.Printf("start scanning tcp/%s:%d\n", t.Host, t.Port)
		b2, err2 := connection.TcpProtocol(t.Host, t.Port, t.TTL)
		errs = append(errs, err2)
		s.PutResult(NewResult(t, b2, "tcp"))
		tcpM.Unlock()
	case 1:
		tlsM.Lock()
		//fmt.Printf("start scanning tls/%s:%d\n", t.Host, t.Port)

		b1, err1 := connection.TlsProtocol(t.Host, t.Port, t.TTL)
		errs = append(errs, err1)
		s.PutResult(NewResult(t, b1, "tls"))
		tlsM.Unlock()
	case 2:
		udpM.Lock()
		
		//fmt.Printf("start scanning udp/%s:%d\n", t.Host, t.Port)

		b3, err3 := connection.UdpProtocol(t.Host, t.Port, t.TTL)
		errs = append(errs, err3)
		s.PutResult(NewResult(t, b3, "udp"))
		udpM.Unlock()
	default:
	}
	
	return errs
	
}

func (s *ScanTaskProcessor) Handle(errs []error) {
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
