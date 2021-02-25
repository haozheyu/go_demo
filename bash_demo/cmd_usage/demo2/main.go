package main

import (
	"os/exec"
	"fmt"
)

func main() {
	var (
		cmd *exec.Cmd
		output []byte
		err error
	)

	// 生成Cmd
	cmd = exec.Command("C:\\Program Files\\Git\\bin\\bash.exe", "-c", "echo hello_world")

	// 执行了命令, 捕获了子进程的输出( pipe )
	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println(err)
		return
	}

	// 打印子进程的输出
	fmt.Println(string(output))
}