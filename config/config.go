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

		// APIKey: Just for testing purposes, API_Keys of agents should be managed by a separate table
		APIKey string `env:"API_KEY" env-default:"test_api_key" env-description:"API key"`
	}
	IPInfo struct {
		DefaultTimeoutInSecs int64  `env:"IP_INFO_DEFAULT_TIMEOUT_IN_SECS" env-default:"5" env-description:"Default timeout in seconds"`
		Token                string `env:"IP_INFO_TOKEN" env-default:"<secret>" env-description:"Token used to connect to IP Info API"`
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
	}
	Tracing struct {
		Enabled      bool    `env:"IS_TRACING_ENABLED"  env-default:"false" env-description:"activate tracing"`
		Endpoint     string  `env:"JAEGER_ENDPOINT" env-default:"http://localhost:14268/api/traces" env-description:"Endpoint to send tracing requests"`
		ServiceName  string  `env:"JAEGER_SERVICE_NAME" env-default:"argus" env-description:"URL for backup generation"`
		SamplerRatio float64 `env:"TRACING_SAMPLE_RATIO" env-default:"0.5" env-description:"ratio to send tracing info"`
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
	maskConfig(&sc.IPInfo.Token)
	maskConfig(&sc.Argus.APIKey)

	return sc
}
