package collector

import (
	"github.com/shirou/gopsutil/v3/disk"

	"monitor-screen/model"
)

func CollectStorage() []model.StorageDrive {
	var drives []model.StorageDrive

	partitions, err := disk.Partitions(false)
	if err != nil {
		return drives
	}

	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}

		drives = append(drives, model.StorageDrive{
			Name:    p.Mountpoint,
			UsedGB:  float64(usage.Used) / (1024 * 1024 * 1024),
			TotalGB: float64(usage.Total) / (1024 * 1024 * 1024),
		})
	}

	enrichDiskBrands(drives)

	return drives
}

func enrichDiskBrands(drives []model.StorageDrive) {
	type diskInfo struct {
		Index uint32
		Model string
	}
	var disks []diskInfo
	err := wmiQuery("SELECT Index, Model FROM Win32_DiskDrive", &disks)
	if err != nil {
		return
	}

	for i := range drives {
		if i < len(disks) {
			model := disks[i].Model
			if model != "" {
				drives[i].Brand = model
			}
		}
	}
}
