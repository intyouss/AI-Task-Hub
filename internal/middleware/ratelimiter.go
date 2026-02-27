package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/intyouss/AI-Task-Hub/config"
	"github.com/intyouss/AI-Task-Hub/pkg"
	"golang.org/x/time/rate"
)

var bucket sync.Map

func RateLimiter(limit ...int) gin.HandlerFunc {
	li := config.Cfg.RateLimiter.Limit
	timer := config.Cfg.RateLimiter.EveryTime
	if len(limit) > 0 {
		li = limit[0]
	}

	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-KEY")

		limiter := rate.NewLimiter(rate.Every(timer), li)
		value, _ := bucket.LoadOrStore(apiKey, limiter)
		if !value.(*rate.Limiter).Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, errors.RateLimitExceededError)
			return
		}
		ctx.Next()
	}
}
