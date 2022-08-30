package connection


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