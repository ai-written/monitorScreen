package collector

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/mem"

	"monitor-screen/model"
)

func CollectMemory() model.MemoryData {
	d := model.MemoryData{}

	v, err := mem.VirtualMemory()
	if err == nil {
		d.UsedGB = float64(v.Used) / (1024 * 1024 * 1024)
		d.TotalGB = float64(v.Total) / (1024 * 1024 * 1024)
	}

	s, err := mem.SwapMemory()
	if err == nil {
		d.SwapTotal = float64(s.Total) / (1024 * 1024 * 1024)
		d.SwapUsed = float64(s.Used) / (1024 * 1024 * 1024)
	}

	collectMemoryWMI(&d)

	return d
}

func collectMemoryWMI(d *model.MemoryData) {
	type memChip struct {
		Capacity            uint64
		Speed               uint32
		MemoryType          uint16
		SMBIOSMemoryType    uint16
		ConfiguredClockSpeed uint32
		Manufacturer        string
	}

	var chips []memChip
	q := "SELECT Capacity, Speed, MemoryType, SMBIOSMemoryType, ConfiguredClockSpeed, Manufacturer FROM Win32_PhysicalMemory"
	err := wmiQuery(q, &chips)
	if err != nil {
		return
	}

	if len(chips) > 0 {
		mt := chips[0].SMBIOSMemoryType
		if mt == 0 {
			mt = chips[0].MemoryType
		}

		switch mt {
		case 20:
			d.Type = "DDR"
		case 21:
			d.Type = "DDR2"
		case 22:
			d.Type = "DDR2"
		case 24:
			d.Type = "DDR3"
		case 26:
			d.Type = "DDR4"
		case 27:
			d.Type = "DDR4"
		case 34:
			d.Type = "DDR5"
		default:
			if mt > 0 {
				d.Type = fmt.Sprintf("DDR(%d)", mt)
			}
		}

		speed := chips[0].Speed
		if speed == 0 {
			speed = chips[0].ConfiguredClockSpeed
		}
		if speed > 0 {
			d.Frequency = fmt.Sprintf("%d MHz", speed)
		}

		d.Brand = chips[0].Manufacturer

		switch {
		case len(chips) >= 4:
			d.Channel = "四通道"
		case len(chips) >= 2:
			d.Channel = "双通道"
		case len(chips) == 1:
			d.Channel = "单通道"
		}
	}
}
