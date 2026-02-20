package middlewares

import (
	"log"
	"net/http"
	"runtime/debug"

	"recipe-api/src/utils"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("🔥 PANIC RECOVERED: %v\n", err)
				log.Printf("Stack Trace:\n%s\n", debug.Stack())

				utils.ErrorResponse(c, http.StatusInternalServerError,
					"Internal server error — something went wrong on our end")
				c.Abort()
			}
		}()
		c.Next()
	}
}
