package server

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/haozheyu/go_demo/HostMonitoringSystem/client/monitoring"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"time"
)

// /cron/workers/
type WorkerMgr struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease

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
	G_workerMgr *WorkerMgr
)

// 获取在线worker列表
func (workerMgr *WorkerMgr) ListWorkers() (rests []Message, err error) {
	var (
		getResp *clientv3.GetResponse
		kv *mvccpb.KeyValue
		value []byte
		msg Message
	)

	// 获取目录下所有Kv
	if getResp, err = workerMgr.kv.Get(context.TODO(), "/client/", clientv3.WithPrefix()); err != nil {
		return
	}

	// 解析每个节点的IP
	for _, kv = range getResp.Kvs {
		// kv.Key : /client/192.168.2.1
		//ip = strings.TrimPrefix(string(kv.Key), "/client/")
		value = kv.Value
		if err = json.Unmarshal(value, &msg);err !=nil {
			log.Printf("clienMgr rest fail",err)
		}
		rests = append(rests, msg)
	}
	return
}

// 获取单个client信息
func (client *WorkerMgr) ClientInfo(ip string)(rest *Message,err error){
	var (
		getResp *clientv3.GetResponse
		kv *mvccpb.KeyValue
		value []byte
	)

	// 获取目录下所有Kv
	if getResp, err = client.kv.Get(context.TODO(), "/client/"+ip); err != nil {
		return
	}

	// 解析每个节点的IP
	for _, kv = range getResp.Kvs {
		// kv.Key : /client/192.168.2.1
		//ip = strings.TrimPrefix(string(kv.Key), "/client/")
		value = kv.Value
		if err = json.Unmarshal(value, &rest);err !=nil {
			log.Printf("clienMgr rest fail",err)
		}
	}
	return
}

func InitWorkerMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
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

	// 得到KV和Lease的API子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	G_workerMgr = &WorkerMgr{
		client :client,
		kv: kv,
		lease: lease,
	}
	return
}