package cmd

type (
	MemInfo struct {
		MemTotal float64
		MemUsed float64
		MemFree float64
	}

	DiskInfo struct {
		DiskTotal float64
		DiskUsed int
		DiskFree float64
	}

	Diskd struct {
		DiskDir string
	}

	Process struct {
		Response bool
	}

	ProcessTime struct {
		RunTime int
	}

	System struct {}

	Alarm struct {
		T string
	}

)
