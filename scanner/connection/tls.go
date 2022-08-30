package connection

import (
	"bytes"
	"crypto/tls"
	"net"
	"strconv"
	"time"

	"github.com/zhzyker/dismap/pkg/logger"
)

func TlsProtocol(host string, port int, timeout int) ([]byte, error) {
	conn, err := ConnProxyTls(host, port, timeout)
	if logger.DebugError(err) {
		return nil, err
	}
	msg := "GET /test HTTP/1.1\r\n\r\n"
	_, err = conn.Write([]byte(msg))
	if logger.DebugError(err) {
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
	if logger.DebugError(err) {
		return nil, err
	}
	msg = "GET /test HTTP/1.1\r\n\r\n"
	_, err = conn.Write([]byte(msg))
	if logger.DebugError(err) {
		return nil, err
	}
	_ = conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	reply = make([]byte, 256)
	_, _ = conn.Read(reply)
	if conn != nil {
		_ = conn.Close()
	}
	return reply, err
}

func ConnProxyTls(host string, port int, timeout int) (net.Conn, error) {
	target := net.JoinHostPort(host, strconv.Itoa(port))
	// scheme, address, proxyUri, err := parse.ProxyParse()
	// TLS does not support proxy function temporarily
	// 2022-02-23 by zhzyker
	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: time.Duration(timeout) * time.Second},
		"tcp",
		target,
		&tls.Config{InsecureSkipVerify: true})
	if logger.DebugError(err) {
		return nil, err
	}
	err = conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	if logger.DebugError(err) {
		if conn != nil {
			_ = conn.Close()
		}
		return nil, err
	}
	return conn, nil
}