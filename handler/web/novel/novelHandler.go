package novel

import (
	"fmt"
	"goWeb/models"
	"goWeb/server"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Detail struct {
	ID    uint
	Score uint
	Title string
	Info  string
}

var pageSize = 20

func NovelList(c *gin.Context) {
	page := c.DefaultQuery("page", "0")
	env := server.Inst()
	pageInt, _ := strconv.Atoi(page)

	var novels []models.Novel
	env.DB.Offset(pageInt * pageSize).Limit(pageSize).Find(&novels)

	var nidArr []int

	for _, val := range novels {
		nidArr = append(nidArr, val.Remote_id)
	}

	rows, err := env.DB.Raw("SELECT n_id, count(id) as total from posts WHERE n_id IN (?) group by n_id", nidArr).Rows()
	if err != nil {
		log.Panicln(err)
	}
	nidMap := make(map[int]int)
	for rows.Next() {
		var nid int
		var total int
		rows.Scan(&nid, &total)
		nidMap[nid] = total
	}

	var result []map[string]interface{}
	for _, val1 := range novels {
		result = append(result, map[string]interface{}{
			"id":     val1.Remote_id,
			"title":  val1.Title,
			"anchor": val1.Anthor,
			"length": nidMap[val1.Remote_id],
		})
	}

	pagePre := 0
	if pageInt > 0 {
		pagePre = pageInt - 1
	}

	c.HTML(http.StatusOK, "novels/list", gin.H{
		"result": result,
		"pre":    pagePre,
		"next":   (pageInt + 1),
	})
}

func NovelDetail(c *gin.Context) {
	pid := c.Param("pid")
	env := server.Inst()
	pidInt, _ := strconv.Atoi(pid)

	var novel models.Novel
	env.DB.Where("remote_id = ?", pidInt).First(&novel)

	var posts []models.Post

	env.DB.Select("id, title").Where("n_id = ?", pidInt).Order("remote_id asc").Find(&posts)

	var detailArr [][]models.Post
	var threeArr []models.Post
	arrLength := len(posts)
	for idx, val := range posts {
		threeArr = append(threeArr, val)
		if idx%3 == 2 {
			detailArr = append(detailArr, threeArr)
			threeArr = []models.Post{}
		} else if idx == (arrLength - 1) {
			detailArr = append(detailArr, threeArr)
		}
	}

	c.HTML(http.StatusOK, "novels/detail", gin.H{
		"novel":   novel,
		"details": detailArr,
	})
}

func PostDetail(c *gin.Context) {
	id := c.Param("id")
	env := server.Inst()
	idInt, _ := strconv.Atoi(id)

	var post models.Post
	var nodel models.Novel
	var prePost models.Post
	var nextPost models.Post

	env.DB.Where("id = ?", idInt).First(&post)
	env.DB.Where("remote_id = ?", post.NID).First(&nodel)
	env.DB.Limit(1).Where("n_id = ? AND remote_id < ?", post.NID, post.Remote_id).Order("remote_id desc").Select("id").First(&prePost)
	env.DB.Limit(1).Where("n_id = ? AND remote_id > ?", post.NID, post.Remote_id).Order("remote_id asc").Select("id").First(&nextPost)

	preURL := "#"
	nextURL := "#"
	if prePost.ID != 0 {
		preURL = fmt.Sprintf("/novel/post/detail/%d", prePost.ID)
	}

	if nextPost.ID != 0 {
		nextURL = fmt.Sprintf("/novel/post/detail/%d", nextPost.ID)
	}

	c.HTML(http.StatusOK, "novels/postDetail", gin.H{
		"post":    post,
		"novel":   nodel,
		"preURL":  preURL,
		"nextURL": nextURL,
	})
}
