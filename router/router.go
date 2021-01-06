package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/luoxiaojun1992/go-skeleton/middlewares/recovery"
	"net/http"
)

func Register() *gin.Engine {
	r := gin.New()

	// Middleware
	r.Use(recovery.CustomRegister(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			errMsg := fmt.Sprintf("error: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": errMsg,
				"msg":     errMsg,
				"code":    1,
			})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World!")
	})

	return r
}
