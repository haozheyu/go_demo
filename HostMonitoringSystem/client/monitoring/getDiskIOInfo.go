package monitoring

import (
	"github.com/shirou/gopsutil/v3/disk"
)

var(
	diskio map[string]disk.IOCountersStat
)

func GetDiskIOInfo() ([]disk.IOCountersStat,error){
	var diskIOlist []disk.IOCountersStat
	if diskio, err = disk.IOCounters(); err !=nil {
		return nil, err
	}
	for _, stat := range diskio {
		diskIOlist= append(diskIOlist, stat)
	}
	return diskIOlist, err
}
