package jg

import (
	"fgv50/tools"
	"fmt"
	"regexp"
)

func TlsHTTPS(ret *tools.Result, inurl string) bool {

	var buff []byte= ret.Buf
	ok, err := regexp.Match(`^HTTP/\d.\d \d*`, buff)
	if err != nil {
		return false
	}
	if ok {
		ret.Protocol = "https"
		httpResult, httpErr := httpIdentifyResult(ret, inurl)
		if httpErr != nil {
			ret.BString = "None"
			return true
		}
		ret.BString = httpResult["http.title"].(string)
		//_, _ := url.Parse(httpResult["http.target"].(string))

		r := httpResult["http.result"].(string)
		c := fmt.Sprintf("[%s]", httpResult["http.code"].(string))
		if len(r) != 0 {
			ret.IdBool = true
			ret.IdString = fmt.Sprintf("%s %s", c, r)
			//result["note"] = httpResult["http.target"].(string)
			return true
		} else {
			ret.IdBool = true
			ret.IdString = c
			//result["note"] = httpResult["http.target"].(string)
			return true
		}
	}
	return false

}
