package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/intyouss/AI-Task-Hub/config"
	"github.com/intyouss/AI-Task-Hub/internal/data"
)

func main() {
	serverRun(":8080")
}

func serverRun(port string) {
	cfg := config.Load()

	switch cfg.Mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)

	srv := &http.Server{Addr: port, Handler: r}

	_, cleanup, err := data.NewDBClient(cfg)
	if err != nil {
		slog.Error("数据库连接失败", "err", err)
		return
	}
	defer cleanup()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("服务器错误", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
