package collector

import (
	"net"
	"strconv"
	"sync"
	"syscall"
	"unsafe"

	"github.com/shirou/gopsutil/v3/process"
)

var staticOnce sync.Once
var cachedIP, cachedDisplay string

type devMode struct {
	dmDeviceName       [32]uint16
	dmSpecVersion      uint16
	dmDriverVersion    uint16
	dmSize             uint16
	dmDriverExtra      uint16
	dmFields           uint32
	dmPositionX        int32
	dmPositionY        int32
	dmDisplayOrientation uint32
	dmDisplayFixedOutput uint32
	dmColor            int16
	dmDuplex           int16
	dmYResolution      int16
	dmTTOption         int16
	dmCollate          int16
	dmFormName         [32]uint16
	dmLogPixels        uint16
	dmBitsPerPel       uint32
	dmPelsWidth        uint32
	dmPelsHeight       uint32
	dmDisplayFlags     uint32
	dmDisplayFrequency uint32
	dmICM              uint32
}

func collectStatic() {
	cachedIP = localIP()
	cachedDisplay = displayInfo()
}

func localIP() string {
	if ip := preferredOutboundIP(); ip != "" {
		return ip
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return ""
}

func preferredOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return ""
	}
	defer conn.Close()
	if udpAddr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
		return udpAddr.IP.String()
	}
	return ""
}

func displayInfo() string {
	dll := syscall.NewLazyDLL("user32.dll")
	procSM := dll.NewProc("GetSystemMetrics")
	w, _, _ := procSM.Call(0)
	h, _, _ := procSM.Call(1)
	if w == 0 || h == 0 {
		return ""
	}

	res := strconv.Itoa(int(w)) + "x" + strconv.Itoa(int(h))

	procEDS := dll.NewProc("EnumDisplaySettingsW")
	var dm devMode
	dm.dmSize = uint16(unsafe.Sizeof(dm))
	ret, _, _ := procEDS.Call(0, uintptr(0xFFFFFFFF), uintptr(unsafe.Pointer(&dm)))
	if ret != 0 && dm.dmDisplayFrequency > 0 {
		res += "@" + strconv.Itoa(int(dm.dmDisplayFrequency)) + "Hz"
	}

	return res
}

func CollectSysInfo() (ipAddr, displayInfo string, procCount, threadCount int) {
	staticOnce.Do(collectStatic)
	if pids, err := process.Pids(); err == nil {
		procCount = len(pids)
	}
	type procInfo struct {
		ThreadCount uint32
	}
	var procs []procInfo
	if err := wmiQuery("SELECT ThreadCount FROM Win32_Process", &procs); err == nil {
		for _, p := range procs {
			threadCount += int(p.ThreadCount)
		}
	}
	return cachedIP, cachedDisplay, procCount, threadCount
}
