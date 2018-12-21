package ny

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	broadcast "github.com/dustin/go-broadcast"
	"github.com/gin-gonic/gin"
)

type Team struct {
	power    uint64
	distance uint64
}

type TeamInfo struct {
	totalDistance uint64
	curAnswer     int
	// 如何确定开始出题时间
	beginAt int
	answerA string
	answerB string
	teamA   *Team
	teamB   *Team
}

type RegisterInput struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nick     string `json:"nick" binding:"required"`
	Mobile   string `json:"mobile" binding:"required"`
	Email    string `json:"email"`
}
type Question struct {
	Qid     int      `json:"qid"`
	Title   string   `json:"title"`
	Options []string `json:"options"`
	Answer  int
}

/*
{
	"no": int,
	"answer_at": 0,
	"answer": "A"
}
*/
func (info *TeamInfo) setAnswer(data map[string]interface{}) {
	mutexTeam.Lock()

	mutexTeam.Unlock()
}

/*
{
	power: int,
	distance: float32
}
*/
func (info *TeamInfo) setTeam(data map[string]interface{}) {
	mutexTeam.Lock()

	mutexTeam.Unlock()
}

func (info *TeamInfo) setBeginAt(beginAt int) {
	mutexTeam.Lock()

	mutexTeam.Unlock()
}

func (info *TeamInfo) getTeamData() map[string]uint64 {
	mutexTeam.Lock()
	defer mutexTeam.Unlock()
	rand.Seed(time.Now().UnixNano())
	if info.teamA.distance < info.totalDistance && info.teamB.distance < info.totalDistance {
		info.teamA.distance += uint64(rand.Intn(20))
		info.teamB.distance += uint64(rand.Intn(20))
	}
	savedPower := map[string]uint64{
		"timestamp": uint64(time.Now().Unix()),
		"powerA":    uint64(rand.Intn(240)),
		"powerB":    uint64(rand.Intn(240)),
		"distanceA": info.teamA.distance,
		"distanceB": info.teamB.distance}
	return savedPower
}

func loadDataFromFile() []byte {
	b, err := ioutil.ReadFile("./questions.json")
	if err != nil {
		log.Printf("%s", err)
		return []byte("")
	}
	return b
}

func (info *TeamInfo) getAnserData() {
	mutexTeam.RLock()
	defer mutexTeam.RUnlock()

}

var (
	answerChannels = make(map[string]broadcast.Broadcaster)
	mutexTeam      sync.RWMutex
	teamInfo       = TeamInfo{
		totalDistance: 1000,
		curAnswer:     1,
		beginAt:       0,
		answerA:       "",
		answerB:       "",
		teamA: &Team{
			power:    0,
			distance: 0,
		},
		teamB: &Team{
			power:    0,
			distance: 0,
		},
	}
	questionsStr = loadDataFromFile()
)

func openListener(roomid string) chan interface{} {
	listener := make(chan interface{})
	getAnswerBroadcast(roomid).Register(listener)
	return listener
}

func closeListener(roomid string, listener chan interface{}) {
	getAnswerBroadcast(roomid).Unregister(listener)
	close(listener)
}

func getAnswerBroadcast(roomid string) broadcast.Broadcaster {
	b, ok := answerChannels[roomid]
	if !ok {
		b = broadcast.NewBroadcaster(10)
		answerChannels[roomid] = b
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
			c.SSEvent("stats", teamInfo.getTeamData())
		}
		return true
	})
}

func GetQuestions(c *gin.Context) {
	var questionJSON []Question
	if err := json.Unmarshal(questionsStr, &questionJSON); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "load questions error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "", "data": questionJSON})
}
