package jg

import (
	"bytes"
	"encoding/hex"
	"fgv50/tools"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func Tcpjudge(ret *tools.Result) bool {
	bt := ret.Buf
	if cap(bt) == 0 || len(bt) == 0 {
		return false
	}

	ok, _ := regexp.Match(`(mysql_native_password|MySQL server|MariaDB server|mysqladmin flush-hosts)`, ret.Buf) //mysql
	if ok {
		ret.Protocol = "mysql"
		return true

	}
	ok, _ = regexp.Match(`(^-ERR(.*)command|^-DENIED.Redis)`, ret.Buf) //redis
	if ok {
		ret.Protocol = "redis"
		return true

	}
	ok, _ = regexp.Match(`^RFB \d`, ret.Buf) //vnc
	if ok {
		ret.Protocol = "vnc"
		return true
	}
	ok, _ = regexp.Match(`(Telnet>|^BeanShell)`, ret.Buf) //telnet
	if ok {
		ret.Protocol = "telnet"
		return true
	} else if strings.Contains(hex.EncodeToString(ret.Buf[0:2]), "fffb") {
		ret.Protocol = "telnet"
		return true
	}
	buff := []byte{0x41, 0x01, 0x02} //snmp
	if bytes.Equal(buff[:], ret.Buf[0:3]) {
		ret.Protocol = "snmp"
		return true
	}
	ok, _ = regexp.Match(`(^220[ -](.*)ESMTP|^421(.*)Service not available|^554 )`, ret.Buf) // smtp
	if ok {
		ret.Protocol = "smtp"
		return true
	}
	ok, _ = regexp.Match(`^\+OK`, ret.Buf) // pop3
	if ok {
		ret.Protocol = "pop3"
		return true
	}
	ok, _ = regexp.Match(`^* OK`, ret.Buf) //  imap
	if ok {
		ret.Protocol = "imap"
		return true
	}
	ok, _ = regexp.Match(`(^220(.*FTP|.*FileZilla)|^421(.*)connections)`, ret.Buf) // ftp
	if ok {
		ret.Protocol = "ftp"
		return true
	}
	ok, _ = regexp.Match(`ActiveMQ`, ret.Buf) //activemq
	if ok {
		ver, err := strconv.ParseUint(hex.EncodeToString(buff[13:17]), 16, 32)
		if err == nil {

			ret.IdString = fmt.Sprintf("[%s]", fmt.Sprintf("Version:%s", strconv.FormatUint(ver, 10)))
		}
		ret.Protocol = "activemq"
		return true
	}
	ok, _ = regexp.Match(`^SSH.\d`, ret.Buf) //ssh
	if ok {
		str := ret.BString
		ret.BString = strings.Split(str, "\\x0d\\x0a")[0]
		ret.Protocol = "ssh"
		return true

	}

	return false
}
