package main

import (
	"github.com/gin-gonic/gin"
	"github.com/intyouss/AI-Task-Hub/config"
	"github.com/intyouss/AI-Task-Hub/internal/middleware"
)

func main() {
	// Parse config
	config.Parse()

	r := gin.Default()
	r.Use(middleware.RateLimiter(20))

	r.Run()
}
