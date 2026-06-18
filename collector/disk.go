package collector

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/disk"

	"monitor-screen/model"
)

var (
	diskIOMu       sync.Mutex
	lastDiskBytes  = make(map[string]ioCounters)
	lastDiskTime   time.Time
	cachedDiskDown float64
	cachedDiskUp   float64
)

type ioCounters struct{ read, write uint64 }

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

func CollectDiskIORate() model.DiskIOData {
	diskIOMu.Lock()
	defer diskIOMu.Unlock()

	counters, err := disk.IOCounters()
	if err != nil {
		return model.DiskIOData{ReadMBps: cachedDiskDown, WriteMBps: cachedDiskUp}
	}

	now := time.Now()
	current := make(map[string]ioCounters)
	var totalRead, totalWrite uint64

	for name, c := range counters {
		current[name] = ioCounters{c.ReadBytes, c.WriteBytes}
		totalRead += c.ReadBytes
		totalWrite += c.WriteBytes
	}

	if len(lastDiskBytes) > 0 && !lastDiskTime.IsZero() {
		elapsed := now.Sub(lastDiskTime).Seconds()
		if elapsed >= 0.5 {
			var diffRead, diffWrite uint64
			for name, cur := range current {
				if prev, ok := lastDiskBytes[name]; ok {
					diffRead += cur.read - prev.read
					diffWrite += cur.write - prev.write
				}
			}
			cachedDiskDown = float64(diffRead) / elapsed / (1024 * 1024)
			cachedDiskUp = float64(diffWrite) / elapsed / (1024 * 1024)
		}
	}

	lastDiskBytes = current
	lastDiskTime = now
	return model.DiskIOData{ReadMBps: cachedDiskDown, WriteMBps: cachedDiskUp}
}
