package jg

import (
	"bytes"
	"encoding/hex"
	"fgv50/scanner/connection"
	"fgv50/tools"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

func TcpConnjudge(ret *tools.Result) bool {
	
	ok := TcpDceRpc(ret)
	if ok {
		return true
	}
	ok = TcpFrp(ret)
	if ok {
		return true
	}
	ok = TcpGIOP(ret)
	if ok {
		return true
	}
	ok = TcpLDAP(ret)
	if ok {
		return true
	}
	ok = TcpMssql(ret)
	if ok {
		return true
	}
	ok = TcpOracle(ret)
	if ok {
		return true
	}
	ok = TcpRDP(ret)
	if ok {
		return true
	}
	ok = TcpRMI(ret)
	if ok {
		return true
	}
	ok = TcpRTSP(ret)
	if ok {
		return true
	}
	ok = TcpSocks(ret)
	if ok {
		return true
	} else {
		
		return false
	}


}

func TcpDceRpc(ret *tools.Result) bool {

	timeout := ret.TTL
	host := ret.Host
	port := ret.Port
	conn, err := connection.ConnProxyTcp(host, port, timeout)
	if err != nil {
		return false
	}
	msg1 := "\x05\x00\x0b\x03\x10\x00\x00\x00\x48\x00\x00\x00\x01\x00\x00\x00\xf8\x0f\xf8\x0f\x00\x00\x00\x00\x01\x00\x00\x00\x00\x00\x01\x00\xc4\xfe\xfc\x99\x60\x52\x1b\x10\xbb\xcb\x00\xaa\x00\x21\x34\x7a\x00\x00\x00\x00\x04\x5d\x88\x8a\xeb\x1c\xc9\x11\x9f\xe8\x08\x00\x2b\x10\x48\x60\x02\x00\x00\x00"
	msg2 := "\x05\x00\x00\x03\x10\x00\x00\x00\x18\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x05\x00"
	_, err = conn.Write([]byte(msg1))
	if err != nil {
		return false
	}
	reply1 := make([]byte, 256)
	_, _ = conn.Read(reply1)
	if hex.EncodeToString(reply1[0:8]) != "05000c0310000000" {
		return false
	}
	_, err = conn.Write([]byte(msg2))
	if err != nil {
		return false
	}
	reply2 := make([]byte, 512)
	_, _ = conn.Read(reply2)
	if conn != nil {
		_ = conn.Close()
	}
	ret.Protocol = "dcerpc"

	c := 0
	zero := make([]byte, 1)
	var buffer bytes.Buffer
	for i := 0; i < len(reply2[42:]); {
		b := reply2[42:][i : i+2]
		i += 2
		if 42+i == len(reply2[42:]) {
			break
		}
		if string(b) == "\x09\x00" {
			break
		}
		if string(b) == "\x07\x00" {
			c += 1
			if c == 6 {
				break
			}
			buffer.Write([]byte("\x7C\x7C"))
			ret.BString = strings.Join([]string{string(buffer.Bytes())}, ",")
			continue
		}
		if bytes.Equal(b[0:1], zero[0:1]) {
			continue
		}
		buffer.Write(b[0:1])
		ret.BString = strings.Join([]string{string(buffer.Bytes())}, ",")
		if c == 6 {
			break
		}
	}
	ret.Buf = reply2
	return true

}

// TODO
func TcpFrp(ret *tools.Result) bool {
	timeout := ret.TTL
	host := ret.Host
	port := ret.Port
	conn, err := connection.ConnProxyTcp(host, port, timeout)
	if err != nil {
		return false
	}
	msg := "\x00\x01\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00"
	_, err = conn.Write([]byte(msg))
	if err != nil {
		return false
	}
	reply := make([]byte, 256)
	_, _ = conn.Read(reply)
	if conn != nil {
		_ = conn.Close()
	}
	var buffer [256]byte
	if bytes.Equal(reply[:], buffer[:]) {
		return false
	} else if hex.EncodeToString(reply[0:12]) != "000100020000000100000000" {
		return false
	}
	ret.Protocol = "frp"
	ret.BString = frpByteToStringParse(reply[0:12])
	ret.Buf = reply
	return true

}
func frpByteToStringParse(p []byte) string {
	var w []string
	var res string
	for i := 0; i < len(p); i++ {
		asciiTo16 := fmt.Sprintf("\\x%s", hex.EncodeToString(p[i:i+1]))
		w = append(w, asciiTo16)
	}
	res = strings.Join(w, "")
	return res
}

// TODO
func TcpGIOP(ret *tools.Result) bool {
	timeout := ret.TTL
	host := ret.Host
	port := ret.Port
	conn, err := connection.ConnProxyTcp(host, port, timeout)
	if err != nil {
		return false
	}
	msg := "\x47\x49\x4f\x50\x01\x02\x00\x03\x00\x00\x00\x17\x00\x00\x00\x02\x00\x00\x00\x00\x00\x00\x00\x0b\x4e\x61\x6d\x65\x53\x65\x72\x76\x69\x63\x65"
	_, err = conn.Write([]byte(msg))
	if err != nil {
		return false
	}
	reply := make([]byte, 256)
	_, _ = conn.Read(reply)
	if conn != nil {
		_ = conn.Close()
	}
	if strings.Contains(hex.EncodeToString(reply[0:4]), "47494f50") == false {
		return false
	}
	ret.Protocol = "giop"
	ret.BString = ByteToStringParse2(reply[0:4])
	ret.Buf = reply
	return true

}

// TODO
func TcpLDAP(ret *tools.Result) bool {
	timeout := ret.TTL
	host := ret.Host
	port := ret.Port
	conn, err := connection.ConnProxyTcp(host, port, timeout)
	if err != nil {
		return false
	}
	msg := "\x30\x0c\x02\x01\x01\x60\x07\x02\x01\x03\x04\x00\x80\x00"
	_, err = conn.Write([]byte(msg))
	if err != nil {
		return false
	}
	reply := make([]byte, 256)
	_, _ = conn.Read(reply)
	if conn != nil {
		_ = conn.Close()
	}
	var buffer [256]byte
	if bytes.Equal(reply[:], buffer[:]) {
		return false
	}
	if strings.Contains(hex.EncodeToString(reply), "010004000400") == false {
		return false
	}
	ret.Protocol = "ldap"
	ret.BString = ByteToStringParse2(reply[0:16])
	ret.Buf = reply
	return true

}

// TODO
func TcpMssql(ret *tools.Result) bool {
	timeout := ret.TTL
	host := ret.Host
	port := ret.Port
	conn, err := connection.ConnProxyTcp(host, port, timeout)
	if err != nil {
		return false
	}
	msg := "\x12\x01\x00\x34\x00\x00\x00\x00\x00\x00\x15\x00\x06\x01\x00\x1b\x00\x01\x02\x00\x1c\x00\x0c\x03\x00\x28\x00\x04\xff\x08\x00\x01\x55\x00\x00\x02\x4d\x53\x53\x51\x4c\x53\x65\x72\x76\x65\x72\x00\x00\x00\x31\x32"
	_, err = conn.Write([]byte(msg))
	if err != nil {
		return false
	}
	reply := make([]byte, 256)
	_, _ = conn.Read(reply)
	if conn != nil {
		_ = conn.Close()
	}
	var buffer [256]byte
	if bytes.Equal(reply[:], buffer[:]) {
		return false
	} else if hex.EncodeToString(reply[0:4]) != "04010025" {
		return false
	} else {
		ret.Protocol = "mssql"
	}
	v, bo := getVersion(reply)
	if bo {
		ret.IdBool = true
		ret.IdString = fmt.Sprintf("[%s]", fmt.Sprintf("Version:%s", v))

	}
	ret.BString = ByteToStringParse1(reply)
	ret.Buf = reply
	return true

}

// TODO
func TcpOracle(ret *tools.Result) bool {
	timeout := ret.TTL
	host := ret.Host
	port := ret.Port
	conn, err := connection.ConnProxyTcp(host, port, timeout)
	if err != nil {
		return false
	}
	msg := "\x00\x5a\x00\x00\x01\x00\x00\x00\x01\x36\x01\x2c\x00\x00\x08\x00\x7f\xff\x7f\x08\x00\x00\x00\x01\x00\x20\x00\x3a\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x34\xe6\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x28\x43\x4f\x4e\x4e\x45\x43\x54\x5f\x44\x41\x54\x41\x3d\x28\x43\x4f\x4d\x4d\x41\x4e\x44\x3d\x56\x45\x52\x53\x49\x4f\x4e\x29\x29"
	_, err = conn.Write([]byte(msg))
	if err != nil {
		return false
	}
	reply := make([]byte, 256)
	_, _ = conn.Read(reply)
	if conn != nil {
		_ = conn.Close()
	}
	ok, _ := regexp.Match(`\(DESCRIPTION=`, ret.Buf)
	if ok {
		ret.Protocol = "oracle"
	} else {
		var buffer [256]byte
		if bytes.Equal(reply[:], buffer[:]) {
			return false
		} else if hex.EncodeToString(reply[0:8]) != "0065000004000000" {
			return false
		} else {
			ret.Protocol = "oracle"
		}
	}
	var vsnnum string
	banStr := ByteToStringParse2(reply)
	grep := regexp.MustCompile(`\(VSNNUM=(\d*)\)`)
	vsnnum = grep.FindStringSubmatch(banStr)[1]
	v, err := strconv.ParseInt(vsnnum, 10, 64)
	if err != nil {
		ret.IdBool = false
	}
	hexVsnnum := strconv.FormatInt(v, 16)
	maj, _ := strconv.ParseUint(hexVsnnum[0:1], 16, 32)
	min, _ := strconv.ParseUint(hexVsnnum[1:2], 16, 32)
	a, _ := strconv.ParseUint(hexVsnnum[2:4], 16, 32)
	b, _ := strconv.ParseUint(hexVsnnum[4:5], 16, 32)
	c, _ := strconv.ParseUint(hexVsnnum[5:7], 16, 32)
	var version string
	if err == nil {
		version = fmt.Sprintf("%s.%s.%s.%s.%s",
			strconv.FormatUint(maj, 10),
			strconv.FormatUint(min, 10),
			strconv.FormatUint(a, 10),
			strconv.FormatUint(b, 10),
			strconv.FormatUint(c, 10),
		)
	} else {
		ret.IdBool = false

	}
	ret.IdBool = true
	ret.IdString = fmt.Sprintf("[%s]", fmt.Sprintf("Version:%s", version))
	ret.BString = banStr
	ret.Buf = reply
	return true

}

// TODO
func TcpRDP(ret *tools.Result) bool {
	timeout := ret.TTL
	host := ret.Host
	port := ret.Port
	conn, err := connection.ConnProxyTcp(host, port, timeout)
	if err != nil {
		return false
	}
	msg2 := "\x03\x00\x00\x2b\x26\xe0\x00\x00\x00\x00\x00\x43\x6f\x6f\x6b\x69\x65\x3a\x20\x6d\x73\x74\x73\x68\x61\x73\x68\x3d\x75\x73\x65\x72\x30\x0d\x0a\x01\x00\x08\x00\x00\x00\x00\x00"
	_, err = conn.Write([]byte(msg2))
	if err != nil {
		return false
	}
	reply := make([]byte, 256)
	_, _ = conn.Read(reply)
	if conn != nil {
		_ = conn.Close()
	}
	var buffer [256]byte
	if bytes.Equal(reply[:], buffer[:]) {
		return false
	} else if hex.EncodeToString(reply[0:8]) != "030000130ed00000" {
		return false
	} else {
		ret.Protocol = "rdp"
	}
	os := map[string]string{}
	/*** msg1 os finger ***
	os["030000130ed000001234000209080002000000"]="Windows 7/Windows Server 2008"
	os["030000130ed00000123400021f080002000000"]="Windows 10/Windows Server 2019"
	os["030000130ed00000123400020f080002000000"]="Windows 8.1/Windows Server 2012 R2"
	*/
	os["030000130ed000001234000209080000000000"] = "Windows 7/Windows Server 2008 R2"
	os["030000130ed000001234000200080000000000"] = "Windows 7/Windows Server 2008"
	os["030000130ed000001234000201080000000000"] = "Windows Server 2008 R2"
	os["030000130ed000001234000207080000000000"] = "Windows 8/Windows server 2012"
	os["030000130ed00000123400020f080000000000"] = "Windows 8.1/Windows Server 2012 R2"
	os["030000130ed000001234000300080001000000"] = "Windows 10/Windows Server 2016"
	os["030000130ed000001234000300080005000000"] = "Windows 10/Windows 11/Windows Server 2019"
	for k, v := range os {
		if k == hex.EncodeToString(reply[0:19]) {
			ret.BString = v
			return true
		}
	}
	ret.BString = hex.EncodeToString(reply[0:19])
	ret.Buf = reply
	return true

}

// TODO
func TcpRMI(ret *tools.Result) bool {
	timeout := ret.TTL
	host := ret.Host
	port := ret.Port
	conn, err := connection.ConnProxyTcp(host, port, timeout)
	if err != nil {
		return false
	}
	msg := "\x4a\x52\x4d\x49\x00\x02\x4b"
	_, err = conn.Write([]byte(msg))
	if err != nil {
		return false
	}
	reply := make([]byte, 256)
	_, _ = conn.Read(reply)
	if conn != nil {
		_ = conn.Close()
	}
	var buffer [256]byte
	if bytes.Equal(reply[:], buffer[:]) {
		return false
	} else if hex.EncodeToString(reply[0:1]) != "4e" {
		return false
	}
	ret.Protocol = "rmi"
	ret.BString = ByteToStringParse1(reply)
	ret.Buf = reply
	return true

}

// TODO
func TcpRTSP(ret *tools.Result) bool {
	var buff []byte
	buff = ret.Buf
	ok, err := regexp.Match(`^RTSP/`, buff)
	if err != nil {
		return false
	}

	if ok {
		ret.Protocol = "rtsp"
		return true
	}

	if rtsp(ret) {
		return true
	}
	return false

}
func rtsp(ret *tools.Result) bool {
	timeout := ret.TTL
	host := ret.Host
	port := ret.Port

	address := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := connection.ConnProxyTcp(host, port, timeout)
	if err != nil {
		return false
	}

	msg := fmt.Sprintf("OPTIONS rtsp://%s RTSP/1.0\r\nCSeq:1\r\n\r\n", address)
	_, err = conn.Write([]byte(msg))
	if err != nil {
		return false
	}

	reply := make([]byte, 256)
	_, _ = conn.Read(reply)
	if conn != nil {
		_ = conn.Close()
	}

	var buffer [256]byte
	if bytes.Equal(reply[:], buffer[:]) {
		return false
	} else if hex.EncodeToString(reply[0:4]) != "52545350" {
		return false
	}
	ret.Protocol = "rtsp"
	ret.BString = ByteToStringParse1(reply)
	ret.Buf = reply
	return true
}

// TODO
func TcpSocks(ret *tools.Result) bool {

	if socks4(ret) {
		ret.Protocol = "socks4"
		return true
	}

	if socks5(ret) {
		ret.Protocol = "socks5"
		return true
	}
	return false

}
func socks4(ret *tools.Result) bool {
	timeout := ret.TTL
	host := ret.Host
	port := ret.Port
	conn, err := connection.ConnProxyTcp(host, port, timeout)
	if err != nil {
		return false
	}

	p1 := strconv.FormatInt(int64(port/256), 16)
	if p1str, _ := strconv.Atoi(p1); p1str < 10 {
		p1 = fmt.Sprintf("0%s", p1)
	}
	p1byte, err := hex.DecodeString(p1)

	p2 := strconv.FormatInt(int64(port%256), 16)
	if p2str, _ := strconv.Atoi(p2); p2str < 10 {
		p2 = fmt.Sprintf("0%s", p2)
	}
	p2byte, err := hex.DecodeString(p2)
	if err != nil {
		return false
	}
	msgByte := []byte{0x04, 0x01}
	msgByte = append(msgByte, p1byte[0])
	msgByte = append(msgByte, p2byte[0])
	msgStr := hex.EncodeToString(msgByte)

	grep := regexp.MustCompile(`(\d*).(\d*).(\d*).(\d*)`)
	ip := grep.FindStringSubmatch(host)[1:5]
	for _, i := range ip {
		if i == "0" {
			msgStr += "00"
			continue
		}
		i64, _ := strconv.ParseInt(i, 10, 64)

		n := strconv.FormatInt(i64, 16)
		if len(n) != 2 {
			n = fmt.Sprintf("0%s", n)
		}
		msgStr += n
	}
	msgStr += "0100"
	hexData, _ := hex.DecodeString(msgStr)
	_, err = conn.Write(hexData)
	if err != nil {
		return false
	}

	reply := make([]byte, 256)
	_, _ = conn.Read(reply)
	if conn != nil {
		_ = conn.Close()
	}

	if string(reply[1]) == "\x5b" {
		ret.BString = ByteToStringParse2(reply[0:8])
		ret.Buf = reply
		return true
	}
	return false
}
func socks5(ret *tools.Result) bool {
	timeout := ret.TTL
	host := ret.Host
	port := ret.Port
	conn, err := connection.ConnProxyTcp(host, port, timeout)
	if err != nil {
		return false
	}

	msgSocks5 := "\x05\x02\x00\x02"
	/*
		\x05 - Version: 5
		\x02 - Authentication Method Count: 2
		\x00 - Method[0]: 0 (No authentication)
		\x02 - Method[1]: 2 (Username/Password)
	*/
	_, err = conn.Write([]byte(msgSocks5))
	if err != nil {
		return false
	}

	reply := make([]byte, 256)
	_, _ = conn.Read(reply)
	if conn != nil {
		_ = conn.Close()
	}

	var buffer bytes.Buffer
	if string(reply[0]) == "\x05" {
		buffer.WriteString(fmt.Sprintf("[%s]", "Version:Socks5"))
	} else {
		return false
	}
	if string(reply[1]) == "\x00" {
		buffer.WriteString(fmt.Sprintf("[%s]", "Method:No Authentication(\\x00)"))
	}
	if string(reply[1]) == "\x02" {
		buffer.WriteString(fmt.Sprintf("[%s]", "Method:Username/Password(\\x02)"))
	}
	ret.IdBool = true
	ret.IdString = buffer.String()
	ret.BString = ByteToStringParse2(reply[0:2])
	ret.Buf = reply
	return true
}

//TODO

func getVersion(reply []byte) (string, bool) {
	m, err := strconv.ParseUint(hex.EncodeToString(reply[29:30]), 16, 32)
	if err != nil {
		return "", false
	}
	s, err := strconv.ParseUint(hex.EncodeToString(reply[30:31]), 16, 32)
	if err != nil {
		return "", false
	}
	r, err := strconv.ParseUint(hex.EncodeToString(reply[31:33]), 16, 32)
	if err != nil {
		return "", false
	}
	v := fmt.Sprintf("%d.%d.%d", m, s, r)
	return v, true
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
		w = append(w, asciiTo16)
	}
	res = strings.Join(w, "")
	if strings.Contains(res, "\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00") {
		s := strings.Split(res, "\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00")
		return s[0]
	}
	return res
}
func ByteToStringParse2(p []byte) string {
	var w []string
	var res string
	for i := 0; i < len(p); i++ {
		if p[i] > 32 && p[i] < 127 {
			w = append(w, string(p[i]))
			continue
		}
		asciiTo16 := fmt.Sprintf("\\x%s", hex.EncodeToString(p[i:i+1]))
		w = append(w, asciiTo16)
	}
	res = strings.Join(w, "")
	return res
}
