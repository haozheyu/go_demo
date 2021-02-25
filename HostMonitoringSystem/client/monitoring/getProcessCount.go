package monitoring

import (
	"github.com/shirou/gopsutil/v3/process"
)
var(
	runPids []int32
)

func RunProcessCount() int32 {
	runPids, err = process.Pids()
	return int32(len(runPids))
}