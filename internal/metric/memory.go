package metric

import (
	"github.com/shirou/gopsutil/v4/mem"
)

// CollectMemoryMetrics collects various memory metrics and returns them.
func CollectMemoryMetrics() (*MemoryData, error) {
	vMem, vMemErr := mem.VirtualMemory()
	if vMemErr != nil {
		return nil, vMemErr
	}

	swapMem, swapMemErr := mem.SwapMemory()
	if swapMemErr != nil {
		return nil, swapMemErr
	}

	return &MemoryData{
		TotalBytes:       vMem.Total,
		AvailableBytes:   vMem.Available,
		UsedBytes:        vMem.Used,
		UsagePercent:     RoundFloatPtr(vMem.UsedPercent/100, 4),
		SwapTotalBytes:   swapMem.Total,
		SwapUsedBytes:    swapMem.Used,
		SwapFreeBytes:    swapMem.Free,
		SwapUsagePercent: RoundFloatPtr(swapMem.UsedPercent/100, 4),
	}, nil
}
