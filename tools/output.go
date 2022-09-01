package tools

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ( 
	resFilePath string
   	vulJsonPath string
)
func init() {
	resFilePath, _ = os.Getwd()
	vulJsonPath, _ = os.Getwd()
	resFilePath = filepath.Join(resFilePath, "/json/res.json")
	vulJsonPath = filepath.Join(vulJsonPath, "/json/ali_cve.json")
}

func Open() *os.File {
	return openFile(resFilePath)
}

func Write(resJson []byte, output *os.File) {
	writeContent(output, string(resJson))
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
