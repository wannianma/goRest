package server

import (
	"goWeb/models"
	"goWeb/token"
	"log"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var once sync.Once
var instance *Env

type appConfig struct {
	Host                  string
	Port                  string
	ListenAddr            string `toml:"-"`
	DBUrl                 string
	RSAPriKey             string
	RSAPubKey             string
	MaxAccessTokenMinute  uint
	MaxRefreshTokenMinute uint
	AllowCORS             bool
}

type tomlConfig struct {
	App appConfig
}

type Env struct {
	appConfig
	Gin          *gin.Engine
	DB           *gorm.DB `toml:"-"`
	TokenManager *token.TokenManager
}

func (e *Env) Drop() {
	if e.DB != nil {
		e.DB.Close()
		e.DB = nil
	}
}

func (e *Env) _db_create() {
	if e.DB != nil {
		e.DB.AutoMigrate(&models.Users{})
	}
}

func _init(fpath string) *Env {
	var conf tomlConfig
	if _, err := toml.DecodeFile(fpath, &conf); err != nil {
		log.Fatalf("Invalid toml file: %s, decode with error: %v", fpath, err)
	}

	env := Env{}
	env.appConfig = conf.App
	env.ListenAddr = env.Host + ":" + env.Port

	db, err := gorm.Open("mysql", env.DBUrl)
	if err != nil {
		log.Fatalf("Error Open Database '%v'", err)
	}
	env.DB = db
	env._db_create()
	env.TokenManager, err = token.New(
		env.RSAPriKey,
		env.RSAPubKey,
		env.MaxAccessTokenMinute*60,
		env.MaxRefreshTokenMinute*60,
	)
	if err != nil {
		log.Fatalf("Error Create TokenManager '%v'", err)
	}
	env.Gin = gin.Default()
	// gin.SetMode(gin.ReleaseMode)

	if env.AllowCORS {
		log.Println("In DEBUG MODE, CORS is ALLOWED")
		env.Gin.Use(corsMiddleware())
	}

	return &env
}

func Inst() *Env {
	once.Do(func() {
		instance = _init("./config.toml")
	})
	return instance
}

// a helper middleware used to by-pass CORS
// should only be used in Development
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
