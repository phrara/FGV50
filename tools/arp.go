package tools

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func ARP() {
	// 通过exec.Command函数执行命令或者shell
	// 第一个参数是命令路径，当然如果PATH路径可以搜索到命令，可以不用输入完整的路径
	// 第二到第N个参数是命令的参数
	// 下面语句等价于执行命令: ls -l ar/
	// cmd := exec.Command("powershell", )
	// err := cmd.Run()

	// if err != nil {
	// 	// 命令执行失败
	// 	panic(err)
	// }

	cmd := exec.Command("powershell", "arp", " -a")
	// 执行命令，并返回结果
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	// 因为结果是字节数组，需要转换成string
	//fmt.Println(string(output))
	//
	filePath := "./tools/arp.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("open file error=%v\n", err)
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	//for i := 0; i < 5; i++ {
	writer.WriteString(string(output))
	//}
	writer.Flush()

}
