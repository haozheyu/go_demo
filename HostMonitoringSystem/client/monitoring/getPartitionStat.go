package monitoring

import (
    "github.com/shirou/gopsutil/v3/disk"
)

var(
    partitions []disk.PartitionStat
    partition disk.PartitionStat
    partitionInfo *disk.UsageStat
)

type PartitionsInfo struct {
    Device  string `json:"device"`
    MountPoint string `json:"mount_point"`
    FsType string `json:"fs_type"`
    Total uint64 `json:"total"`
    Free uint64 `json:"free"`
    Used uint64 `json:"used"`
    UsedPercent float64 `json:"used_percent"`
    InodesUsedPercent float64 `json:"inodes_used_percent"`
}

func GetPartitions() ([]PartitionsInfo,error) {
    var (
    	pt  PartitionsInfo
    	pts []PartitionsInfo
    )
    if partitions, err = disk.Partitions(false); err != nil {
        return nil, err
    }
    for _, partition = range partitions {
        if partitionInfo, err = disk.Usage(partition.Device); err !=nil{
            return nil, err
        }
        pt.Device = partition.Device
        pt.MountPoint = partition.Mountpoint
        pt.FsType = partition.Fstype
        pt.Total = partitionInfo.Total
        pt.Free = partitionInfo.Free
        pt.Used = partitionInfo.Used
        pt.UsedPercent = partitionInfo.UsedPercent
        pt.InodesUsedPercent = partitionInfo.InodesUsedPercent
        pts = append(pts, pt)
    }
    return pts, err
}
