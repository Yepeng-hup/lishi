package cmd

import (
	"github.com/shirou/gopsutil/cpu"
	"time"
)

const cpuAll int = 100

func (s *System) CpuSy()int{
	c2, _ := cpu.Percent(time.Duration(time.Second), false)
	cpuFrees := cpuAll - int(c2[0])
	return cpuFrees
}
