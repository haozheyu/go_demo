package client

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/haozheyu/go_demo/HostMonitoringSystem/client/monitoring"
	"log"
	"time"
)

// 注册节点到etcd： /job/exec
type JobExec struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
	localIP string // 本机IP
	RegMassage Message
}

var (
	G_JobExec *JobExec
)

// 获取etcd /job/exec
func (je *JobExec) GetJobExec() (jei *JobExecuteInfo,err error){
	var (
		ctx context.Context
		cancel context.CancelFunc
		getResp  *clientv3.GetResponse
	)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()
	if getResp,err = je.kv.Get(ctx,"/job/exec");err!=nil {
		return nil, err
	}
	for _,ev := range getResp.Kvs{
		if err = json.Unmarshal(ev.Value, &jei);err!=nil{
			return nil, err
		}
	}
	return jei, err
}

func (je *JobExec) DelJobExec() (jei *JobExecuteInfo,err error){
	var (
		ctx context.Context
		cancel context.CancelFunc
		delResp  *clientv3.DeleteResponse
	)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()
	if delResp,err = je.kv.Delete(ctx,"/job/exec",clientv3.WithPrevKV());err!=nil {
		return nil, err
	}
	if len(delResp.PrevKvs) != 0 {
		for _,ev := range delResp.PrevKvs{
			if err = json.Unmarshal(ev.Value, &jei);err!=nil{
				return nil, err
			}
		}
	}
	return jei, err
}

func RunExec(){
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
		localIp string
		jobexec *JobExecuteInfo
		err error
	)

	// 初始化配置
	config = clientv3.Config{
		Endpoints: G_config.EtcdEndpoints, // 集群地址
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond, // 连接超时
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		log.Println("clientv3.New(config) run fail")
	}

	// 本机IP
	if localIp, err = monitoring.GetLocalIP(); err != nil {
		log.Println("monitoring.GetLocalIP run fail")
	}

	// 得到KV和Lease的API子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	G_JobExec = &JobExec{
		client: client,
		kv: kv,
		lease: lease,
		localIP: localIp,
	}
	for {
		time.Sleep(10 * time.Second)
		jobexec, err = G_JobExec.GetJobExec()
		if err != nil {
			continue
		}
		if jobexec == nil || jobexec.IP != G_JobExec.localIP {
			continue
		}
		//执行命令
		if err = ExecuteJob(jobexec);err!=nil {
			continue
		}
		//删除job/exec
		time.Sleep(2 * time.Second)
		jobexec, err = G_JobExec.DelJobExec()
		if err != nil {
			continue
		}
		log.Println("删除的任务是：",jobexec)
	}
}