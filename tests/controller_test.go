package tests

import (
	"bytes"
	// "json"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"url-shortener/internal/controller"
	"url-shortener/internal/model"

	"url-shortener/tests/mocks"
)

func TestShortenEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockService := new(mocks.MockShortenerService)
	mockService.On("Shortening", "https://example.com").Return("test_short_url", nil).Once()

	handler := handler.NewHandler(mockService, nil)
	handler.Register(router)

	reqBody := model.LongURL{URL: "https://example.com"}
	req := httptest.NewRequest(http.MethodPost, "/shorten", convertToJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"short_url": "test_short_url"}`, w.Body.String())

	mockService.AssertExpectations(t)
}

func convertToJSON(data interface{}) *bytes.Buffer {
	body, _ := json.Marshal(data)
	return bytes.NewBuffer(body)
}
