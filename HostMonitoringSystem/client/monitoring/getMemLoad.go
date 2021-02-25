package monitoring

import "github.com/shirou/gopsutil/v3/mem"

var(
	vmstata *mem.VirtualMemoryStat
)

func GetMemLoad() (*mem.VirtualMemoryStat,error) {
	if vmstata, err = mem.VirtualMemory(); err != nil {
		return nil, err
	}
	return vmstata, err
}
