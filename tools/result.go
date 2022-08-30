package tools

type Result struct {
	Host string
	Port int
	Protocol string
	TTL int
	Buf []byte
	Type string
	IdBool   bool
	IdString  string
	BString  string
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