package workers

import (
	"runtime"
	"sync"
	"time"
)

var (
	mutexStats sync.RWMutex
	savedStats map[string]uint64
)

func StatsWorker() {
	c := time.Tick(1 * time.Second)
	var lastMallocs uint64
	var lastFrees uint64
	for range c {
		var stats runtime.MemStats
		runtime.ReadMemStats(&stats)

		mutexStats.Lock()
		savedStats = map[string]uint64{
			"timestamp":  uint64(time.Now().Unix()),
			"HeapInuse":  stats.HeapInuse,
			"StackInuse": stats.StackInuse,
			"Mallocs":    stats.Mallocs - lastMallocs,
			"Frees":      stats.Frees - lastFrees,
			// "Connected":  connectedUsers(),
		}
		lastMallocs = stats.Mallocs
		lastFrees = stats.Frees
		mutexStats.Unlock()
	}
}

// func connectedUsers() uint64 {
// 	connected := users.Get("connected") - users.Get("disconnected")
// 	if connected < 0 {
// 		return 0
// 	}
// 	return uint64(connected)
// }

// Stats returns savedStats data.
func Stats() map[string]uint64 {
	mutexStats.RLock()
	defer mutexStats.RUnlock()
	return savedStats
}
