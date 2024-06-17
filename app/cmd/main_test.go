package main

import (
	"net/http"
	"net/http/httptest"
	"ratelimiter/handlers"
	redisdb "ratelimiter/redis"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func setupRedis() *redisdb.RedisRateLimiter {
	rateLimiter := redisdb.RedisRateLimiter{}
	rateLimiter.Client = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // use default Addr
		Password: "",               // no password set
		DB:       0,                // use default DB
	})
	rateLimiter.Client.FlushDB() // Limpa o banco de dados Redis antes dos testes
	return &rateLimiter
}

func TestRateLimitMiddleware(t *testing.T) {
	rateLimiter := setupRedis()
	defer rateLimiter.Close()

	rateLimitConfig := handlers.RateLimitConfig{
		LimiteRequisicoesToken: "3",
		LimiteRequisicoesIP:    "3",
		TempoExpiracao:                  time.Second * 10,
		TempoBloqueio:          time.Second * 20,
	}

	middleware := handlers.RateLimitMiddleware(rateLimiter, rateLimitConfig)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Requisição bem-sucedida"))
	}))

	tests := []struct {
		name           string
		apiKey         string
		expectedStatus int
		sleep          time.Duration
	}{
		{"Primeira requisição com API_KEY", "test-api-key", http.StatusOK, 0},
		{"Segunda requisição com API_KEY", "test-api-key", http.StatusOK, 0},
		{"Terceira requisição com API_KEY", "test-api-key", http.StatusOK, 0},
		{"Quarta requisição API_KEY", "test-api-key", http.StatusTooManyRequests, 0},
		{"Requisição após período de bloqueio", "test-api-key", http.StatusOK, time.Second * 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.sleep > 0 {
				time.Sleep(tt.sleep)
			}

			req, err := http.NewRequest("GET", "/api/", nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.apiKey != "" {
				req.Header.Set("API_KEY", tt.apiKey)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestRateLimitMiddlewareByIP(t *testing.T) {
	rateLimiter := setupRedis()
	defer rateLimiter.Close()

	rateLimitConfig := handlers.RateLimitConfig{
		LimiteRequisicoesToken: "3",
		LimiteRequisicoesIP:    	"3",
		TempoExpiracao:                  time.Second * 10,
		TempoBloqueio:          time.Second * 20,
	}

	middleware := handlers.RateLimitMiddleware(rateLimiter, rateLimitConfig)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Requisição bem-sucedida"))
	}))

	tests := []struct {
		name           string
		ip             string
		expectedStatus int
		sleep          time.Duration
	}{
		{"Primeira requisição com IP", "192.168.0.1:1234", http.StatusOK, 0},
		{"Segunda requisição com  IP", "192.168.0.1:1234", http.StatusOK, 0},
		{"Terceira requisição com IP", "192.168.0.1:1234", http.StatusOK, 0},
		{"Quarta requisição com IP", "192.168.0.1:1234", http.StatusTooManyRequests, 0},
		{"Requisição depois do período de bloqueio ", "192.168.0.1:1234", http.StatusOK, time.Second * 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.sleep > 0 {
				time.Sleep(tt.sleep)
			}

			req, err := http.NewRequest("GET", "/api/", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.RemoteAddr = tt.ip

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
