package ny

import (
	"io"
	"math/rand"
	"time"

	broadcast "github.com/dustin/go-broadcast"
	"github.com/gin-gonic/gin"
	"github.com/manucorporat/stats"
)

var (
	ips          = stats.New()
	messages     = stats.New()
	roomChannels = make(map[string]broadcast.Broadcaster)
)

func openListener(roomid string) chan interface{} {
	listener := make(chan interface{})
	room(roomid).Register(listener)
	return listener
}

func closeListener(roomid string, listener chan interface{}) {
	room(roomid).Unregister(listener)
	close(listener)
}

func room(roomid string) broadcast.Broadcaster {
	b, ok := roomChannels[roomid]
	if !ok {
		b = broadcast.NewBroadcaster(10)
		roomChannels[roomid] = b
	}
	return b
}

func StreamData(c *gin.Context) {
	roomid := "bb"
	listener := openListener(roomid)
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		closeListener(roomid, listener)
		ticker.Stop()
	}()

	c.Stream(func(w io.Writer) bool {
		select {
		case msg := <-listener:
			c.SSEvent("message", msg)
		case <-ticker.C:
			c.SSEvent("stats", getRanPower())
		}
		return true
	})
}

func getRanPower() map[string]uint64 {
	rand.Seed(time.Now().UnixNano())
	savedPower := map[string]uint64{
		"timestamp": uint64(time.Now().Unix()),
		"powerA":    uint64(rand.Intn(240)),
		"powerB":    uint64(rand.Intn(240))}
	return savedPower
}
