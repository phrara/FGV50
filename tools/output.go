package tools

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var resFilePath string

func init() {
	resFilePath, _ = os.Getwd()
	resFilePath = filepath.Join(resFilePath, "/json/res.json")
}

func Open() *os.File {
	return openFile(resFilePath)
}

func Write(reslist []*Result, output *os.File, total int) {
	cnts := make([]map[string]interface{}, 0, total)
	for _, v := range reslist {
		if v.Protocol == "" {
			continue
		}
		content := make(map[string]interface{})
		content["time"]=time.Now().Format("2006-01-02 15:04:05")
		content["type"]=v.Type
		content["protocol"]=v.Protocol
		content["host"]=v.Host
		content["port"]=v.Port
		content["idstring"]=v.IdString
		content["banner"]=v.BString
		cnts = append(cnts, content)
		
	}
	mjson,_ :=json.Marshal(cnts)
	writeContent(output, string(mjson))

}

func Close(file *os.File) {
	err := file.Close()
	if err != nil {
		fmt.Println("Close file" + file.Name() + "exception")
	} else {
		fmt.Println("The identification results are saved in " + file.Name())
	}
}

func openFile(name string) *os.File {
	osFile, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("Failed to open file " + name)
	}
	fmt.Println("Open file success")
	return osFile
}

func writeContent(file *os.File, content string) {
	_, err := file.Write([]byte(content + "\n"))
	if err != nil {
		fmt.Println("Write failed: " + content)
	} else {
		fmt.Println("Write Success: " + content)
	}
}

func ByteToStringParse1(p []byte) string {
	var w []string
	var res string
	for i := 0; i < len(p); i++ {
		if p[i] > 32 && p[i] < 127 {
			w = append(w, string(p[i]))
			continue
		}
		asciiTo16 := fmt.Sprintf("\\x%s", hex.EncodeToString(p[i:i+1]))
		w = append(w,asciiTo16)
	}
	res = strings.Join(w, "")
	if strings.Contains(res, "\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00") {
		s := strings.Split(res, "\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00")
		return s[0]
	}
	return res
}
