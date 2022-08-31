package jg

import (
	"fgv50/tools"
)

func IdentifyTcp(ret *tools.Result) bool {
	

	ok := Tcpjudge(ret)
	if ok {
		return true
	}
	ok = TcpConnjudge(ret)
	if ok {
		return true
	}
	 ok = TcpSMB(ret)
	  
	 if ok {
	 	return true
	 }

	return false
}

func IdentifyTls(result *tools.Result) bool {
	protocol := result.Protocol
	runAll := true
	if protocol != "" {
		runAll = false
	}
	// if protocol == "http" || protocol == "https" || runAll {
	// 	if TlsHTTPS(result) {
	// 		return true
	// 	}
	// 	return true
	// }
	if protocol == "rdp" || runAll {
		if TlsRDP(result) {
			return true
		}
	}
	if protocol == "redis-ssl" || runAll {
		if TlsRedisSsl(result) {
			return true
		}
	}
	return false
}

func IdentifyUdp(result *tools.Result) bool {
	protocol := result.Protocol
	runAll := true
	if protocol != "" {
		runAll = false
	}
	if protocol == "nbns" || runAll {
		if UdpNbns(result) {
			return true
		}
	}
	return false
}

