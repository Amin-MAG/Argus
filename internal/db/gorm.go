package db

import (
	"argus/config"
	"argus/pkg/logger"
	"context"
	"fmt"
	"net/url"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	cfg config.Config
	db  *gorm.DB
}

// NewGormDB establishes a new connection to a database,
// It creates our schema, auto migrate and finally, creates new instance of DB
func NewGormDB(ctx context.Context, cfg config.Config, l logger.Logger) (DB, error) {
	c := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Name,
		cfg.Database.Password,
	)

	// Create the Gorm config
	gormCfg := &gorm.Config{}
	if cfg.Logger.SQLTraceLogEnable {
		gormCfg.Logger = l
	}

	// Create a new connection
	db, err := gorm.Open(postgres.Open(c), gormCfg)
	if err != nil {
		return nil, err
	}

	// Migration
	err = db.AutoMigrate(
		&Agent{},
	)

	return &GormDB{
		cfg: cfg,
		db:  db,
	}, nil
}

func NewGormDBWithURI(ctx context.Context, uri string, logger logger.Logger) (DB, error) {
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	password, _ := parsedURL.User.Password()

	var cfg config.Config
	cfg.Database.Host = parsedURL.Hostname()
	cfg.Database.Port = parsedURL.Port()
	cfg.Database.Name = strings.TrimLeft(parsedURL.Path, "/")
	cfg.Database.Username = parsedURL.User.Username()
	cfg.Database.Password = password

	return NewGormDB(ctx, cfg, logger)
}
