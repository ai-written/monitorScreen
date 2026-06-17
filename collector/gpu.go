package collector

import (
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"monitor-screen/model"
)

func CollectGPU() model.GPUData {
	d := model.GPUData{}

	cmd := exec.Command("nvidia-smi",
		"--query-gpu=name,memory.total,temperature.gpu,utilization.gpu,clocks.current.graphics,memory.used,fan.speed,power.draw",
		"--format=csv,noheader,nounits",
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	out, err := cmd.Output()
	if err != nil {
		return d
	}

	line := strings.TrimSpace(string(out))
	fields := strings.Split(line, ",")
	if len(fields) < 8 {
		return d
	}

	for i := range fields {
		fields[i] = strings.TrimSpace(fields[i])
	}

	d.Model = fields[0]

	if v, err := strconv.ParseFloat(fields[1], 64); err == nil {
		d.MemTotal = v
		d.VRAMSpec = formatVRAM(v)
	}
	if v, err := strconv.ParseFloat(fields[2], 64); err == nil {
		d.Temp = v
	}
	if v, err := strconv.ParseFloat(fields[3], 64); err == nil {
		d.Usage = v
	}
	if v, err := strconv.ParseFloat(fields[4], 64); err == nil {
		d.Clock = v
	}
	if v, err := strconv.ParseFloat(fields[5], 64); err == nil {
		d.MemUsed = v
	}
	if fields[6] != "[N/A]" {
		if v, err := strconv.ParseFloat(fields[6], 64); err == nil {
			d.FanSpeed = v
		}
	}
	if fields[7] != "[N/A]" {
		if v, err := strconv.ParseFloat(fields[7], 64); err == nil {
			d.Power = v
		}
	}

	return d
}

func formatVRAM(mb float64) string {
	if mb >= 1024 {
		return strconv.FormatFloat(mb/1024, 'f', -1, 64) + " GB GDDR6X"
	}
	return strconv.FormatFloat(mb, 'f', 0, 64) + " MB"
}
