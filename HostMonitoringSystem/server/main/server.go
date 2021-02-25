package main

import (
	"flag"
	"fmt"
	"github.com/haozheyu/go_demo/HostMonitoringSystem/server"
	"runtime"
)

var (
	confFile string // 配置文件路径
)

// 解析命令行参数
func initArgs() {
	// master -config ./master.json -xxx 123 -yyy ddd
	// master -h
	flag.StringVar(&confFile, "config", "./server.json", "指定server.json")
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
	if err = server.InitConfig(confFile); err != nil {
		goto ERR
	}

	// 初始化服务发现模块
	if err = server.InitWorkerMgr(); err != nil {
		goto ERR
	}

	// 让其中指定节点执行命令
	if err = server.InitJobMgr(); err != nil {
		goto ERR
	}
	//
	//// 日志管理器
	//if err =InitLogMgr(); err != nil {
	//	goto ERR
	//}
	//
	////  任务管理器
	//if err = InitJobMgr(); err != nil {
	//	goto ERR
	//}
	//
	// 启动Api HTTP服务
	if err = server.InitApiServer(); err != nil {
		goto ERR
	}

	select {}

	return

ERR:
	fmt.Println(err)
}
