package middleware

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

var bucket sync.Map

const (
	LIMIT = 1
	BURST
)

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		if v, ok := bucket.Load(apiKey); !ok {
			bucket.Store(apiKey, newLimiter())
		} else if !v.(*rate.Limiter).Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func newLimiter() *rate.Limiter {
	return rate.NewLimiter(rate.Limit(LIMIT), BURST)
}
