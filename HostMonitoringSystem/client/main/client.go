package main

import (
	"flag"
	"fmt"
	"github.com/haozheyu/go_demo/HostMonitoringSystem/client"
	"runtime"
	"time"
)

var (
	confFile string // 配置文件路径
)

// 解析命令行参数
func initArgs() {
	// worker -config ./worker.json
	// worker -h
	flag.StringVar(&confFile, "config", "./client.json", "client.json")
	flag.Parse()
}

// 初始化线程数量
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)

	// 初始化命令行参数
	initArgs()

	// 初始化线程
	initEnv()

	// 加载配置
	if err = client.InitConfig(confFile); err != nil {
		goto ERR
	}

	// 服务注册
	go client.InitRegister()
    // 暴露web接口
	go client.ExportWeb()

	//// 定时命令执行器
	//if err = worker.InitJobMgr(); err != nil {
	//	goto ERR
	//}

	// 正常退出
	for {
		time.Sleep(1 * time.Second)
	}

	return

ERR:
	fmt.Println(err)
}

