package metric

type ApiResponse struct {
	Cpu    CpuData     `json:"cpu"`
	Memory MemoryData  `json:"memory"`
	Disk   []*DiskData `json:"disk"`
	Host   HostData    `json:"host"`
}

type CpuData struct {
	PhysicalCore int      `json:"physical_core"` // Physical cores
	LogicalCore  int      `json:"logical_core"`  // Logical cores aka Threads
	Frequency    float64  `json:"frequency"`     // Frequency in mHz
	Temperature  *float32 `json:"temperature"`   // Temperature in Celsius (nil if not available)
	FreePercent  float64  `json:"free_percent"`  // Free percentage                               //* 1 - (Total - Idle / Total)
	UsagePercent float64  `json:"usage_percent"` // Usage percentage                              //* Total - Idle / Total
}

type MemoryData struct {
	TotalBytes       uint64   `json:"total_bytes"`        // Total space in bytes
	AvailableBytes   uint64   `json:"available_bytes"`    // Available space in bytes
	UsedBytes        uint64   `json:"used_bytes"`         // Used space in bytes      //* Total - Free - Buffers - Cached
	UsagePercent     *float64 `json:"usage_percent"`      // Usage Percent            //* (Used / Total) * 100.0
	SwapTotalBytes   uint64   `json:"swap_total_bytes"`   // Total swap space in bytes
	SwapUsedBytes    uint64   `json:"swap_used_bytes"`    // Used swap space in bytes
	SwapFreeBytes    uint64   `json:"swap_free_bytes"`    // Free swap space in bytes
	SwapUsagePercent *float64 `json:"swap_usage_percent"` // Swap usage percent       //* (SwapUsed / SwapTotal) * 100.0
}

type DiskData struct {
	ReadSpeedBytes  *uint64  `json:"read_speed_bytes"`  // Read speed in bytes per second
	WriteSpeedBytes *uint64  `json:"write_speed_bytes"` // Write speed in bytes per second
	TotalBytes      *uint64  `json:"total_bytes"`       // Total space of "/" in bytes
	FreeBytes       *uint64  `json:"free_bytes"`        // Free space of "/" in bytes
	UsagePercent    *float64 `json:"usage_percent"`     // Usage Percent of "/"
}

type HostData struct {
	Hostname       string `json:"hostname"`       // Hostname
	Os             string `json:"os"`             // Operating System
	Platform       string `json:"platform"`       // Platform Name
	KernelVersion  string `json:"kernel_version"` // Kernel Version
	Uptime         uint64 `json:"uptime"`         // Uptime in seconds
	Virtualization string `json:"virtualization"` // Virtualization system and role
}

func GetAllSystemMetrics() (*ApiResponse, error) {
	cpu, cpuErr := CollectCpuMetrics()
	memory, memErr := CollectMemoryMetrics()
	disk, diskErr := CollectDiskMetrics()
	host, hostErr := GetHostInformation()

	if cpuErr != nil {
		return nil, cpuErr
	}

	if memErr != nil {
		return nil, memErr
	}

	if diskErr != nil {
		return nil, diskErr
	}

	if hostErr != nil {
		return nil, hostErr
	}

	return &ApiResponse{
		Cpu:    *cpu,
		Memory: *memory,
		Disk:   disk,
		Host:   *host,
	}, nil
}
