package jg

import (
	"bytes"
	"fgv50/tools"
	"fmt"
	"strconv"
	"strings"
	"time"
	"fgv50/scanner/connection"
	"github.com/zhzyker/dismap/pkg/logger"
)

func UdpNbns(result *tools.Result) bool {
	return nbnsIdentifyResult(result)
}

func nbnsIdentifyResult(result *tools.Result) bool {
	host := result.Host
	port := result.Port
	timeout := result.TTL
	conn, err := connection.ConnProxyUdp(host, port, timeout)
	if err!=nil {
		return false
	}
	msg := []byte{
		0x0, 0x00, 0x0, 0x10, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x20, 0x43, 0x4b, 0x41, 0x41,
		0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41,
		0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x0, 0x0,
		0x21, 0x0, 0x1,
	}
	_, err = conn.Write(msg)
	if err!=nil {
		if conn != nil {
			_ = conn.Close()
		}
		return false
	}
	reply := make([]byte, 256)
	err = conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	if err!=nil {
		if conn != nil {
			_ = conn.Close()
		}
		return false
	}
	_, _ = conn.Read(reply)
	if conn != nil {
		_ = conn.Close()
	}

	var buffer [256]byte
	if bytes.Equal(reply[:], buffer[:]) {
		return false
	}
	var n int
	NumberFoNames, _ := strconv.Atoi(convert([]byte{reply[56:57][0]}[:]))
	var flagGroup string
	var flagUnique string
	var flagDC string

	for i := 0; i < NumberFoNames; i++ {
		data := reply[n+57+18*i : n+57+18*i+18]
		if string(data[16:17]) == "\x84" || string(data[16:17]) == "\xC4" {
			if string(data[15:16]) == "\x1C" {
				flagDC = "Domain Controllers"
			}
			if string(data[15:16]) == "\x00" {
				flagGroup = nbnsByteToStringParse(data[0:16])
			}
			if string(data[14:16]) == "\x02\x01" {
				flagGroup = nbnsByteToStringParse(data[0:16])
			}
		} else if string(data[16:17]) == "\x04" || string(data[16:17]) == "\x44" || string(data[16:17]) == "\x64" {
			if string(data[15:16]) == "\x1C" {
				flagDC = "Domain Controllers"
			}
			if string(data[15:16]) == "\x00" {
				flagUnique = nbnsByteToStringParse(data[0:16])
			}
			if string(data[15:16]) == "\x20" {
				flagUnique = nbnsByteToStringParse(data[0:16])
			}

		}
	}
	if flagGroup == "" && flagUnique == "" {
		return false
	}

	result.BString = flagGroup + "\\" + flagUnique
	result.IdString = fmt.Sprintf("[%s]", logger.LightRed(flagDC))
	if len(flagDC) != 0 {
		result.IdBool = true
	} else {
		result.IdBool = false
	}
	result.Protocol = "nbns"
	result.Buf = reply
	return true
}

func convert(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, "")
}

func nbnsByteToStringParse(p []byte) string {
	var w []string
	var res string
	for i := 0; i < len(p); i++ {
		if p[i] > 32 && p[i] < 127 {
			w = append(w, string(p[i]))
			continue
		}
	}
	res = strings.Join(w, "")
	return res
}
