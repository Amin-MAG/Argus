package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskString(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"password", "pa****rd"},
		{"1234", "****"},
		{"token1234", "to*****34"},
		{"", "****"},
	}

	for _, tc := range testCases {
		output := maskString(tc.input)
		assert.Equal(t, tc.expected, output, "Expected maskString(%s) to be %s", tc.input, tc.expected)
	}
}

func TestMaskConfig(t *testing.T) {
	field := "password"
	maskConfig(&field)
	expected := "pa****rd"
	assert.Equal(t, expected, field, "Expected maskConfig to modify field to %s", expected)
}

func TestSecureClone(t *testing.T) {
	c := Config{
		Argus: struct {
			IsProductionMode bool   `env:"IS_PRODUCTION_MODE" env-default:"false" env-description:"Is in production mode"`
			Version          string `env:"ARGUS_VERSION" env-default:"local" env-description:"The version of Argus Service"`
			GinMode          string `env:"GIN_MODE" env-default:"debug" env-description:"Gin framework mode (release or debug)"`
			Port             string `env:"SERVING_PORT" env-default:"8081" env-description:"Port number for Argus API"`
			APIKey           string `env:"API_KEY" env-default:"test_api_key" env-description:"API key"`
		}{
			IsProductionMode: true,
			Version:          "1.0.0",
			GinMode:          "release",
			Port:             "8080",
		},
		IPInfo: struct {
			DefaultTimeoutInSecs int64  `env:"IP_INFO_DEFAULT_TIMEOUT_IN_SECS" env-default:"5" env-description:"Default timeout in seconds"`
			Token                string `env:"IP_INFO_TOKEN" env-default:"<secret>" env-description:"Token used to connect to IP Info API"`
		}{
			DefaultTimeoutInSecs: 10,
			Token:                "my-secret-token",
		},
		Database: struct {
			Host     string `env:"POSTGRES_HOST" env-default:"localhost" env-description:"Database host for service"`
			Port     string `env:"POSTGRES_PORT" env-default:"5432" env-description:"Database port for service"`
			Name     string `env:"POSTGRES_NAME" env-default:"argus_db" env-description:"Database name for service"`
			Username string `env:"POSTGRES_USERNAME" env-default:"argus" env-description:"Database username for service"`
			Password string `env:"POSTGRES_PASSWORD" env-default:"test" env-description:"Database password for service"`
		}{
			Host:     "db.example.com",
			Port:     "5432",
			Name:     "prod_db",
			Username: "admin",
			Password: "secure-password",
		},
		Logger: struct {
			Level              string `env:"LOGGER_LEVEL" env-default:"debug" env-description:"Log Level for application logger"`
			SQLTraceLogEnable  bool   `env:"LOGGER_SQL_TRACE_LOG_ENABLE" env-default:"false" env-description:"Does the logger print low level SQL logs"`
			IsReportCallerMode bool   `env:"LOGGER_IS_REPORT_CALLER_MODE" env-default:"false" env-description:"Does the logger have report caller"`
			IsPrettyPrint      bool   `env:"LOGGER_PRETTY_PRINT" env-default:"false" env-description:"Pretty JSON Print flag"`
			FilePath           string `env:"LOGGER_FILE" env-default:"service.log" env-description:"The file that stores the logs of service"`
		}{
			Level:              "info",
			SQLTraceLogEnable:  true,
			IsReportCallerMode: true,
			IsPrettyPrint:      false,
			FilePath:           "/var/log/argus.log",
		},
	}

	sc := c.SecureClone()

	// Check sensitive fields are masked
	assert.True(t, strings.HasPrefix(sc.Database.Password, "se**"), "Expected Database.Password to be masked, got %s", sc.Database.Password)
	assert.True(t, strings.HasPrefix(sc.IPInfo.Token, "my****"), "Expected IPInfo.Token to be masked, got %s", sc.IPInfo.Token)

	// Check non-sensitive fields remain unchanged
	assert.Equal(t, c.Argus.Version, sc.Argus.Version, "Expected Argus.Version to be %s, got %s", c.Argus.Version, sc.Argus.Version)
	assert.Equal(t, c.Logger.Level, sc.Logger.Level, "Expected Logger.Level to be %s, got %s", c.Logger.Level, sc.Logger.Level)
}
