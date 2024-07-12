package handlers

import (
	"argus/internal/db"
	"argus/pkg/logger"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"testing"
	"time"
)

const (
	dbImage = "docker.io/postgres:16-alpine"

	dbName     = "test_db"
	dbUser     = "test_uer"
	dbPassword = "test_pass"
)

var postgresContainer *postgres.PostgresContainer

func TestMain(m *testing.M) {
	var code int

	// Setup the logger
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	logger.SetupLogger(log)

	// Run the container
	ctx := context.Background()

	var err error
	postgresContainer, err = postgres.Run(ctx,
		dbImage,
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		logger.WithError(err).Fatal(ctx, "can not start the postgres test container")
	}
	logger.Info(ctx, "successfully started the postgresql container")

	// Execute the tests
	code = m.Run()

	// Terminate the container
	err = postgresContainer.Terminate(ctx)
	if err != nil {
		logger.WithError(err).Warn(ctx, "can not terminate the postgresql container")
	}
	logger.Info(ctx, "successfully terminated the postgresql container")

	os.Exit(code)
}

func getTestDatabase(ctx context.Context, t *testing.T) db.DB {
	// Setup the logger
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	logger.SetupLogger(log)

	// Get the connection string
	uri, err := postgresContainer.ConnectionString(ctx)
	assert.NoError(t, err)

	// Create the db instance of the database
	tdb, err := db.NewGormDBWithURI(ctx, uri, logger.GetLogger())
	assert.NoError(t, err, "can not connect to the database", err)

	return tdb
}
