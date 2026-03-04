package config

import (
	"embed"
	"log/slog"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Database Database
	Mode     string
}

type Database struct {
	Driver string `yaml:"driver"`
	Dsn    string `yaml:"dsn"`
}

//go:embed config.yaml
var fs embed.FS

func Load() Config {
	file, err := fs.ReadFile("config.yaml")
	if err != nil {
		slog.Error("配置文件加载失败", "error", err)
		panic(err)
	}

	var Cfg Config
	err = yaml.Unmarshal(file, &Cfg)
	if err != nil {
		panic(err)
	}
	slog.Info("配置文件加载成功")
	return Cfg
}
