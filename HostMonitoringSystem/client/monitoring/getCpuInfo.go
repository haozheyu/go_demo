package monitoring

import (
	"github.com/shirou/gopsutil/v3/cpu"
)

var(
	G_CpuInfo *CpuInfo //注册一个单利
	cpuinfo []cpu.InfoStat
	err error
)

type CpuInfo struct {
	Cores int32 `json:"cores"`
	ModeName string `json:"mode_name"`
	Mhz float64 `json:"mhz"`
	CacheSize int32 `json:"cache_size"`
}

func GetCpuInfo() (*CpuInfo,error){
	var c CpuInfo
	if cpuinfo, err = cpu.Info(); err !=nil {
		return nil, err
	}
	c.Cores = cpuinfo[0].Cores
	c.ModeName = cpuinfo[0].ModelName
	c.Mhz = cpuinfo[0].Mhz
	c.CacheSize = cpuinfo[0].CacheSize
	G_CpuInfo = &c
	return G_CpuInfo, err
}