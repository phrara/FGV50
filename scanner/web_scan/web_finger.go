package webscan

import (
	"encoding/json"
	"os"
)

var FPS *FPList

type Fingerprint struct {
	Cms      string
	Method   string
	Location string
	Keyword  []string
}

type FPList struct {
	FP []Fingerprint
}

func NewFP(cms, meth, loc string, kw []string) *Fingerprint {
	return &Fingerprint{
		Cms:      cms,
		Method:   meth,
		Location: loc,
		Keyword:  kw,
	}
}

func LoadWebFP(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var fps FPList
	err = json.Unmarshal(data, &fps)
	if err != nil {
		return err
	}

	FPS = &fps
	return nil

}