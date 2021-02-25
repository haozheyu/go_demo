package monitoring

import (
	"encoding/json"
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	ip,err := GetLocalIP()
	t.Log(ip,err)
}

func TestGetHostInfo(t *testing.T) {
	info, err2 := GetHostInfo()
	t.Log(info,err2)
}

func TestGetPartitions(t *testing.T) {
	getPartitions, _ := GetPartitions()
	marshal, _ := json.Marshal(getPartitions)
	t.Log(string(marshal))
}
func TestGetCpuLoad(t *testing.T) {
	load, err2 := GetCpuLoad()
	t.Log(load,err2)
}
func TestGetMemLoad(t *testing.T) {
	load, err2 := GetMemLoad()
	t.Log(load,err2)
}
func TestGetDiskIOInfo(t *testing.T) {
	info, err2 := GetDiskIOInfo()
	t.Log(info,err2)
}
func TestGetProcessStat(t *testing.T) {
	count := RunProcessCount()
	t.Log(count)
}
func TestGetNetStat(t *testing.T) {
	netStat,err	 := GetNetStat()
	t.Log(netStat.TimeWaitNumber,netStat.ESTABLISHEDNumber,netStat.NetConnectionCount,netStat.BytesRecv,netStat.PacketsRecv,netStat.Dropin, err)
}
