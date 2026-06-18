package collector

import (
	"sort"
	"strings"

	"monitor-screen/model"
)

func EnrichFromLHM(bridgeOutput *BridgeOutput, cpu *model.CPUData, gpu *model.GPUData, drives []model.StorageDrive) []model.FanInfo {
	if bridgeOutput == nil {
		return nil
	}

	enrichCPUBridge(bridgeOutput, cpu)
	enrichGPUBridge(bridgeOutput, gpu, cpu.Usage)
	enrichStorageBridge(bridgeOutput, drives)
	fans := collectFansBridge(bridgeOutput)

	return fans
}

func enrichCPUBridge(output *BridgeOutput, d *model.CPUData) {
	d.PackageTemp = findCPUValue(output, "temperature", "package", "cpu package", "core", "cpu core", "cpu")
	if d.PackageTemp == 0 {
		d.PackageTemp = findCPUValue(output, "temperature", "temp")
	}

	d.ClockSpeed = findCPUValue(output, "clock", "clock", "core", "cpu core")
	if d.ClockSpeed == 0 {
		d.ClockSpeed = findCPUValue(output, "clock", "speed")
	}

	d.Vcore = findCPUValue(output, "voltage", "vcore", "core", "cpu core", "cpu vcore")
	if d.Vcore == 0 {
		d.Vcore = findCPUValue(output, "voltage", "voltage")
	}

	d.FanSpeed = findCPUValue(output, "control", "fan", "cpu fan", "pump")
	if d.FanSpeed == 0 {
		d.FanSpeed = findCPUValue(output, "fan", "fan")
	}
	if d.FanSpeed == 0 {
		d.FanSpeed = findAnyFan(output)
	}

	d.Power = findCPUValue(output, "power", "package", "cpu package", "cpu")

	d.MBTemp = findMBTemp(output)
	d.VRMTemp = findVRMTemp(output)
}

func findMBTemp(output *BridgeOutput) float64 {
	for key, sensors := range *output {
		lowerKey := strings.ToLower(key)
		if !strings.Contains(lowerKey, "motherboard") {
			continue
		}
		for _, s := range sensors {
			if strings.Contains(strings.ToLower(s.Type), "temperature") && s.Value > 0 {
				return s.Value
			}
		}
	}
	return 0
}

func findVRMTemp(output *BridgeOutput) float64 {
	for key, sensors := range *output {
		lowerKey := strings.ToLower(key)
		if !strings.Contains(lowerKey, "motherboard") && !strings.Contains(lowerKey, "vrm") {
			continue
		}
		for _, s := range sensors {
			sn := strings.ToLower(s.Name)
			st := strings.ToLower(s.Type)
			if strings.Contains(st, "temperature") && (strings.Contains(sn, "vrm") || strings.Contains(sn, "mos")) && s.Value > 0 {
				return s.Value
			}
		}
	}
	return 0
}

func findVRAMType(sensors []BridgeSensor) string {
	for _, s := range sensors {
		sn := strings.ToLower(s.Name)
		if strings.Contains(sn, "memory type") || strings.Contains(sn, "vram type") {
			return s.Name
		}
	}
	return ""
}

func findCPUValue(output *BridgeOutput, sensorType string, names ...string) float64 {
	for key, sensors := range *output {
		if !strings.Contains(strings.ToLower(key), "cpu") {
			continue
		}

		for _, s := range sensors {
			sName := strings.ToLower(s.Name)
			sType := strings.ToLower(s.Type)

			if !strings.Contains(sType, strings.ToLower(sensorType)) {
				continue
			}

			for _, name := range names {
				if strings.Contains(sName, strings.ToLower(name)) {
					if s.Value > 0 {
						return s.Value
					}
				}
			}

			if len(names) == 1 && strings.Contains(sName, names[0]) && s.Value > 0 {
				return s.Value
			}
		}
	}
	return 0
}

