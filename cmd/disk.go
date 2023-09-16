package cmd

import (
	"github.com/shirou/gopsutil/disk"
)


func (s *System)DiskSy(dataDir string)*DiskInfo {
	d, _ := disk.Usage(dataDir)
	disks := &DiskInfo{
		DiskUsed: int(d.Used/1024/1024/1024),
		DiskFree: float64(d.Free/1024/1024/1024),
		DiskTotal: float64(d.Total/1024/1024/1024),
	}
	return disks
}

