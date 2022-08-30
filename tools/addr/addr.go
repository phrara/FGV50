package addr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func Conbine(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}

func Split(addr string) (string, int) {
	s := strings.Split(addr, ":")
	p, _ := strconv.ParseInt(s[1], 10, 64)
	return s[0], int(p)
}

func AddrCheck(addr string) bool {
	ip, port := Split(addr)
	compile, _ := regexp.Compile(`((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)`)
	if b := compile.MatchString(ip); b {	
		if port < 0 || port > 65535 {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}