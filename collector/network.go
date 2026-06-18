package collector

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/net"
)

var (
	netMu       sync.Mutex
	lastBytes   = make(map[string]recvSent)
	lastNetTime time.Time
	netDown     float64
	netUp       float64
)

type recvSent struct{ recv, sent uint64 }

func CollectNetworkRate() (downMbps, upMbps float64) {
	netMu.Lock()
	defer netMu.Unlock()

	counters, err := net.IOCounters(true)
	if err != nil {
		return netDown, netUp
	}

	now := time.Now()
	var totalRecv, totalSent uint64
	current := make(map[string]recvSent)

	for _, c := range counters {
		totalRecv += c.BytesRecv
		totalSent += c.BytesSent
		current[c.Name] = recvSent{c.BytesRecv, c.BytesSent}
	}

	if len(lastBytes) > 0 && !lastNetTime.IsZero() {
		elapsed := now.Sub(lastNetTime).Seconds()
		if elapsed >= 0.5 {
			var diffRecv, diffSent uint64
			for name, cur := range current {
				if prev, ok := lastBytes[name]; ok {
					diffRecv += cur.recv - prev.recv
					diffSent += cur.sent - prev.sent
				}
			}
			netDown = float64(diffRecv) * 8 / elapsed / 1e6
			netUp = float64(diffSent) * 8 / elapsed / 1e6
		}
	}

	lastBytes = current
	lastNetTime = now
	return netDown, netUp
}

func CollectConnectionCount() int {
	conns, err := net.Connections("tcp")
	if err != nil {
		return 0
	}
	count := 0
	for _, c := range conns {
		if c.Status == "ESTABLISHED" {
			count++
		}
	}
	return count
}
