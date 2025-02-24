package metric

import (
	"github.com/shirou/gopsutil/v4/mem"
)

// CollectMemoryMetrics collects various memory metrics and returns them along with any errors encountered.
func CollectMemoryMetrics() (*MemoryData, []CustomErr) {
	var memErrors []CustomErr
	defaultMemoryData := &MemoryData{
		TotalBytes:       0,
		AvailableBytes:   0,
		UsedBytes:        0,
		FreeBytes:        0,
		UsagePercent:     RoundFloatPtr(0, 4),
		BufferedBytes:    0,
		CachedBytes:      0,
		SharedBytes:      0,
		SwapTotal:        0,
		SwapFree:         0,
		SwapUsed:         0,
		SwapUsagePercent: RoundFloatPtr(0, 4),
	}

	// Collect virtual memory metrics
	vMem, vMemErr := mem.VirtualMemory()
	if vMemErr != nil {
		memErrors = append(memErrors, CustomErr{
			Metric: []string{"memory.virtual"},
			Error:  vMemErr.Error(),
		})
		return defaultMemoryData, memErrors
	}

	// Collect swap memory metrics
	swap, swapErr := mem.SwapMemory()
	if swapErr != nil {
		memErrors = append(memErrors, CustomErr{
			Metric: []string{"memory.swap"},
			Error:  swapErr.Error(),
		})
		// Continue with virtual memory data even if swap fails
	}

	memData := &MemoryData{
		// RAM statistics
		TotalBytes:     vMem.Total,
		AvailableBytes: vMem.Available,
		UsedBytes:      vMem.Used,
		FreeBytes:      vMem.Free,
		UsagePercent:   RoundFloatPtr(vMem.UsedPercent/100, 4),
		// Detailed RAM statistics
		BufferedBytes: vMem.Buffers,
		CachedBytes:   vMem.Cached,
		SharedBytes:   vMem.Shared,
	}

	// Add swap statistics if available
	if swap != nil {
		memData.SwapTotal = swap.Total
		memData.SwapFree = swap.Free
		memData.SwapUsed = swap.Used
		memData.SwapUsagePercent = RoundFloatPtr(swap.UsedPercent/100, 4)
	}

	return memData, memErrors
}
