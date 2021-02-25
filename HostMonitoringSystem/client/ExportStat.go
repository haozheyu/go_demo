package client

import (
	"encoding/json"
	"fmt"
	"github.com/haozheyu/go_demo/HostMonitoringSystem/client/monitoring"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"net/http"
	"strconv"
	"time"
)

func ExportWeb(){ //应用监控处理
    var (
    	resp Register
    	ip string
    	err error
		cpuinfo *monitoring.CpuInfo
		cpuload *monitoring.CPULoad
		diskstat []disk.IOCountersStat
		hostinfo *host.InfoStat
		vmstat *mem.VirtualMemoryStat
		netstat *monitoring.NetStat
		partstat []monitoring.PartitionsInfo
		rep []byte
    	export string
    )
	http.HandleFunc("/monitor", func(writer http.ResponseWriter, request *http.Request) {
		ip, err = monitoring.GetLocalIP()
		resp.localIP = ip
		resp.RegMassage.Timestamp = time.Now()
		if cpuinfo, err = monitoring.GetCpuInfo();err ==nil{resp.RegMassage.CPUInfo = cpuinfo}
		if cpuload,err = monitoring.GetCpuLoad();err == nil{resp.RegMassage.CPULoad = cpuload}
		if diskstat,err = monitoring.GetDiskIOInfo(); err == nil {resp.RegMassage.DiskIOStat= diskstat}
		if hostinfo,err = monitoring.GetHostInfo();err == nil {resp.RegMassage.HostInfo = hostinfo}
		if vmstat,err = monitoring.GetMemLoad();err == nil {resp.RegMassage.VMStat = vmstat}
		if netstat,err = monitoring.GetNetStat();err == nil {resp.RegMassage.NetStat = netstat}
		if partstat,err = monitoring.GetPartitions();err == nil {resp.RegMassage.PartitionStat = partstat}
		resp.RegMassage.ProcessCount = monitoring.RunProcessCount()
		if rep, err = json.Marshal(resp.RegMassage);err != nil{
			log.Printf("exprot json 解析失败",err.Error())
		}
		writer.Write(rep)
	})
	export = strconv.Itoa(G_config.Export)
	http.ListenAndServe(fmt.Sprintf(":%s",export),nil)
}
