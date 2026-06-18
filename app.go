package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"monitor-screen/collector"
	"monitor-screen/model"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/yusufpapurcu/wmi"
)

//go:embed icon.ico
var iconICO []byte

var (
	user32DLL                    = syscall.NewLazyDLL("user32.dll")
	procGetCursorPos             = user32DLL.NewProc("GetCursorPos")
	procGetWindowRect            = user32DLL.NewProc("GetWindowRect")
	procSetWindowPos             = user32DLL.NewProc("SetWindowPos")
	procFindWindowW              = user32DLL.NewProc("FindWindowW")
	procSendMessageW             = user32DLL.NewProc("SendMessageW")
	procCreateIconFromResourceEx = user32DLL.NewProc("CreateIconFromResourceEx")
	procEnumDisplayMonitors      = user32DLL.NewProc("EnumDisplayMonitors")
	procGetMonitorInfoW          = user32DLL.NewProc("GetMonitorInfoW")
)

type winPoint struct{ X, Y int32 }
type winRect struct{ Left, Top, Right, Bottom int32 }

type monitorInfoEx struct {
	cbSize    uint32
	rcMonitor winRect
	rcWork    winRect
	dwFlags   uint32
}

type displayMonitor struct {
	Rect      winRect
	IsPrimary bool
}

func getMonitors() []displayMonitor {
	var monitors []displayMonitor

	callback := syscall.NewCallback(func(hMonitor, hdc, lprcMonitor, dwData uintptr) uintptr {
		var mi monitorInfoEx
		mi.cbSize = uint32(unsafe.Sizeof(mi))
		ret, _, _ := procGetMonitorInfoW.Call(hMonitor, uintptr(unsafe.Pointer(&mi)))
		if ret == 0 {
			return 0
		}
		ptr := (*[]displayMonitor)(unsafe.Pointer(dwData))
		*ptr = append(*ptr, displayMonitor{
			Rect:      mi.rcMonitor,
			IsPrimary: mi.dwFlags&1 != 0,
		})
		return 1
	})

	procEnumDisplayMonitors.Call(0, 0, callback, uintptr(unsafe.Pointer(&monitors)))
	return monitors
}

type App struct {
	ctx      context.Context
	cfg      *model.Config

	mu          sync.RWMutex
	dashboard   *model.DashboardData
	quitting    bool
	cancelColl  context.CancelFunc
	cpuWarmedUp bool
	fullscreen  bool
	hwnd        uintptr
	dragging    bool
	dragOffX    int
	dragOffY    int
}

func (a *App) getHWND() uintptr {
	if a.hwnd != 0 {
		return a.hwnd
	}
	title, _ := syscall.UTF16PtrFromString("Monitor Screen")
	h, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(title)))
	if h != 0 {
		a.hwnd = h
	}
	return h
}

func (a *App) setWindowIcon() {
	hwnd := a.getHWND()
	if hwnd == 0 {
		time.Sleep(500 * time.Millisecond)
		hwnd = a.getHWND()
	}
	if hwnd == 0 {
		return
	}

	icoData := iconICO
	res := findICOIconResource(icoData, 32)
	if res == nil {
		res = findICOIconResource(icoData, 16)
	}
	if res == nil {
		return
	}

	hIcon, _, _ := procCreateIconFromResourceEx.Call(
		uintptr(unsafe.Pointer(&res.data[0])),
		uintptr(len(res.data)),
		1,
		0x00030000,
		uintptr(res.width),
		uintptr(res.height),
		0,
	)
	if hIcon != 0 {
		WM_SETICON := uintptr(0x0080)
		ICON_SMALL := uintptr(0)
		ICON_BIG := uintptr(1)
		procSendMessageW.Call(hwnd, WM_SETICON, ICON_SMALL, hIcon)
		procSendMessageW.Call(hwnd, WM_SETICON, ICON_BIG, hIcon)
	}
}

type icoIconResource struct {
	width  int
	height int
	data   []byte
}

