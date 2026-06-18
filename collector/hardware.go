package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type BridgeSensor struct {
	Name  string  `json:"n"`
	Type  string  `json:"t"`
	Value float64 `json:"v"`
}

type BridgeOutput map[string][]BridgeSensor

var bridgeProcess *os.Process

func EnsureDriverLoaded(lhmExe string) {
	if lhmExe == "" {
		lhmExe = "lhm\\LibreHardwareMonitor.exe"
	}
	if !filepath.IsAbs(lhmExe) {
		if exe, err := os.Executable(); err == nil {
			lhmExe = filepath.Join(filepath.Dir(exe), lhmExe)
		}
	}
	if _, err := os.Stat(lhmExe); err != nil {
		return
	}
	cmd := exec.Command(lhmExe)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Start()
	time.Sleep(1500 * time.Millisecond)
	cmd.Process.Kill()
}

func StartBridge(exePath string) {
	if !filepath.IsAbs(exePath) {
		if exe, err := os.Executable(); err == nil {
			exePath = filepath.Join(filepath.Dir(exe), exePath)
		}
	}
	if bridgeProcess != nil {
		return
	}
	cmd := exec.Command(exePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		log.Printf("bridge: failed to start %s: %v", exePath, err)
		return
	}
	bridgeProcess = cmd.Process
	time.Sleep(1 * time.Second)
}

func StopBridge() {
	if bridgeProcess != nil {
		pid := bridgeProcess.Pid
		bridgeProcess = nil
		killPID(pid)
	}
	killByName("sensor_bridge.exe")
}

func killPID(pid int) {
	cmd := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", pid))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()
}

func killByName(name string) {
	cmd := exec.Command("taskkill", "/F", "/IM", name)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()
}

func RunSensorBridge(exePath string) *BridgeOutput {
	if bridgeProcess != nil && isBridgeAlive() {
		return queryBridge()
	}

	if !filepath.IsAbs(exePath) {
		if exe, err := os.Executable(); err == nil {
			exePath = filepath.Join(filepath.Dir(exe), exePath)
		}
	}

	if bridgeProcess != nil {
		killPID(bridgeProcess.Pid)
		bridgeProcess = nil
	}

	log.Printf("bridge: starting %s", exePath)
	cmd := exec.Command(exePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Start(); err != nil {
		log.Printf("bridge: error starting %s: %v", exePath, err)
		return nil
	}
	bridgeProcess = cmd.Process

	time.Sleep(2 * time.Second)

	if !isBridgeAlive() {
		log.Printf("bridge: process exited immediately, falling back to one-shot mode")
		bridgeProcess = nil
		return runBridgeOneShot(exePath)
	}

	return queryBridge()
}

func isBridgeAlive() bool {
	if bridgeProcess == nil {
		return false
	}
	resp, err := http.Get("http://127.0.0.1:16533/")
	if err != nil {
		return false
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()
	return true
}

func QueryBridgeIfAlive(exePath string) *BridgeOutput {
	if bridgeProcess == nil || !isBridgeAlive() {
		return nil
	}
	return queryBridge()
}

func queryBridge() *BridgeOutput {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://127.0.0.1:16533/")
	if err != nil {
		log.Printf("bridge: HTTP query failed: %v", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("bridge: read body failed: %v", err)
		return nil
	}

	var output BridgeOutput
	if err := json.Unmarshal(body, &output); err != nil {
		log.Printf("bridge: parse error: %v", err)
		return nil
	}

	log.Printf("bridge: got %d hardware entries via HTTP", len(output))
	return &output
}

func runBridgeOneShot(exePath string) *BridgeOutput {
	log.Printf("bridge: running one-shot %s", exePath)
	cmd := exec.Command(exePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		log.Printf("bridge: error running %s: %v", exePath, err)
		return nil
	}

	var output BridgeOutput
	if err := json.Unmarshal(out, &output); err != nil {
		log.Printf("bridge: parse error: %v", err)
		return nil
	}

	log.Printf("bridge: got %d hardware entries", len(output))
	return &output
}

func findSensor(output *BridgeOutput, hwType, sensorName, sensorType string) float64 {
	if output == nil {
		return 0
	}

	for key, sensors := range *output {
		if !strings.Contains(strings.ToLower(key), strings.ToLower(hwType)) {
			continue
		}

		for _, s := range sensors {
			sName := strings.ToLower(s.Name)
			sType := strings.ToLower(s.Type)

			if sensorType != "" && !strings.Contains(sType, strings.ToLower(sensorType)) {
				continue
			}

			if strings.Contains(sName, strings.ToLower(sensorName)) {
				return s.Value
			}
		}
	}
	return 0
}
