package db

import (
	"argus/config"
	"argus/pkg/logger"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	cfg config.Config
	db  *gorm.DB
}

// NewGormDB establishes a new connection to a database,
// It creates our schema, auto migrate and finally, creates new instance of DB
func NewGormDB(cfg config.Config, l logger.Logger) (DB, error) {
	c := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
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
