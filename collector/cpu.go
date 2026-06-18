package collector

import (
	"fmt"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"

	"monitor-screen/model"
)

func CollectCPU(prevUsage *float64) model.CPUData {
	d := model.CPUData{}

	info, err := cpu.Info()
	if err == nil && len(info) > 0 {
		d.Model = strings.TrimSpace(info[0].ModelName)
		mhz := info[0].Mhz
		if mhz >= 1000 {
			d.MaxFreq = fmt.Sprintf("%.1f GHz", mhz/1000)
		} else if mhz > 0 {
			d.MaxFreq = fmt.Sprintf("%.0f MHz", mhz)
		}
		if info[0].CacheSize > 0 {
			cacheKB := info[0].CacheSize
			if cacheKB >= 1024 {
				d.CacheL3 = fmt.Sprintf("%.0f MB", float64(cacheKB)/1024)
			} else {
				d.CacheL3 = fmt.Sprintf("%d KB", cacheKB)
			}
		}
	}

	physical, _ := cpu.Counts(false)
	logical, _ := cpu.Counts(true)
	if physical > 0 && logical > 0 {
		d.CoresThreads = fmt.Sprintf("%dC / %dT", physical, logical)
	}

	percentages, err := cpu.Percent(500*time.Millisecond, false)
	if err == nil && len(percentages) > 0 {
		d.Usage = percentages[0]
		if prevUsage != nil {
			*prevUsage = d.Usage
		}
	}

	perCore, err := cpu.Percent(0, true)
	if err == nil && len(perCore) > 0 {
		d.PerCoreUsage = perCore
	}

	return d
}
