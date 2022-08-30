package connection

import (
	"net"
	"strconv"
	"time"
)

var udpPort = []int{53, 111, 123, 137, 138, 139, 12345}

func UdpProtocol(host string, port int, timeout int) ([]byte, error) {
	if isContainInt(udpPort, port) {
		return make([]byte, 256), nil
	}
	return make([]byte, 256), nil
}

func isContainInt(items []int, item int) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func ConnProxyUdp(host string, port int, timeout int) (net.Conn, error) {
	target := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("udp", target, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	/*
	scheme, address, proxyUri, err := parse.ProxyParse()
	if logger.DebugError(err) {
		return nil, err
	}
	var conn net.Conn

	if proxyUri == "" {
		conn, err = net.DialTimeout("udp", target, time.Duration(timeout)*time.Second)
		if logger.DebugError(err) {
			return nil, err
		}
	}

	if scheme == "http" {
		target := net.JoinHostPort(host, strconv.Itoa(port))
		conn, err = net.DialTimeout("udp", target, time.Duration(timeout)*time.Second)
		if logger.DebugError(err) {
			return nil, err
		}
	}

	if scheme == "socks5" {
		dialer, err := socks5.NewClient(address, "", "", timeout, timeout)
		if logger.DebugError(err) {
			logger.Error("Cannot initialize socks5 proxy")
			return nil, err
		}
		conn, err = dialer.Dial("udp", target)
		if logger.DebugError(err) {
			return nil, err
		}
	}
*/
	err = conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	if err != nil {
		if conn != nil {
			_ = conn.Close()
		}
		return nil, err
	}
	return conn, nil
}