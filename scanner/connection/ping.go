package connection

import (
	"github.com/go-ping/ping"
)
func Ping(ip string) bool {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		return false
	}
	pinger.Count = 1
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return false
	} else {
		return true
	}
}