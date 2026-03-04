package data

import (
	"context"
	"regexp"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/intyouss/AI-Task-Hub/config"
	"github.com/intyouss/AI-Task-Hub/ent"
	_ "github.com/mattn/go-sqlite3"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type Data struct {
	Client *ent.Client
}

func NewDBClient(config config.Config) (*Data, func(), error) {
	client, err := setClient(config.Database)
	if err != nil {
		return nil, nil, err
	}

	if config.Mode == "debug" {
		client = client.Debug()
	}

	cleanup := func() {
		client.Close()
	}
	return &Data{
		Client: client,
	}, cleanup, nil
}

func setClient(dbConfig config.Database) (*ent.Client, error) {
	reg := regexp.MustCompile(`^file:\./(?P<name>[^.]+)\.db(?:\?.+)?$`)
	match := reg.FindStringSubmatch(dbConfig.Dsn)
	dbName := ""
	if len(match) > 0 {
		dbName = match[reg.SubexpIndex("name")]
	}
	otelDB, err := otelsql.Open(dbConfig.Driver, dbConfig.Dsn,
		otelsql.WithAttributes(semconv.DBSystemSqlite),
		otelsql.WithDBName(dbName))
	if err != nil {
		return nil, err
	}

	// 连接池配置
	otelDB.SetMaxOpenConns(20)                // 最大打开连接数
	otelDB.SetMaxIdleConns(5)                 // 最大空闲连接数
	otelDB.SetConnMaxLifetime(20 * time.Hour) // 连接最大生命周期
	otelDB.SetConnMaxIdleTime(10 * time.Minute)

	driver := sql.OpenDB(dialect.SQLite, otelDB)
	client := ent.NewClient(ent.Driver(driver))
	err = client.Schema.Create(context.Background())
	if err != nil {
		return nil, err
	}
	return client, nil
}
