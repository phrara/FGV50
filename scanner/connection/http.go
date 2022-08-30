package connection

import (
	"crypto/tls"
	"net/http"
	"time"
)

func ConnProxyHttp(request *http.Request, timeout int) (*http.Response, error) {
	
	
	// var tr *http.Transport
	/*
	scheme, address, pUri, err := parse.ProxyParse()
	if logger.DebugError(err) {
		return nil, err
	}
	if scheme == "http" {
		proxyUri, err := url.Parse(pUri)
		if logger.DebugError(err) {
			logger.Error("Cannot initialize http proxy")
			return nil, err
		}
		tr = &http.Transport{
			Proxy:           http.ProxyURL(proxyUri),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else if scheme == "socks5" {
		dialSocksProxy, err := proxy.SOCKS5("tcp", address, nil, &net.Dialer{
			Timeout:   time.Duration(timeout) * time.Second,
			KeepAlive: 10 * time.Second,
		})
		if logger.DebugError(err) {
			logger.Error("Cannot initialize socks5 proxy")
			return nil, err
		}
		tr = &http.Transport{
			Dial:            dialSocksProxy.Dial,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	*/
	
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	response, err := client.Do(request)
	return response, err
}