func findICOIconResource(icoData []byte, targetSize int) *icoIconResource {
	if len(icoData) < 6 {
		return nil
	}
	count := int(icoData[4]) | (int(icoData[5]) << 8)
	offset := 6
	for i := 0; i < count && offset+16 <= len(icoData); i++ {
		w := int(icoData[offset])
		if w == 0 {
			w = 256
		}
		imgSize := int(icoData[offset+8]) | (int(icoData[offset+9]) << 8) | (int(icoData[offset+10]) << 16) | (int(icoData[offset+11]) << 24)
		imgOffset := int(icoData[offset+12]) | (int(icoData[offset+13]) << 8) | (int(icoData[offset+14]) << 16) | (int(icoData[offset+15]) << 24)
		if w == targetSize && imgOffset+imgSize <= len(icoData) {
			return &icoIconResource{
				width:  w,
				height: w,
				data:  icoData[imgOffset : imgOffset+imgSize],
			}
		}
		offset += 16
	}
	return nil
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	initLogging()

	cfg, err := model.LoadConfig("config.yaml")
	if err != nil {
		log.Printf("failed to load config: %v, using defaults", err)
		cfg = defaultConfig()
	}
	a.cfg = cfg

	monitors := getMonitors()
	if len(monitors) >= 2 {
		for _, m := range monitors {
			if !m.IsPrimary {
				runtime.WindowSetPosition(ctx, int(m.Rect.Left), int(m.Rect.Top))
				break
			}
		}
	}
	runtime.WindowFullscreen(ctx)
	a.fullscreen = true

	if a.cfg.LHM.Enabled == "true" {
		go func() {
			collector.EnsureDriverLoaded(a.cfg.LHM.LHMExe)
			bridgeExe := a.cfg.LHM.BridgeExe
			if bridgeExe == "" {
				bridgeExe = "lhm\\sensor_bridge.exe"
			}
			collector.StartBridge(bridgeExe)
		}()
	}

	go a.setWindowIcon()

	var collCtx context.Context
	collCtx, a.cancelColl = context.WithCancel(context.Background())
	go a.runCollectors(collCtx)

	go a.runSystray()
}

func (a *App) shutdown(ctx context.Context) {
	a.quitting = true
	collector.StopBridge()
	if a.cancelColl != nil {
		a.cancelColl()
	}
}

func (a *App) beforeClose(ctx context.Context) bool {
	collector.StopBridge()
	return false
}

func (a *App) ShowWindow() {
	if a.ctx != nil {
		runtime.WindowShow(a.ctx)
		runtime.WindowUnminimise(a.ctx)
	}
}

func (a *App) QuitApp() {
	a.mu.Lock()
	a.quitting = true
	a.mu.Unlock()
	collector.StopBridge()
	if a.cancelColl != nil {
		a.cancelColl()
	}
	if a.ctx != nil {
		runtime.Quit(a.ctx)
	}
}

func (a *App) GetDashboard() *model.DashboardData {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.dashboard
}

func (a *App) Shutdown() {
	a.QuitApp()
}

func (a *App) GetIsFullscreen() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.fullscreen
}

func (a *App) ToggleFullscreen() {
	if a.ctx == nil {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.fullscreen {
		runtime.WindowUnfullscreen(a.ctx)
		runtime.WindowSetSize(a.ctx, 1280, 820)
		runtime.WindowCenter(a.ctx)
		a.fullscreen = false
	} else {
		runtime.WindowFullscreen(a.ctx)
		a.fullscreen = true
	}
}

func (a *App) DragStart(_, _ int) {
	var pt winPoint
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))

	var rect winRect
	hwnd := a.getHWND()
	procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&rect)))

	a.mu.Lock()
	a.dragOffX = int(pt.X) - int(rect.Left)
	a.dragOffY = int(pt.Y) - int(rect.Top)
	a.dragging = true
	a.mu.Unlock()
}

func (a *App) DragMove(_, _ int) {
	a.mu.RLock()
	if !a.dragging {
		a.mu.RUnlock()
		return
	}
	offX, offY := a.dragOffX, a.dragOffY
	a.mu.RUnlock()

	var pt winPoint
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))

	newX := int(pt.X) - offX
	newY := int(pt.Y) - offY

	hwnd := a.getHWND()
	procSetWindowPos.Call(hwnd, 0, uintptr(newX), uintptr(newY), 0, 0, 0x0001|0x0004)
}

