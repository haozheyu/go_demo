package main

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/zieckey/etcdsync"
	"log"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err error
	)

	// 客户端配置
	config = clientv3.Config{
		Endpoints: []string{"192.168.32.201:2379"},
		DialTimeout: 5 * time.Second,
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	client = client

	m, err := etcdsync.New("/cron/jobs/job1", 100, []string{"http://192.168.43.201:2379"})
	if m == nil || err != nil {
		log.Printf("etcdsync.New failed")
		return
	}
	err = m.Lock()
	if err != nil {
		log.Printf("etcdsync.Lock failed")
	} else {
		log.Printf("etcdsync.Lock OK")
	}

	log.Printf("Get the lock. Do something here.")

	err = m.Unlock()
	if err != nil {
		log.Printf("etcdsync.Unlock failed")
	} else {
		log.Printf("etcdsync.Unlock OK")
	}
}

