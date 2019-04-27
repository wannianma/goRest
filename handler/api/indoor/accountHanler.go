package account

import (
	"goWeb/models"
	"goWeb/server"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var json LoginInput
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if json.Account != "manu" || json.Password != "123" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	name := c.DefaultQuery("name", "Guest")
	action := c.DefaultQuery("action", "Guest")
	message := name + " is " + action
	c.String(http.StatusOK, message)
}

type RegisterInput struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nick     string `json:"nick" binding:"required"`
	Mobile   string `json:"mobile" binding:"required"`
	Email    string `json:"email"`
}

func Register(c *gin.Context) {
	var json RegisterInput
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	salt := string(_genSalt(4))
	passwd, _ := bcrypt.GenerateFromPassword([]byte(json.Password+salt), bcrypt.MinCost)
	log.Println(json.Nick)
	env := server.Inst()
	user := models.Users{
		UserName: json.Nick,
		Birthday: time.Now(),
		Mobile:   json.Mobile,
		Email:    json.Email,
		Password: string(passwd),
		Salt:     salt}

	if err := env.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
}

func _genSalt(num int) []byte {
	return []byte("abce")
}
