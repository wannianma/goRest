package novel

import (
	"encoding/json"
	"fmt"
	"goWeb/server"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Novel struct {
	ID     uint
	Title  string
	Anthor string
	Intro  string
}

type Detail struct {
	ID    uint
	Score uint
	Title string
	Info  string
}

func NovelList(c *gin.Context) {
	page := c.DefaultQuery("page", "0")
	env := server.Inst()
	client := env.RClient

	start, _ := strconv.ParseInt(page, 10, 64)
	rs, err := client.LRange("list:post", start*100, (start+1)*100).Result()
	if err != nil {
		log.Println(err)
	}
	var result []map[string]string
	for _, val := range rs {
		strs := strings.Split(val, "-")
		rs2, _ := client.Keys(strs[0] + ":post:*").Result()
		length := fmt.Sprintf("%d", len(rs2))

		result = append(result, map[string]string{
			"id":     strs[0],
			"title":  strs[1],
			"length": length,
		})
	}

	c.HTML(http.StatusOK, "novels/list", gin.H{
		"result": result,
	})
}

func NovelDetail(c *gin.Context) {
	pid := c.Param("pid")
	fmt.Println(pid)
	env := server.Inst()
	client := env.RClient

	rs, err := client.Get(pid + ":post:0").Result()
	if err != nil {
		log.Println(err)
	}
	var novel Novel
	if err = json.Unmarshal([]byte(rs), &novel); err != nil {
		log.Println(err)
	}

	rsArr, err := client.Keys(pid + ":post:*").Result()
	if err != nil {
		log.Println(err)
	}
	sort.Sort(sort.StringSlice(rsArr))

	var detailArr [][]string
	threeArr := []string{}
	for idx, val := range rsArr {
		if result := strings.Index(val, ":0"); result > 0 {
			continue
		}
		rsDetail, _ := client.Get(val).Result()
		var detail Detail
		json.Unmarshal([]byte(rsDetail), &detail)
		threeArr = append(threeArr, detail.Title)
		if idx%3 == 0 {
			detailArr = append(detailArr, threeArr)
			threeArr = []string{}
		}
	}

	c.HTML(http.StatusOK, "novels/detail", gin.H{
		"novel":   novel,
		"details": detailArr,
	})
}
