package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// t := time.Now()

		// before request

		c.Next()

		// after request
		// latency := time.Since(t)
		// log.Print(latency)

		log.Println("A function was ran")
		// access the status we are sending
		// status := c.Writer.Status()
		// log.Println(status)

	}
}
