package handlers

import (
	"argus/config"
	"argus/internal/db"
	"argus/internal/iputil"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Test handler
func TestPing(t *testing.T) {
	ctx := context.Background()
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		db           db.DB
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Database Up",
			db:           getTestDatabase(ctx, t),
			expectedCode: http.StatusOK,
			expectedBody: `{"database_status":"up"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			argusIpClient, err := iputil.NewArgusIPClient(&iputil.MockIPInfoClient{})
			assert.NoError(t, err)
			assert.NotNil(t, argusIpClient)

			gh := NewGinHandler(config.Config{}, tt.db, argusIpClient)

			router := gin.Default()
			router.GET("/ping", gh.Ping)

			req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
