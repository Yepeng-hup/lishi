package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sysmons/cmd"
	"sysmons/use"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Warning: please enter parameters or -h to view the usage of parameters.")
		os.Exit(1)
	}
	fs := flag.NewFlagSet("cmd", flag.ExitOnError)
	cmdParameterNameList := os.Args[1:]
	c := cmd.NewCmdConfig(fs)

	disk := &cmd.Diskd{
		DiskDir: *c.DiskDir,
	}

	alarm := &cmd.Alarm{
		T: *c.Title,
	}

	sys := &cmd.System{}


	if strings.Contains(cmdParameterNameList[0], "=") {
		os.Exit(1)
	} else {
		for _, v := range cmdParameterNameList {
			switch v {
				case "-cpu":
					use.UseCPU(sys, alarm, c)
				case "-m":
					use.UseMemory(sys, alarm, c)
				case "-diskDataDir":
					use.UseDisk(disk, sys, alarm, c)
				case "-n":
					use.UseProcess(sys, alarm, c)
				default:
					continue
			}
		}
	}

}
