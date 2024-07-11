package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Check the logged-in user
		username, loggedIn := c.Get("username")
		if loggedIn {
			username = username.(string)
		} else {
			username = "_"
		}

		// Process the request
		c.Next()

		// Log the request
		log.l.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"statusCode": c.Writer.Status(),
			"latency":    time.Since(start).String(),
			"clientIP":   c.ClientIP(),
			"user":       username,
		}).Info("Request handled")
	}
}
