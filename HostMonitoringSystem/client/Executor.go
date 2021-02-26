package client

import (
	"context"
	"github.com/haozheyu/go_demo/HostMonitoringSystem/client/monitoring"
	"os/exec"
)

// 任务执行状态
type JobExecuteInfo struct {
	IP string   `json:"ip"`// 任务节点
	Exec string `json:"exec"`
}
// 任务执行结果
type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo	// 执行的信息
	Output []byte // 脚本输出
	Err error // 脚本错误原因
}

// 执行一个任务
func ExecuteJob(info *JobExecuteInfo) (err error){
	var (
		cmd *exec.Cmd
		output []byte
		result *JobExecuteResult
		ip string
	)
	// 任务结果
	result = &JobExecuteResult{
		ExecuteInfo: info,
		Output: make([]byte, 0),
	}
	if ip, err = monitoring.GetLocalIP();err == nil {
		if info.IP == ip {
			// 执行shell命令
			cmd = exec.CommandContext(context.TODO(), G_config.BashDir, "-c", info.Exec)
			// 执行并捕获输出
			output, err = cmd.CombinedOutput()
			// 将输出传递
			result.Output = output
		}
	}
	return err
}