package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"rate-limiter/limiter"
	"rate-limiter/storage"

	"github.com/gin-gonic/gin"
)

func setupTestRouter(limiter *limiter.RateLimiter) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimitMiddleware(limiter))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	return router
}

func TestRateLimitMiddleware_AllowRequest(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	config := &limiter.Config{
		IPLimit:            5,
		IPBlockDuration:    300,
		TokenLimit:         10,
		TokenBlockDuration: 300,
	}

	rateLimiter := limiter.NewRateLimiter(mockStorage, config)
	router := setupTestRouter(rateLimiter)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Verificar que a requisição foi permitida
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestRateLimitMiddleware_BlockRequest(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	config := &limiter.Config{
		IPLimit:            1,
		IPBlockDuration:    300,
		TokenLimit:         10,
		TokenBlockDuration: 300,
	}

	rateLimiter := limiter.NewRateLimiter(mockStorage, config)
	router := setupTestRouter(rateLimiter)

	req1, _ := http.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("First request should be allowed, got status %d", w1.Code)
	}

	// Essa aqui tem que ser bloqueada
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "192.168.1.1:12345"
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("Second request should be blocked, got status %d", w2.Code)
	}
}

func TestRateLimitMiddleware_TokenBased(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	config := &limiter.Config{
		IPLimit:            1,
		IPBlockDuration:    300,
		TokenLimit:         2,
		TokenBlockDuration: 300,
	}

	rateLimiter := limiter.NewRateLimiter(mockStorage, config)
	router := setupTestRouter(rateLimiter)

	req1, _ := http.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	req1.Header.Set("API_KEY", "test-token")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("First request with token should be allowed, got status %d", w1.Code)
	}

	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "192.168.1.1:12345"
	req2.Header.Set("API_KEY", "test-token")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Second request with token should be allowed, got status %d", w2.Code)
	}

	req3, _ := http.NewRequest("GET", "/test", nil)
	req3.RemoteAddr = "192.168.1.1:12345"
	req3.Header.Set("API_KEY", "test-token")
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusTooManyRequests {
		t.Errorf("Third request with token should be blocked, got status %d", w3.Code)
	}
}

func TestRateLimitMiddleware_BlockedIP(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	config := &limiter.Config{
		IPLimit:            5,
		IPBlockDuration:    300,
		TokenLimit:         10,
		TokenBlockDuration: 300,
	}

	mockStorage.SetBlocked("ip:192.168.1.1", true)

	rateLimiter := limiter.NewRateLimiter(mockStorage, config)
	router := setupTestRouter(rateLimiter)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Request from blocked IP should be blocked, got status %d", w.Code)
	}
}

func TestRateLimitMiddleware_BlockedToken(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	config := &limiter.Config{
		IPLimit:            5,
		IPBlockDuration:    300,
		TokenLimit:         10,
		TokenBlockDuration: 300,
	}

	mockStorage.SetBlocked("token:test-token", true)

	rateLimiter := limiter.NewRateLimiter(mockStorage, config)
	router := setupTestRouter(rateLimiter)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	req.Header.Set("API_KEY", "test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Request with blocked token should be blocked, got status %d", w.Code)
	}
}

func TestRateLimitMiddleware_ClientIPExtraction(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	config := &limiter.Config{
		IPLimit:            5,
		IPBlockDuration:    300,
		TokenLimit:         10,
		TokenBlockDuration: 300,
	}

	rateLimiter := limiter.NewRateLimiter(mockStorage, config)
	router := setupTestRouter(rateLimiter)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "10.0.0.1")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Request should be allowed, got status %d", w.Code)
	}

	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.Header.Set("X-Real-IP", "10.0.0.2")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Request should be allowed, got status %d", w2.Code)
	}
}
