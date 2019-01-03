package main

import (
	"fmt"
	"hello/cralwer/models"
	"log"
	"regexp"
	"strconv"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/gocolly/redisstorage"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func getPostId(text string) int {
	reg := regexp.MustCompile(`bid=(\d+)`)
	matchArr := reg.FindStringSubmatch(text)
	pid := 0
	if len(matchArr) > 1 {
		pid, _ = strconv.Atoi(matchArr[1])
	}
	return int(pid)
}

func getDetailId(text string) (int, int) {
	reg := regexp.MustCompile(`id=(\d+)&cid=(\d+)`)
	matchArr := reg.FindStringSubmatch(text)
	pid := 0
	score := 0
	if len(matchArr) > 2 {
		pid, _ = strconv.Atoi(matchArr[1])
		score, _ = strconv.Atoi(matchArr[2])
	}
	return int(pid), int(score)
}

func checkIsHasPost(url string) bool {
	reg := regexp.MustCompile(`www\.cuiweijuxs\.com/0_(\d+)/(\d+)\.html`)
	matchArr := reg.FindStringSubmatch(url)
	pid := 0
	score := 0
	if len(matchArr) > 2 {
		pid, _ = strconv.Atoi(matchArr[1])
		score, _ = strconv.Atoi(matchArr[2])
	} else {
		return false
	}
	detail := models.Post{
		NID:       pid,
		Remote_id: score,
	}
	db.First(&detail, "remote_id = ? and n_id = ?", score, pid)
	if detail.ID > 0 {
		return true
	}
	return false
}

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:Manager@(127.0.0.1:3306)/novel?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
}

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("www.cuiweijuxs.com"),
	)

	// create the redis storage
	storage := &redisstorage.Storage{
		Address:  "127.0.0.1:6379",
		Password: "",
		DB:       0,
		Prefix:   "httpbin_test",
	}

	// add storage to the collector
	err := c.SetStorage(storage)
	if err != nil {
		panic(err)
	}

	// delete previous data from storage
	if err := storage.Clear(); err != nil {
		log.Fatal(err)
	}

	// close redis client
	defer storage.Client.Close()

	q, _ := queue.New(
		4, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		if !checkIsHasPost(e.Request.AbsoluteURL(link)) {
			q.AddURL(e.Request.AbsoluteURL(link))
		}
	})

	// list 页
	c.OnHTML("#maininfo", func(e *colly.HTMLElement) {
		title := e.ChildText("#info h1")
		anthor := e.ChildText("#info > p:nth-child(2)")
		intro := e.ChildText("#intro p")
		text := e.ChildAttr("#info > p:nth-child(3) > a:nth-child(1)", "onclick")
		pid := getPostId(text)
		// title = Decode(title)
		// intro = ConvertToByte(intro, "gbk", "utf8")
		// anthor = ConvertToByte(anthor, "GB18030", "utf8")
		post := models.Novel{
			Remote_id: pid,
			Title:     title,
			Anthor:    anthor,
			Intro:     intro,
			Nums:      0,
		}

		db.First(&post, "remote_id = ?", pid)
		if post.ID > 0 {
			log.Printf("ID %d remote_id %d is existed\n", post.ID, pid)
		} else {
			db.Create(&post)
		}
	})

	// detail 页
	c.OnHTML("div.content_read .box_con", func(e *colly.HTMLElement) {
		title := e.ChildText(".bookname h1")
		content := e.ChildText("#content")
		text := e.ChildAttr("div.bookname > div.bottem1 > a:nth-child(5)", "onclick")
		pid, score := getDetailId(text)

		// title = ConvertToByte(title, "GB18030", "utf8")
		// content = ConvertToByte(content, "GB18030", "utf8")

		detail := models.Post{
			NID:       pid,
			Remote_id: score,
			Title:     title,
			Info:      content,
		}
		db.First(&detail, "remote_id = ? and n_id = ?", score, pid)
		if detail.ID > 0 {
			log.Printf("Post ID %d n_id %d remote_id %d is existed\n", detail.ID, pid, score)
		} else {
			db.Create(&detail)
			// db.Update()
		}

		fmt.Printf("$$$ %d - %d - %s \n", pid, score, title)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	q.AddURL("https://www.cuiweijuxs.com/0_69/")
	q.AddURL("https://www.cuiweijuxs.com/dushixiaoshuo/")
	q.AddURL("https://www.cuiweijuxs.com/paihangbang/")
	q.AddURL("https://www.cuiweijuxs.com")
	q.AddURL("https://www.cuiweijuxs.com/wanben/1_1")
	q.Run(c)
}
