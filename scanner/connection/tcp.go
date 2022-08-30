package connection

import (
	"bytes"
	"net"
	"strconv"
	"time"
)

func TcpProtocol(host string, port int, timeout int) ([]byte, error) {
	conn, err := ConnProxyTcp(host, port, timeout)
	if err != nil {
		return nil, err
	}
	_ = conn.SetDeadline(time.Now().Add(time.Duration(2) * time.Second))
	reply := make([]byte, 256)
	_, err = conn.Read(reply)
	var buffer [256]byte
	if err == nil && !bytes.Equal(reply[:], buffer[:]) {
		if conn != nil {
			_ = conn.Close()
		}
		return reply, nil
	}
	conn, err = ConnProxyTcp(host, port, timeout)
	if err != nil {
		return nil, err
	}
	msg := "GET /test HTTP/1.1\r\n\r\n"
	_, err = conn.Write([]byte(msg))
	if err != nil {
		return nil, err
	}
	_ = conn.SetDeadline(time.Now().Add(time.Duration(2) * time.Second))
	reply = make([]byte, 256)
	_, _ = conn.Read(reply)
	if conn != nil {
		_ = conn.Close()
	}
	return reply, nil
}

func ConnProxyTcp(host string, port int, timeout int) (net.Conn, error) {
	target := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", target, time.Duration(timeout)*time.Second)
	if err != nil {
			return nil, err
	}
	/*
	scheme, address, proxyUri, err := parse.ProxyParse()
	if err != nil {
		return nil, err
	}

	var conn net.Conn
	if proxyUri == "" {
		conn, err = net.DialTimeout("tcp", target, time.Duration(timeout)*time.Second)
		if err != nil {
			return nil, err
		}
	}
	if scheme == "http" {
		conn, err = net.DialTimeout("tcp", target, time.Duration(timeout)*time.Second)
		if err != nil {
			return nil, err
		}
	}
	if scheme == "socks5" {
		tcp, err := proxy.SOCKS5("tcp", address, nil, &net.Dialer{
			Timeout:   time.Duration(timeout) * time.Second,
			KeepAlive: 10 * time.Second,
		})
		if err != nil {
			logger.Error("Cannot initialize socks5 proxy")
			return nil, err
		}
		conn, err = tcp.Dial("tcp", target)
		if err != nil {
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