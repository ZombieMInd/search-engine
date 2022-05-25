package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZombieMInd/go-logger/internal/logger"
	"github.com/ZombieMInd/go-logger/internal/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestServer_HandleLog(t *testing.T) {
	s := NewServer(teststore.New())
	s.configLogger(&Config{})
	s.InitServices(&Config{})
	initRouter(s)
	var events []interface{}
	event := logger.Events{}
	events = append(events, event)
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"user_uuid": "bd4cb967-a824-4ada-ad75-f74820793819",
				"timestamp": 2987428975,
				"events":    events,
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/log", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
