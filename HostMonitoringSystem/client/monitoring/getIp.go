package monitoring

import "net"

func GetLocalIP() (ipv4 string, err error) {
	var (
		addres []net.Addr
		addr net.Addr
		ipNet *net.IPNet // IP地址
		isIpNet bool
	)
	// 获取所有网卡
	if addres, err = net.InterfaceAddrs(); err != nil {
		return
	}
	// 取第一个非lo的网卡IP
	for _, addr = range addres {
		// 这个网络地址是IP地址: ipv4, ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String()	// 192.168.1.1
				return
			}
		}
	}
	return
}
