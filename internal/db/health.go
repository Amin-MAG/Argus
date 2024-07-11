package db

import (
	"context"
)

// Ping checks the status of database connection
func (gdb *GormDB) Ping(ctx context.Context) error {
	gormDB, err := gdb.db.DB()
	if err != nil {
		return err
	}

	return gormDB.PingContext(ctx)
}
