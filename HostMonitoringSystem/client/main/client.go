package main

import (
	"flag"
	"fmt"
	"github.com/haozheyu/go_demo/HostMonitoringSystem/client"
	"runtime"
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
	// go client.InitRegister()
    // 暴露web接口
	go client.ExportWeb()
	// 命令执行器
	// go client.RunExec()




	// 正常退出
	select { }

	return

ERR:
	fmt.Println(err)
}

