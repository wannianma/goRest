package ny

import (
	"encoding/json"
	"fmt"
	"goWeb/server"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
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
	anserAt uint64
	answerA int
	answerB int
	teamA   *Team
	teamB   *Team
}

type Question struct {
	Qid     int      `json:"qid"`
	Title   string   `json:"title"`
	Options []string `json:"options"`
	Answer  int
}

func (info *TeamInfo) setAnswerA(data map[string]int) {
	mutexTeam.Lock()
	defer mutexTeam.Unlock()
	now := uint64(time.Now().Unix())
	if data["qid"] > info.curAnswer && now > info.anserAt {
		info.curAnswer = data["qid"]
		info.answerA = data["aid"]
	}
	getAnswerBroadcast("A").Submit(fmt.Sprintf("%d:%d", data["qid"], data["aid"]))
}

func (info *TeamInfo) setAnswerB(data map[string]int) {
	mutexTeam.Lock()
	defer mutexTeam.Unlock()
	now := uint64(time.Now().Unix())
	if data["qid"] > info.curAnswer && now > info.anserAt {
		info.curAnswer = data["qid"]
		info.answerB = data["aid"]
	}
	getAnswerBroadcast("B").Submit(fmt.Sprintf("%d:%d", data["qid"], data["aid"]))
}

func (info *TeamInfo) setTeamA(data map[string]uint64) {
	mutexTeam.Lock()
	defer mutexTeam.Unlock()
	if info.teamA.distance < info.totalDistance {
		info.teamA.distance = data["distance"]
		info.teamA.power = data["power"]
	}
}

func (info *TeamInfo) setTeamB(data map[string]uint64) {
	mutexTeam.Lock()
	defer mutexTeam.Unlock()
	if info.teamB.distance < info.totalDistance {
		info.teamB.distance = data["distance"]
		info.teamB.power = data["power"]
	}
}

func (info *TeamInfo) getTeamData() map[string]uint64 {
	mutexTeam.RLock()
	defer mutexTeam.RUnlock()
	savedPower := map[string]uint64{
		"timestamp": uint64(time.Now().Unix()),
		"powerA":    info.teamA.power,
		"powerB":    info.teamB.power,
		"distanceA": info.teamA.distance,
		"distanceB": info.teamB.distance}
	return savedPower
}

func (info *TeamInfo) getCurAnswer() int {
	mutexTeam.RLock()
	defer mutexTeam.RUnlock()
	return info.curAnswer
}

func (info *TeamInfo) getStartTime() uint64 {
	mutexTeam.RLock()
	defer mutexTeam.RUnlock()
	return info.anserAt
}

func (info *TeamInfo) setStartTime() {
	mutexTeam.Lock()
	defer mutexTeam.Unlock()
	env := server.Inst()

	info.anserAt = uint64(time.Now().Unix()) + 90
	info.totalDistance = uint64(env.TotalDistance)
}

func (info *TeamInfo) getAnserData() {
	mutexTeam.RLock()
	defer mutexTeam.RUnlock()
}

func loadDataFromFile() []byte {
	wd, _ := os.Getwd()
	b, err := ioutil.ReadFile(path.Join(wd, "questions.json"))
	if err != nil {
		log.Printf("%s", err)
		return []byte("")
	}
	return b
}

var (
	answerChannels = make(map[string]broadcast.Broadcaster)
	mutexTeam      sync.RWMutex
	teamInfo       = TeamInfo{
		totalDistance: 1000,
		curAnswer:     0,
		anserAt:       0,
		answerA:       0,
		answerB:       0,
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

// StreamData sse stream push event
func StreamData(c *gin.Context) {
	listenerA := openListener("A")
	listenerB := openListener("B")
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		closeListener("A", listenerA)
		closeListener("B", listenerB)
		ticker.Stop()
	}()

	c.Stream(func(w io.Writer) bool {
		select {
		case msg := <-listenerA:
			c.SSEvent("answerA", msg)
		case msg := <-listenerB:
			c.SSEvent("answerB", msg)
		case <-ticker.C:
			c.SSEvent("stats", teamInfo.getTeamData())
		}
		return true
	})
}

// GetQuestions 拉取题目列表
func GetQuestions(c *gin.Context) {
	var questionJSON []Question
	if err := json.Unmarshal(questionsStr, &questionJSON); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "load questions error"})
		return
	}
	data := map[string]interface{}{
		"questions": questionJSON,
		"curAnswer": teamInfo.getCurAnswer(),
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": data})
}

func TeamStart(c *gin.Context) {
	teamInfo.setStartTime()
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": ""})
}

func PushData(c *gin.Context) {
	teamName := c.Param("name")
	distance, _ := strconv.ParseUint(c.Query("distance"), 10, 64)
	power, _ := strconv.ParseUint(c.Query("power"), 10, 64)
	speed, _ := strconv.ParseUint(c.Query("speed"), 10, 64)

	teamData := map[string]uint64{
		"distance": distance,
		"power":    power,
		"speed":    speed,
	}

	if teamName == "A" {
		teamInfo.setTeamA(teamData)
	} else {
		teamInfo.setTeamB(teamData)
	}

	state := map[string]interface{}{
		"curAnswer": teamInfo.getCurAnswer(),
		"startAt":   teamInfo.getStartTime(),
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": state})
}

func PushAnswer(c *gin.Context) {
	teamName := c.Param("name")
	qid, _ := strconv.Atoi(c.Query("qid"))
	aid, _ := strconv.Atoi(c.Query("aid"))

	answerData := map[string]int{
		"qid": qid,
		"aid": aid,
	}

	if teamName == "A" {
		teamInfo.setAnswerA(answerData)
	} else {
		teamInfo.setAnswerB(answerData)
	}

	state := map[string]int{
		"curAnswer": teamInfo.getCurAnswer(),
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": state})
}

func ResetTeam(c *gin.Context) {
	mutexTeam.Lock()
	defer mutexTeam.Unlock()
	teamInfo = TeamInfo{
		totalDistance: 1000,
		curAnswer:     0,
		anserAt:       0,
		answerA:       0,
		answerB:       0,
		teamA: &Team{
			power:    0,
			distance: 0,
		},
		teamB: &Team{
			power:    0,
			distance: 0,
		},
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": ""})
}