func (a *App) DragEnd() {
	a.mu.Lock()
	a.dragging = false
	a.mu.Unlock()
}

func (a *App) buildDashboard() *model.DashboardData {
	d := &model.DashboardData{}

	uptime := getRealUptime()
	d.System.Uptime = formatUptime(uptime)

	if !a.cpuWarmedUp {
		collector.CollectCPU(nil)
		a.cpuWarmedUp = true
	}
	var cpuUsage float64
	d.CPU = collector.CollectCPU(&cpuUsage)

	d.GPU = collector.CollectGPU()
	d.Memory = collector.CollectMemory()
	d.Storage = collector.CollectStorage()

	if a.cfg.LHM.Enabled == "true" {
		bridgeExe := a.cfg.LHM.BridgeExe
		if bridgeExe == "" {
			bridgeExe = "lhm\\sensor_bridge.exe"
		}
		bridge := collector.QueryBridgeIfAlive(bridgeExe)
		if bridge != nil {
			fans := collector.EnrichFromLHM(bridge, &d.CPU, &d.GPU, d.Storage)
			if len(fans) > 0 {
				d.Fans = fans
			}
			d.TotalPower = collector.TotalPowerFromLHM(bridge)
		}
	}

	return d
}

func (a *App) runCollectors(ctx context.Context) {
	go func() {
		var dummy float64
		collector.CollectCPU(&dummy)
		a.mu.Lock()
		a.cpuWarmedUp = true
		a.mu.Unlock()
	}()

	time.Sleep(100 * time.Millisecond)

	a.emitDashboard(ctx)

	ticker := time.NewTicker(time.Duration(a.cfg.Collector.Interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.emitDashboard(ctx)
		}
	}
}

func (a *App) emitDashboard(ctx context.Context) {
	db := a.buildDashboard()

	a.mu.Lock()
	a.dashboard = db
	a.mu.Unlock()

	data, err := json.Marshal(db)
	if err == nil && a.ctx != nil {
		runtime.EventsEmit(a.ctx, "dashboard-update", string(data))
	}
}

func formatUptime(seconds uint64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60

	if days > 0 {
		return fmt.Sprintf("%d天 %d小时 %d分钟", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%d小时 %d分钟", hours, minutes)
	}
	return fmt.Sprintf("%d分钟", minutes)
}

func getRealUptime() uint64 {
	type osInfo struct {
		LastBootUpTime string
	}
	var infos []osInfo
	err := wmiQuery("SELECT LastBootUpTime FROM Win32_OperatingSystem", &infos)
	if err == nil && len(infos) > 0 && len(infos[0].LastBootUpTime) >= 14 {
		s := infos[0].LastBootUpTime
		year, _ := strconv.Atoi(s[0:4])
		month, _ := strconv.Atoi(s[4:6])
		day, _ := strconv.Atoi(s[6:8])
		hour, _ := strconv.Atoi(s[8:10])
		min, _ := strconv.Atoi(s[10:12])
		sec, _ := strconv.Atoi(s[12:14])
		bootTime := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.Local)
		elapsed := time.Since(bootTime)
		if elapsed > 0 {
			return uint64(elapsed.Seconds())
		}
	}

	uptime, err := host.Uptime()
	if err == nil {
		return uptime
	}
	return 0
}

func wmiQuery(query string, dst interface{}) error {
	return wmi.Query(query, dst)
}

func initLogging() {
	exeDir := "."
	if exe, err := os.Executable(); err == nil {
		exeDir = filepath.Dir(exe)
	}
	logPath := filepath.Join(exeDir, "monitor-screen.log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		log.SetOutput(f)
		log.Printf("=== monitor-screen started ===")
	}
}

func defaultConfig() *model.Config {
	return &model.Config{
		Collector: model.CollectorConfig{
			Interval: 1,
		},
		LHM: model.LHMConfig{
			Enabled:   "false",
			BridgeExe: "lhm\\sensor_bridge.exe",
			LHMExe:    "lhm\\LibreHardwareMonitor.exe",
		},
	}
}
