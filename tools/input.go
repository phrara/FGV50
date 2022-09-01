package tools

import "os"

func ReadVulJson() []byte {
	b, _ := os.ReadFile(vulJsonPath)
	return b
}

func ReadResJson() []byte {
	b, _ := os.ReadFile(resFilePath)
	return b
}
