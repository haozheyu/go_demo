package client

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/haozheyu/go_demo/HostMonitoringSystem/client/monitoring"
	"github.com/robfig/cron/v3"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"time"
)

// 注册节点到etcd： /cron/workers/IP地址
type Register struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
	localIP string // 本机IP
	RegMassage Message
}
type Message struct {
	Ip string `json:"ip"`
	Timestamp int64 `json:"timestamp"`
	CPUInfo *monitoring.CpuInfo `json:"cpu_info"`
	CPULoad *monitoring.CPULoad `json:"cpu_load"`
	DiskIOStat []disk.IOCountersStat `json:"disk_io_stat"`
	HostInfo *host.InfoStat `json:"host_info"`
	VMStat *mem.VirtualMemoryStat `json:"vm_stat"`
	NetStat *monitoring.NetStat `json:"net_stat"`
	PartitionStat []monitoring.PartitionsInfo `json:"partition_stat"`
	ProcessCount int32 `json:"process_count"`
}

var (
	G_register *Register
)

// 注册到/client/IP, 并自动续租
func (register *Register) keepOnline() {
	var (
		regKey string
		err error
		cpuinfo *monitoring.CpuInfo
		cpuload *monitoring.CPULoad
		diskstat []disk.IOCountersStat
		hostinfo *host.InfoStat
		vmstat *mem.VirtualMemoryStat
		netstat *monitoring.NetStat
		partstat []monitoring.PartitionsInfo
		rep []byte
	)
	regKey = "/client/" + register.localIP
	register.RegMassage.Timestamp = time.Now().Unix()
	if cpuinfo, err = monitoring.GetCpuInfo();err ==nil{register.RegMassage.CPUInfo = cpuinfo}
	if cpuload,err = monitoring.GetCpuLoad();err == nil{register.RegMassage.CPULoad = cpuload}
	if diskstat,err = monitoring.GetDiskIOInfo(); err == nil {register.RegMassage.DiskIOStat= diskstat}
	if hostinfo,err = monitoring.GetHostInfo();err == nil {register.RegMassage.HostInfo = hostinfo}
	if vmstat,err = monitoring.GetMemLoad();err == nil {register.RegMassage.VMStat = vmstat}
	if netstat,err = monitoring.GetNetStat();err == nil {register.RegMassage.NetStat = netstat}
	if partstat,err = monitoring.GetPartitions();err == nil {register.RegMassage.PartitionStat = partstat}
	register.RegMassage.ProcessCount = monitoring.RunProcessCount()
	register.RegMassage.Ip = register.localIP
	if rep, err = json.Marshal(register.RegMassage);err !=nil {
		log.Println("etcd put rep marshal fail")
	}
	// 注册到etcd
	if _, err = register.kv.Put(context.TODO(), regKey, string(rep)); err != nil {
		log.Println("node register fail")
	}
}

func InitRegister() {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
		localIp string
		err error
		s *cron.Cron
	)

	// 初始化配置
	config = clientv3.Config{
		Endpoints: G_config.EtcdEndpoints, // 集群地址
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond, // 连接超时
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		return
	}

	// 本机IP
	if localIp, err = monitoring.GetLocalIP(); err != nil {
		return
	}

	// 得到KV和Lease的API子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	G_register = &Register{
		client: client,
		kv: kv,
		lease: lease,
		localIP: localIp,
	}
	// 服务注册
	s = cron.New()
    if _,err = s.AddFunc(G_config.Crontab,G_register.keepOnline);err!=nil{
    	log.Printf("计划任务解析失败 %s:",err.Error())
	}
    s.Start()
	select {}
}