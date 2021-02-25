package client

import (
	"io/ioutil"
	"encoding/json"
)

// 程序配置
type Config struct {
	EtcdEndpoints []string `json:"etcdEndpoints"`
	EtcdDialTimeout int `json:"etcdDialTimeout"`
	Crontab string `json:"crontab"`
}

var (
	// 单例
	G_config *Config
)

// 加载配置
func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf Config
	)

	// 1, 把配置文件读进来
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	// 2, 做JSON反序列化
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}

	// 3, 赋值单例
	G_config = &conf

	return
}