package monitoring

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"time"
)

type CPULoad struct {
	TotalPercent float64
	PerPercents []float64
	CPUIOWait float64
}

func GetCpuLoad() (*CPULoad,error) {
	var (
		total    []float64
		per      []float64
		cpuload  CPULoad
		timeStat []cpu.TimesStat
	)
	if per, err = cpu.Percent(1* time.Second,true); err !=nil {
		return nil, err
	}
	if total, err = cpu.Percent(1* time.Second,false); err !=nil {
		return nil, err
	}
	if timeStat, err = cpu.Times(false); err !=nil {
		return nil, err
	}
	cpuload.TotalPercent = total[0]
	cpuload.PerPercents = per
	cpuload.CPUIOWait = timeStat[0].Iowait
	return &cpuload, err
}