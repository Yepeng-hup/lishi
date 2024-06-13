package cmd

import (
	"fmt"
	"sysmons/core"
)

func CleCache(cacheSignal int)error{
	switch  {
	case cacheSignal == 1:
		_, err := core.RunCommand(`sync && echo 1 > /proc/sys/vm/drop_caches`)
		return err
	case cacheSignal == 2:
		_, err := core.RunCommand(`sync && echo 2 > /proc/sys/vm/drop_caches`)
		return err
	case cacheSignal == 3:
		_, err := core.RunCommand(`sync && echo 3 > /proc/sys/vm/drop_caches`)
		return err
	case cacheSignal == 0:
		return nil
	default:
		return fmt.Errorf("%s","Signal type error and to [1, 2, 3].")
	}
}
