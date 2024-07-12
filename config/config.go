package config

import "strings"

type Config struct {
	Argus struct {
		// Global
		IsProductionMode bool   `env:"IS_PRODUCTION_MODE" env-default:"false" env-description:"Is in production mode"`
		Version          string `env:"ARGUS_VERSION" env-default:"local" env-description:"The version of Argus Service"`

		// Gin
		GinMode string `env:"GIN_MODE" env-default:"debug" env-description:"Gin framework mode (release or debug)"`
		Port    string `env:"SERVING_PORT" env-default:"8081" env-description:"Port number for Argus API"`
	}
	Database struct {
		Host     string `env:"POSTGRES_HOST" env-default:"localhost" env-description:"Database host for service"`
		Port     string `env:"POSTGRES_PORT" env-default:"5432" env-description:"Database port for service"`
		Name     string `env:"POSTGRES_NAME" env-default:"argus_db" env-description:"Database name for service"`
		Username string `env:"POSTGRES_USERNAME" env-default:"argus" env-description:"Database username for service"`
		Password string `env:"POSTGRES_PASSWORD" env-default:"test" env-description:"Database password for service"`
	}
	Logger struct {
		Level              string `env:"LOGGER_LEVEL" env-default:"debug" env-description:"Log Level for application logger"`
		SQLTraceLogEnable  bool   `env:"LOGGER_SQL_TRACE_LOG_ENABLE" env-default:"false" env-description:"Does the logger print low level SQL logs"`
		IsReportCallerMode bool   `env:"LOGGER_IS_REPORT_CALLER_MODE" env-default:"false" env-description:"Does the logger have report caller"`
		IsPrettyPrint      bool   `env:"LOGGER_PRETTY_PRINT" env-default:"false" env-description:"Pretty JSON Print flag"`
		FilePath           string `env:"LOGGER_FILE" env-default:"service.log" env-description:"The file that stores the logs of service"`
	}
}

// maskString Masks sensitive information with asterisks from string
func maskString(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + strings.Repeat("*", len(s)-4) + s[len(s)-2:]
}

// maskConfig masks the data
func maskConfig(field *string) {
	*field = maskString(*field)
}

// SecureClone creates a secure instance of Config with masking sensitive information
func (c Config) SecureClone() Config {
	sc := c

	// Censor critical values
	maskConfig(&sc.Database.Password)

	return sc
}
