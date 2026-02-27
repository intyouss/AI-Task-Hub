package config

import (
	"embed"
	"log/slog"
	"time"

	"github.com/goccy/go-yaml"
)

var Cfg Config

type Config struct {
	RateLimiter RateLimiter
}

type RateLimiter struct {
	Enable    bool          `yaml:"enable"`
	Limit     int           `yaml:"limit"`
	EveryTime time.Duration `yaml:"every_time"`
}

//go:embed config.yaml
var fs embed.FS

func Parse() {
	file, err := fs.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, &Cfg)
	if err != nil {
		panic(err)
	}
	slog.Info("配置文件加载成功")
}
