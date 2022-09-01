package tools

type Result struct {
	Time string `json:"time"`
	Host string `json:"host"`
	Port int `json:"port"`
	Protocol string `json:"protocol"`
	TTL int `json:"-"`
	Buf []byte `json:"-"`
	Type string `json:"type"`
	IdBool   bool `json:"-"`
	IdString  string `json:"idstring"`
	BString  string `json:"banner"`
}

func NewRes(host, proto, typ, idstr, bstr string, port, ttl int, buf []byte, idb bool) *Result {
	return &Result{
		Host: host,
		Port: port,
		Protocol: proto,
		TTL: ttl,
		Buf: buf,
		Type: typ,
		IdBool: idb,   
		IdString: idstr,
		BString: bstr,  
	}
}

type HardWare struct {
	Host string `json:"host"`
	Mac string `json:"mac"`
	Dev string `json:"dev"`
}