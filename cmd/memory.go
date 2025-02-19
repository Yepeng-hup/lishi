package cmd

import (
	"github.com/shirou/gopsutil/mem"
)

const gm uint64 = 1074000000

func (s *System) MemSy() *MemInfo {
	m, _ := mem.VirtualMemory()
	if m.Total < gm {
		total := float64(m.Total/1024/1024) / float64(1024)
		active := float64(m.Active/1024/1024) / float64(1024)
		free := float64(m.Free/1024/1024) / float64(1024)
		mems := &MemInfo{
			MemTotal: total,
			MemUsed:  active,
			MemFree:  free,
		}
		return mems
	} else {
		total := float64(m.Total) / 1024 / 1024 / 1024
		active := float64(m.Active) / 1024 / 1024 / 1024
		free := float64(m.Free) / 1024 / 1024 / 1024
		mems := &MemInfo{
			MemTotal: total,
			MemUsed:  active,
			MemFree:  free,
		}
		return mems
	}
}
