package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPing(t *testing.T) {
	ctx := context.Background()

	// Get test database instance
	tdb := getTestDatabase(ctx, t)

	// Ping the database
	err := tdb.Ping(ctx)
	assert.NoError(t, err, "error pinging database")
}
