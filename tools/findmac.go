package tools

import (
	"fmt"
	"os/exec"
	"strings"
)

func FindMac(ip string)(string,string) {
	fmt.Println("find the mac of "+ip)
	cmd := exec.Command("python3","./python/find_mac_dev.py",ip)
	// 执行命令，并返回结果
	output,err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("find mac and dev success")
		res:=strings.Split(string(output), "\n")
		mac:=strings.ReplaceAll(res[0],"\r","")
		dev:=strings.ReplaceAll(res[1],"\r","")
		fmt.Println(mac)
		fmt.Println(dev)
		return mac,dev
	}
	return "",""
}
