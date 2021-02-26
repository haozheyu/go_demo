package server

import (
	"encoding/json"
	"github.com/haozheyu/go_demo/crontab/common"
	"net"
	"net/http"
	"strconv"
	"time"
)

// 任务的HTTP接口
type ApiServer struct {
	httpServer *http.Server
}

var (
	// 单例对象
	G_apiServer *ApiServer
)

type ResponseMsg struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Count int    `json:"count"`
	Data  interface{} `json:"data"`
}

// POST 保存执行任务接口
func handleExecJob(resp http.ResponseWriter, req *http.Request) {
	var (
		err error
		postJob string
		oldJob string
		bytes []byte
		ip string
	)

	// 1, 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	// 2, 取表单中的job字段
	postJob = req.PostForm.Get("exec")
	ip = req.PostForm.Get("ip")
	if len(ip)== 0 {goto ERR}
	// 4, 保存到etcd
	if oldJob, err = G_jobMgr.SaveJob(postJob,ip); err != nil {
		goto ERR
	}
	// 5, {"code":0,"msg":"",count:1,"data":[{"ip":"127.0.0.1", },{}]
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	// 6, 返回异常应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// 获取注册的client节点列表
func handleWorkerList(resp http.ResponseWriter, req *http.Request) {
	var (
		msg []Message
		err error
		bytes []byte
		data ResponseMsg
	)

	if msg, err = G_workerMgr.ListWorkers(); err != nil {
		goto ERR
	}
    data.Code = 0
    data.Msg = ""
    data.Count = len(msg)
    data.Data = msg
	// 正常应答
	if bytes, err = json.Marshal(data);err == nil{
		resp.Write(bytes)
	}
	return

ERR:
	data.Code = -1
	if bytes, err := json.Marshal(data);err == nil{
		data.Msg = err.Error()
		resp.Write(bytes)
	}
}

// 单个节点信息返回
func handleWorkerInfo(resp http.ResponseWriter,req *http.Request){
	var (
		msg *Message
		err error
		bytes []byte
		data ResponseMsg
		ip string
	)
	ip = req.PostFormValue("ip")
	if msg, err = G_workerMgr.ClientInfo(ip); err != nil {
		goto ERR
	}
	data.Code = 0
	data.Msg = ""
	data.Data = msg
	// 正常应答
	if bytes, err = json.Marshal(data);err == nil{
		resp.Write(bytes)
	}
	return

ERR:
	data.Code = -1
	if bytes, err := json.Marshal(data);err == nil{
		data.Msg = err.Error()
		resp.Write(bytes)
	}
}

// 初始化服务
func InitApiServer() (err error){
	var (
		mux *http.ServeMux
		listener net.Listener
		httpServer *http.Server
		staticDir http.Dir	// 静态文件根目录
		staticHandler http.Handler	// 静态文件的HTTP回调
	)

	// 配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/exec", handleExecJob) //单个节点执行命令
	//mux.HandleFunc("/job/delete", handleJobDelete)
	//mux.HandleFunc("/job/list", handleJobList)
	//mux.HandleFunc("/job/kill", handleJobKill)
	//mux.HandleFunc("/job/log", handleJobLog)
	mux.HandleFunc("/worker/node", handleWorkerInfo) //获取单个节点信息
	mux.HandleFunc("/worker/list", handleWorkerList) //获取注册节点列表

	//  /index.html

	// 静态文件目录
	staticDir = http.Dir(G_config.WebRoot)
	staticHandler = http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler))	//   ./webroot/index.html

	// 启动TCP监听
	if listener, err = net.Listen("tcp", ":" + strconv.Itoa(G_config.Export)); err != nil {
		return
	}

	// 创建一个HTTP服务
	httpServer = &http.Server{
		ReadTimeout: time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
		Handler: mux,
	}

	// 赋值单例
	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	// 启动了服务端
	go httpServer.Serve(listener)

	return
}