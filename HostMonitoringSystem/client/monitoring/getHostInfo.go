package monitoring

import (
	"github.com/shirou/gopsutil/v3/host"
)

//{"hostname":"DESKTOP-3EHS48G","uptime":704788,"bootTime":1613448053,"procs":225,"os":"windows","platform":"Microsoft Windows 10 Pro","platformFamily":"Standalone Workstation","platformVersion":"10.0.19042 Build 19042","kernelVersion":"10.0.19042 Build 19042","kernelArch":"x86_64","virtualizationSystem":"","virtualizationRole":"","hostId":"67097ed2-552d-4675-99d1-32b30bb75547"}
var (
	HostInfo *host.InfoStat
)

func GetHostInfo() (*host.InfoStat,error){
	HostInfo, err = host.Info()
	return HostInfo, err
}


