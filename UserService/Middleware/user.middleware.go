package Middleware

import (
	service "deneme.com/bng-go/Service"
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

func VerifyJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if err := c.ShouldBindHeader(&token); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		_, err := service.ParseAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}
		c.Next()
	}
}
