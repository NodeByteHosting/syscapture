package metric

import "time"

// MetricsSlice represents a slice of Metric interfaces.
type MetricsSlice []Metric

func (m MetricsSlice) isMetric() {}

// Metric interface to be implemented by all metric types.
type Metric interface {
	isMetric()
}

// APIResponse represents the structure of the API response.
type APIResponse struct {
	Data   Metric      `json:"data"`
	Errors []CustomErr `json:"errors"`
}

// AllMetrics represents all collected system metrics.
type AllMetrics struct {
	CPU     CPUData      `json:"cpu"`
	Memory  MemoryData   `json:"memory"`
	Disk    MetricsSlice `json:"disk"`
	Host    HostData     `json:"host"`
	Network NetworkData  `json:"network"`
}

func (a AllMetrics) isMetric() {}

// CustomErr represents a custom error structure.
type CustomErr struct {
	Metric []string `json:"metric"`
	Error  string   `json:"err"`
}

// CPUData represents the collected CPU metrics.
type CPUData struct {
	PhysicalCore     int       `json:"physical_core"`     // Physical cores
	LogicalCore      int       `json:"logical_core"`      // Logical cores aka Threads
	Frequency        float64   `json:"frequency"`         // Frequency in mHz
	CurrentFrequency int       `json:"current_frequency"` // Current Frequency in mHz
	Temperature      []float32 `json:"temperature"`       // Temperature in Celsius (nil if not available)
	FreePercent      float64   `json:"free_percent"`      // Free percentage
	UsagePercent     float64   `json:"usage_percent"`     // Usage percentage
}

func (c CPUData) isMetric() {}

// MemoryData represents the collected memory metrics.
type MemoryData struct {
	TotalBytes     uint64   `json:"total_bytes"`     // Total space in bytes
	AvailableBytes uint64   `json:"available_bytes"` // Available space in bytes
	UsedBytes      uint64   `json:"used_bytes"`      // Used space in bytes
	UsagePercent   *float64 `json:"usage_percent"`   // Usage Percent
}

func (m MemoryData) isMetric() {}

// DiskData represents disk metrics for a single partition
type DiskData struct {
	Device       string   `json:"device"`
	Mountpoint   string   `json:"mountpoint"`
	TotalBytes   *uint64  `json:"total_bytes"`
	FreeBytes    *uint64  `json:"free_bytes"`
	UsagePercent *float64 `json:"usage_percent"`
	IOStats      *IOStats `json:"io_stats,omitempty"`
}

// IOStats represents disk I/O statistics
type IOStats struct {
	ReadCount      uint64        `json:"read_count"`
	WriteCount     uint64        `json:"write_count"`
	ReadBytes      uint64        `json:"read_bytes"`
	WriteBytes     uint64        `json:"write_bytes"`
	ReadTime       time.Duration `json:"read_time"`
	WriteTime      time.Duration `json:"write_time"`
	IopsInProgress uint64        `json:"iops_in_progress"`
}

func (d DiskData) isMetric() {}

// HostData represents the collected host information.
type HostData struct {
	Os            string    `json:"os"`
	Platform      string    `json:"platform"`
	KernelVersion string    `json:"kernel_version"`
	Hostname      string    `json:"hostname"`
	Uptime        uint64    `json:"uptime"`
	BootTime      time.Time `json:"boot_time"`
	ProcCount     uint64    `json:"process_count"`
	Users         []string  `json:"users"`
}

func (h HostData) isMetric() {}

// NetworkData represents the collected network metrics.
type NetworkData struct {
	BytesSent   uint64 `json:"bytes_sent"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
}

func (n NetworkData) isMetric() {}

// GetAllSystemMetrics collects all system metrics and returns them along with any errors encountered.
func GetAllSystemMetrics() (AllMetrics, []CustomErr) {
	cpu, cpuErr := CollectCPUMetrics()
	memory, memErr := CollectMemoryMetrics()
	disk, diskErr := CollectDiskMetrics()
	host, hostErr := GetHostInformation()
	network, networkErr := CollectNetworkMetrics()

	var errors []CustomErr

	if cpuErr != nil {
		errors = append(errors, cpuErr...)
	}

	if memErr != nil {
		errors = append(errors, memErr...)
	}

	if diskErr != nil {
		errors = append(errors, diskErr...)
	}

	if hostErr != nil {
		errors = append(errors, hostErr...)
	}

	if networkErr != nil {
		errors = append(errors, networkErr...)
	}

	return AllMetrics{
		CPU:     *cpu,
		Memory:  *memory,
		Disk:    disk,
		Host:    *host,
		Network: *network,
	}, errors
}
