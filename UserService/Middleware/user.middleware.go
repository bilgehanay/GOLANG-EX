package Middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

const (
	maxRequest = 2
	perMinute  = 1 * time.Minute
)

var (
	ipRequestCounts = make(map[string]int)
	mutex           = &sync.Mutex{}
)

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		mutex.Lock()
		defer mutex.Unlock()

		count := ipRequestCounts[ip]
		if count >= maxRequest {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "too many requests"})
			return
		}

		ipRequestCounts[ip] = count + 1
		time.AfterFunc(perMinute, func() {
			mutex.Lock()
			defer mutex.Unlock()

			ipRequestCounts[ip] -= 1
		})
		c.Next()
	}
}
