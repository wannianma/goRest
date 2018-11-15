package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	name := c.DefaultQuery("name", "Guest")
	action := c.DefaultQuery("action", "Guest")
	message := name + " is " + action
	c.String(http.StatusOK, message)
}