func enrichGPUBridge(output *BridgeOutput, d *model.GPUData, cpuUsage float64) {
	gpuKeys := make([]string, 0, len(*output))
	for key := range *output {
		if strings.Contains(strings.ToLower(key), "gpu") {
			gpuKeys = append(gpuKeys, key)
		}
	}

	sortGPUKeys(gpuKeys)

	for _, key := range gpuKeys {
		sensors := (*output)[key]

		for _, s := range sensors {
			sName := strings.ToLower(s.Name)
			sType := strings.ToLower(s.Type)

			if d.Temp == 0 && strings.Contains(sName, "core") && strings.Contains(sType, "temperature") {
				d.Temp = s.Value
			}
			if strings.Contains(sName, "memory junction") && strings.Contains(sType, "temperature") {
				if d.Temp == 0 {
					d.Temp = s.Value
				}
			}
			if d.Clock == 0 && strings.Contains(sType, "clock") &&
				!strings.Contains(sName, "memory") {
				d.Clock = s.Value
			}
			if d.Usage == 0 && strings.Contains(sType, "load") &&
				(strings.Contains(sName, "core") || strings.Contains(sName, "3d")) {
				d.Usage = s.Value
			}
			if d.Power == 0 && strings.Contains(sType, "power") {
				d.Power = s.Value
			}
			if d.FanSpeed == 0 && strings.Contains(sType, "fan") &&
				strings.Contains(sName, "fan") && s.Value > 0 {
				d.FanSpeed = s.Value
			}
			if d.FanSpeed == 0 && strings.Contains(sName, "fan") &&
				strings.Contains(sType, "control") && s.Value > 0 {
				d.FanSpeed = s.Value
			}
			if d.FanSpeed == 0 && strings.Contains(sName, "fan") &&
				strings.Contains(sType, "fan") {
				d.FanSpeed = s.Value
			}
			if d.MemUsed == 0 && strings.Contains(sType, "smalldata") &&
				strings.Contains(sName, "used") {
				d.MemUsed = s.Value
			}
			if d.MemTotal == 0 && strings.Contains(sType, "smalldata") &&
				strings.Contains(sName, "total") {
				d.MemTotal = s.Value
			}
		}

		if d.Model == "" {
			parts := strings.SplitN(key, "|", 2)
			if len(parts) == 2 {
				d.Model = strings.TrimSpace(parts[1])
			}
		}

		if d.VRAMType == "" {
			d.VRAMType = findVRAMType(sensors)
		}
	}
}

func sortGPUKeys(keys []string) {
	gpuPriority := func(key string) int {
		k := strings.ToLower(key)
		if strings.Contains(k, "nvidia") {
			return 0
		}
		if strings.Contains(k, "amd") || strings.Contains(k, "radeon") {
			return 1
		}
		return 2
	}
	sort.Slice(keys, func(i, j int) bool {
		return gpuPriority(keys[i]) < gpuPriority(keys[j])
	})
}

func enrichStorageBridge(output *BridgeOutput, drives []model.StorageDrive) {
	for key, sensors := range *output {
		lowerKey := strings.ToLower(key)
		if !strings.Contains(lowerKey, "storage") &&
			!strings.Contains(lowerKey, "hdd") &&
			!strings.Contains(lowerKey, "ssd") &&
			!strings.Contains(lowerKey, "nvme") {
			continue
		}

		for _, s := range sensors {
			sName := strings.ToLower(s.Name)
			if strings.Contains(sName, "temperature") || strings.Contains(sName, "temp") {
				for i := range drives {
					if drives[i].Temp == 0 && s.Value > 0 {
						drives[i].Temp = s.Value
						break
					}
				}
			}
		}

		_ = key
	}
}

func collectFansBridge(output *BridgeOutput) []model.FanInfo {
	var fans []model.FanInfo

	for key, sensors := range *output {
		for _, s := range sensors {
			sName := strings.ToLower(s.Name)

			if !strings.Contains(sName, "fan") {
				continue
			}
			if strings.Contains(sName, "control") || strings.Contains(sName, "max") || s.Value <= 0 {
				continue
			}

			fan := model.FanInfo{
				Name: s.Name,
				RPM:  s.Value,
			}

			for _, s2 := range sensors {
				s2Name := strings.ToLower(s2.Name)
				if strings.Contains(s2Name, "control") || strings.Contains(s2Name, "percent") {
					fan.Percent = s2.Value
					break
				}
			}

			fans = append(fans, fan)
		}

		_ = key
	}

	return fans
}

func TotalPowerFromLHM(bridgeOutput *BridgeOutput) float64 {
	if bridgeOutput == nil {
		return 0
	}
	var total float64
	for _, sensors := range *bridgeOutput {
		for _, s := range sensors {
			if strings.Contains(strings.ToLower(s.Type), "power") && s.Value > 0 {
				total += s.Value
			}
		}
	}
	return total
}

func findAnyFan(output *BridgeOutput) float64 {
	for _, sensors := range *output {
		for _, s := range sensors {
			sType := strings.ToLower(s.Type)
			sName := strings.ToLower(s.Name)
			if strings.Contains(sType, "fan") && !strings.Contains(sName, "control") && !strings.Contains(sName, "max") && s.Value > 0 {
				return s.Value
			}
		}
	}
	return 0
}
