package webscan

import (
	"fgv50/scanner/connection"
	"fgv50/tools"
	"fgv50/tools/workerpool"
	"os"
	"path/filepath"
)

func init() {
	finger_path, _ = os.Getwd()
	finger_path = filepath.Join(finger_path, "/json/finger.json")
	LoadWebFP(finger_path)
}
const (
	DefaultResQueLen = 100
)
var (
	finger_path string
)

type WebScanTask struct {
	url string
	urlNum int
	cycle int

}


type WebScanProcessor struct {
	resC chan *tools.Result
	*workerpool.Task
}

func NewWebScanProcessor(resQueLen int) *WebScanProcessor {
	return &WebScanProcessor{
		resC: make(chan *tools.Result, resQueLen),
	}
}

func (w *WebScanProcessor) Process(task any, qid int) []error {
	errs := make([]error, 0, 20)
	t := task.(WebScanTask)
	res, err := fingerScan(t)
	errs = append(errs, err)
	w.PutResult(res)


	//TODO
	return errs
}

func (w *WebScanProcessor) Handle(errs []error) {

}

func fingerScan(t WebScanTask) (*tools.Result, error) {
	var data *connection.Resps
	data, err := connection.Httprequest(t.url, t.cycle, t.urlNum)
	if err != nil {
		return nil, err
	}
	if data.Statuscode == 400 {
		data, err = connection.Httprequest(t.url, t.cycle, t.urlNum)
		if err != nil {
			return nil, err
		}
	}
	t.cycle++
	for _, j := range data.Jsurl {
		if j != "" {
			// TODO
			panic("TODO")
		}
	}
	return &tools.Result{}, nil
}


func (w *WebScanProcessor) PutResult(res *tools.Result)  {

	if w.resC != nil {
		w.resC <- res
	} else {
		w.resC = make(chan *tools.Result, DefaultResQueLen)
		w.resC <- res
	}
}

func (w *WebScanProcessor) GetResults() chan *tools.Result {
	if w.resC != nil {
		return w.resC
	} else {
		w.resC = make(chan *tools.Result)
		close(w.resC)
		return w.resC
	}
}