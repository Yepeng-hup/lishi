package cmd

import (
	"github.com/shirou/gopsutil/mem"
)


func (s *System) MemSy() *MemInfo {
	m, _ := mem.VirtualMemory()
	mems := &MemInfo{
		MemTotal: float64(m.Total/1024/1024/1024),
		MemUsed: float64(m.Active/1024/1024/1024),
		MemFree: float64(m.Free/1024/1024/1024),
	}
	return mems
}