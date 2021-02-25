package monitoring

import (
	"github.com/shirou/gopsutil/v3/net"
)

var (
	iostat []net.IOCountersStat
	connections []net.ConnectionStat
)

type NetStat struct {
	NetConnectionCount int64
	TimeWaitNumber int64
	ESTABLISHEDNumber int64
	BytesSent uint64
	BytesRecv uint64
	PacketsRecv uint64
	PacketsSent uint64
	Dropin uint64
	Dropout uint64
}

func GetNetStat() (*NetStat,error) {
	var (
		netstat          NetStat
		netConnection    net.ConnectionStat
		totalCount       = 0
		establishedCount = 0
		timewaitCount    = 0
	)

	if iostat, err = net.IOCounters(false); err !=nil{
		return nil, err
	}
	netstat.BytesSent =  iostat[0].BytesSent
	netstat.BytesRecv = iostat[0].BytesRecv
	netstat.PacketsRecv = iostat[0].PacketsRecv
	netstat.PacketsSent = iostat[0].PacketsSent
	netstat.Dropin = iostat[0].Dropin
	netstat.Dropout = iostat[0].Dropout

	if connections, err = net.Connections("all"); err !=nil {
		return nil, err
	}
    for totalCount,netConnection = range connections {
    	if netConnection.Status == "ESTABLISHED" {
    		establishedCount ++
		}else if netConnection.Status == "TIME_WAIT"{
			timewaitCount ++
		}
	}
	netstat.NetConnectionCount = int64(totalCount)
	netstat.ESTABLISHEDNumber = int64(establishedCount)
	netstat.TimeWaitNumber = int64(timewaitCount)
	return &netstat, err
}

