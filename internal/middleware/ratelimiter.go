package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/intyouss/AI-Task-Hub/pkg"
	"golang.org/x/time/rate"
)

var bucket sync.Map

const LIMIT = 10

func RateLimiter(limit ...int) gin.HandlerFunc {
	li := LIMIT
	if len(limit) > 0 {
		li = limit[0]
	}

	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-KEY")

		limiter := rate.NewLimiter(rate.Every(time.Second), li)
		value, _ := bucket.LoadOrStore(apiKey, limiter)
		if !value.(*rate.Limiter).Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, errors.RateLimitExceededError)
			return
		}
		ctx.Next()
	}
}
