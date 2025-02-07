package metric

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/v4/disk"
)

// CollectDiskMetrics collects various disk metrics and returns them.
func CollectDiskMetrics() ([]*DiskData, error) {
	var diskData []*DiskData

	partitions, partitionsErr := disk.Partitions(false)
	if partitionsErr != nil {
		return nil, partitionsErr
	}

	for _, partition := range partitions {
		diskUsage, diskUsageErr := disk.Usage(partition.Mountpoint)
		if diskUsageErr != nil {
			log.Printf("Unable to get disk usage for %s: %v", partition.Mountpoint, diskUsageErr)
			continue
		}

		readSpeed, writeSpeed, speedErr := getDiskSpeed(partition.Mountpoint)
		if speedErr != nil {
			log.Printf("Unable to get disk speed for %s: %v", partition.Mountpoint, speedErr)
			continue
		}

		diskData = append(diskData, &DiskData{
			ReadSpeedBytes:  readSpeed,
			WriteSpeedBytes: writeSpeed,
			TotalBytes:      &diskUsage.Total,
			FreeBytes:       &diskUsage.Free,
			UsagePercent:    RoundFloatPtr(diskUsage.UsedPercent/100, 4),
		})
	}

	return diskData, nil
}

// getDiskSpeed calculates the read and write speed for the given mountpoint.
func getDiskSpeed(mountpoint string) (*uint64, *uint64, error) {
	ioCountersStart, err := disk.IOCounters(mountpoint)
	if err != nil {
		return nil, nil, err
	}

	time.Sleep(1 * time.Second)

	ioCountersEnd, err := disk.IOCounters(mountpoint)
	if err != nil {
		return nil, nil, err
	}

	readSpeed := ioCountersEnd[mountpoint].ReadBytes - ioCountersStart[mountpoint].ReadBytes
	writeSpeed := ioCountersEnd[mountpoint].WriteBytes - ioCountersStart[mountpoint].WriteBytes

	return &readSpeed, &writeSpeed, nil
}
