package main

import (
	"argus/config"
	"argus/internal/db"
	"argus/internal/routes"
	"argus/pkg/logger"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load the config
	var cfg config.Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf(`
___  ____ ____ _    ___  
|  \ | . \|  _\|| \ | _\ 
| . \|  <_| [ \||_|\[__ \
|/\_/|/\_/|___/|___/|___/

Starting Argus-%s...
Configuration: %+v
`, cfg.Argus.Version, cfg)

	// Setup the logger
	log := logrus.New()
	logLevel, err := logrus.ParseLevel(cfg.Logger.Level)
	if err != nil {
		logLevel = logrus.DebugLevel
	}
	log.SetLevel(logLevel)
	log.SetReportCaller(cfg.Logger.IsReportCallerMode)
	if cfg.Logger.IsPrettyPrint {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	}
	logger.SetupLogger(log)
	logger.Info("logger is setup successfully")

	// Create new instance of dg.DB
	gormDB, err := db.NewGormDB(cfg, logger.GetLogger())
	if err != nil {
		log.WithError(err).Fatal("error in connecting to the postgres database")
	}
	log.Info("connected to the argus database")

	// Create Gin HTTP Server
	s, err := routes.NewGinServer(cfg, gormDB)
	if err != nil {
		log.WithError(err).Fatal("error in creating API server")
	}
	log.Info("created the gin server")

	// Start Listening and Serving
	log.WithField("port", cfg.Argus.Port).Info("the server is going to be started")
	log.WithError(s.ListenAndServe()).Fatal("")
}
