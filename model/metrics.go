package model

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Collector CollectorConfig `yaml:"collector"`
	LHM       LHMConfig       `yaml:"lhm"`
}

type CollectorConfig struct {
	Interval int `yaml:"interval"`
}

type LHMConfig struct {
	Enabled   string `yaml:"enabled"`
	BridgeExe string `yaml:"bridge_exe"`
	LHMExe    string `yaml:"lhm_exe"`
}

type DashboardData struct {
	System     SystemInfo     `json:"system"`
	CPU        CPUData        `json:"cpu"`
	GPU        GPUData        `json:"gpu"`
	Memory     MemoryData     `json:"memory"`
	Storage    []StorageDrive `json:"storage"`
	Fans       []FanInfo      `json:"fans"`
	TotalPower float64        `json:"total_power"`
	Network    NetworkData    `json:"network"`
	DiskIO     DiskIOData     `json:"disk_io"`
}

type SystemInfo struct {
	Time         string `json:"time"`
	Date         string `json:"date"`
	Uptime       string `json:"uptime"`
	ProcessCount int    `json:"process_count"`
	ThreadCount  int    `json:"thread_count"`
	IPAddress    string `json:"ip_address"`
	DisplayInfo  string `json:"display_info"`
}

type CPUData struct {
	Model        string    `json:"model"`
	MaxFreq      string    `json:"max_freq"`
	CoresThreads string    `json:"cores_threads"`
	CacheL3      string    `json:"cache_l3"`
	PackageTemp  float64   `json:"package_temp"`
	Usage        float64   `json:"usage"`
	ClockSpeed   float64   `json:"clock_speed"`
	Vcore        float64   `json:"vcore"`
	FanSpeed     float64   `json:"fan_speed"`
	Power        float64   `json:"power"`
	PerCoreUsage []float64 `json:"per_core_usage"`
	MBTemp       float64   `json:"mb_temp"`
	VRMTemp      float64   `json:"vrm_temp"`
}

type GPUData struct {
	Model       string  `json:"model"`
	VRAMSpec    string  `json:"vram_spec"`
	VRAMType    string  `json:"vram_type"`
	PCIEVersion string  `json:"pcie_version"`
	Temp        float64 `json:"temp"`
	Usage       float64 `json:"usage"`
	Clock       float64 `json:"clock"`
	MemUsed     float64 `json:"mem_used"`
	MemTotal    float64 `json:"mem_total"`
	FanSpeed    float64 `json:"fan_speed"`
	Power       float64 `json:"power"`
}

type MemoryData struct {
	UsedGB    float64 `json:"used_gb"`
	TotalGB   float64 `json:"total_gb"`
	Type      string  `json:"type"`
	Frequency string  `json:"frequency"`
	Channel   string  `json:"channel"`
	Brand     string  `json:"brand"`
}

type StorageDrive struct {
	Name    string  `json:"name"`
	Brand   string  `json:"brand"`
	Temp    float64 `json:"temp"`
	UsedGB  float64 `json:"used_gb"`
	TotalGB float64 `json:"total_gb"`
}

type FanInfo struct {
	Name    string  `json:"name"`
	RPM     float64 `json:"rpm"`
	Percent float64 `json:"percent"`
}

type NetworkData struct {
	DownloadMbps    float64 `json:"download_mbps"`
	UploadMbps      float64 `json:"upload_mbps"`
	ConnectionCount int     `json:"connection_count"`
}

type DiskIOData struct {
	ReadMBps  float64 `json:"read_mbps"`
	WriteMBps float64 `json:"write_mbps"`
}

func LoadConfig(path string) (*Config, error) {
	exeDir, err := exeDirectory()
	if err == nil {
		path = filepath.Join(exeDir, path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Collector.Interval == 0 {
		cfg.Collector.Interval = 1
	}
	return &cfg, nil
}

func exeDirectory() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}
