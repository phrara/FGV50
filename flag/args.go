package flag

import (
	"csl/err"
	"csl/tools/addr"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

var DefaultPorts = []int{21, 22, 23, 25, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 110, 135, 137, 138, 139, 143, 389, 443, 445, 587, 631, 800, 801, 808,
	880, 888, 1000, 1024, 1025, 1080, 1099, 1389, 1433, 1521, 2383, 3306, 3307, 3388, 3389, 3443, 5000, 5357, 5560, 5800, 5900, 6379, 7000, 7001, 7007,
	7010, 7788, 8000, 8001, 8002, 8003, 8004, 8005, 8006, 8007, 8008, 8009, 8010, 8011, 8030, 8060, 8070, 8080, 8081, 8082, 8083, 8084, 8085, 8086, 8087,
	8088, 8089, 8090, 8091, 8092, 8093, 8094, 8095, 8096, 8097, 8098, 8099, 8161, 8175, 8188, 8189, 8200, 8443, 8445, 8448, 8554, 8800, 8848, 8880, 8881,
	8888, 8899, 8983, 8989, 9000, 9001, 9002, 9008, 9010, 9043, 9060, 9080, 9081, 9082, 9083, 9084, 9085, 9086, 9087, 9088, 9089, 9090, 9091, 9092, 9093,
	9094, 9095, 9096, 9097, 9099, 9443, 9600, 9628, 9800, 9999, 11001, 13443, 49155, 50050, 61616}

type Args struct {
	Url   *url.URL `json:"url"`
	Hosts []string `json:"hosts"`
	Ports []int    `json:"ports"`
	TTL int `json:"ttl"`
}

func NewArgs(u, ns, ip string, port int, ttl int) (*Args, error) {
	a := new(Args)
	var ports []int
	if u != "" {
		u2, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		a.Url = u2
		a.TTL = ttl
		return a, nil
	} else {
		hosts, err1 := ParseNetworkSegment(ns)
		if err1 != nil {
			return nil, err1
		}

		if ip != "" {
			if b := addr.AddrCheck(addr.Conbine(ip, 0)); b {
				hosts = append(hosts, ip)
			} else {
				return nil, err.ErrIllFormedIP
			}
		}
		if port != -1 {
			if port >= 0 || port <= 65535 {
				ports = make([]int, 0, 1)
				ports = append(ports, port)
			} else {
				return nil, err.ErrPortOutRange
			}
		} else {
			ports = DefaultPorts
		}
		a.Hosts = hosts
		a.Ports = ports
		a.TTL = ttl
		return a, nil
	}
}

func (a *Args) PrintAllAddrs() {
	for _, h := range a.Hosts {
		for _, p := range a.Ports {
			fmt.Printf("%s:%d\n", h, p)
		}
	}
}

func ParseNetworkSegment(ns string) ([]string, error) {
	h := make([]string, 0, 15)
	if ns == "" {
		return h, nil
	}
	if strings.Contains(ns, "~") && strings.Contains(ns, "+") {
		return nil, err.ErrIllFormedNS
	}
	if strings.Contains(ns, "~") {
		s := strings.Split(ns, "~")
		if len(s) != 2 {
			return nil, err.ErrIllFormedNS
		} else {
			if b := addr.AddrCheck(addr.Conbine(s[0], 0)); !b {
				return nil, err.ErrIllFormedIP
			}
			ipSlice := strings.Split(s[0], ".")
			sti, er := strconv.ParseInt(ipSlice[3], 10, 64)
			if er != nil {
				return nil, err.ErrIllFormedNS
			}
			edi, er := strconv.ParseInt(s[1], 10, 64)
			if er != nil {
				return nil, err.ErrIllFormedNS
			}
			if edi > 65535 {
				return nil, err.ErrPortOutRange
			}
			for i := sti; i <= edi; i++ {
				seg4 := strconv.FormatInt(i, 10)
				ip := ip(append(ipSlice[:3], seg4)...)
				h = append(h, ip)
			}
			return h, nil
		}
	} else if strings.Contains(ns, "+") {
		s := strings.Split(ns, "+")
		if len(s) != 2 {
			return nil, err.ErrIllFormedNS
		} else {
			ipSlice := strings.Split(s[0], ".")
			sti, er := strconv.ParseInt(ipSlice[3], 10, 64)
			if er != nil {
				return nil, err.ErrIllFormedNS
			}
			inc, er := strconv.ParseInt(s[1], 10, 64)
			if er != nil {
				return nil, err.ErrIllFormedNS
			}
			edi := sti + inc
			if edi > 65535 {
				return nil, err.ErrPortOutRange
			}
			for i := sti; i <= edi; i++ {
				seg4 := strconv.FormatInt(i, 10)
				ip := ip(append(ipSlice[:3], seg4)...)
				h = append(h, ip)
			}
			return h, nil
		}
	} else {
		return nil, err.ErrIllFormedNS
	}
}

func ip(seg ...string) string {
	return fmt.Sprintf("%s.%s.%s.%s", seg[0], seg[1], seg[2], seg[3])
}
