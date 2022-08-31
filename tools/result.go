package tools

type Result struct {
	Host string `json:"host"`
	Port int `json:"port"`
	Protocol string `json:"protocol"`
	TTL int `json:"ttl"`
	Buf []byte `json:"buf"`
	Type string `json:"type"`
	IdBool   bool `json:"id_bool"`
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