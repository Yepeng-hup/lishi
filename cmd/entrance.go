package cmd

import (
	"flag"
	"os"
)

type CmdConfig struct {
	DiskDir *string
	DiskHorizon *int
	MemHorizon  *float64
	CleCacheNum *int
	CPU *int
	DDToken *string
	DDTokenFile *string
	Title *string
	ProcessNum *int
	ProcessName *string
}


func NewCmdConfig(fs *flag.FlagSet)*CmdConfig{
	diskDir := fs.String("diskDataDir", "/", "Specify the storage directory to monitor. The default is / .")
	diskHorizon := fs.Int("d", 50, "Specify how many utilization of the disk to send an alarm,The default is 50G.")
	memHorizon := fs.Float64("m", 2.0, "Specify the number of gigabytes of available memory to send an alarm. The default is 2.0G.")
	cleCacheNum := fs.Int("c", 0, "Specify the number[1,2,3] clear system cache.[0] do nothing.The default is 0 .")
	cpu := fs.Int("cpu", 20, "Specify how much the CPU is lower than to send an alarm, 20% by default.")
	ddToken := fs.String("token", "", "Specify the Token to send DingDing.")
	ddTokenFile := fs.String("token_filePath", "", "File path with token written.")
	title := fs.String("t", "", "push DindDing Keyword Title.")
	processNum := fs.Int("p", 1, "Count the total number of processes.The default is 1 .")
	processName := fs.String("n", "nil", "Process name supports wildcard.The default is nil.")
	fs.Parse(os.Args[1:])

	return &CmdConfig{
        DiskDir: diskDir,
        DiskHorizon: diskHorizon,
        MemHorizon:  memHorizon,
        CleCacheNum: cleCacheNum,
        CPU:         cpu,
        DDToken:     ddToken,
        DDTokenFile: ddTokenFile,
        Title:       title,
        ProcessNum:  processNum,
        ProcessName: processName,
    }
}
