package handlers

import (
	"github.com/go-redis/redis"
	"net/http"
	"time"
)

type RateLimit struct {
	DB *redis.Client
}

func (rate *RateLimit) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter, err := rate.DB.Get(r.RemoteAddr).Int64()
		if err == redis.Nil {
			err = rate.DB.Set(r.RemoteAddr, 1, time.Second*20).Err()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			counter = 1

		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			if counter >= 10 {
				http.Error(w, "Limite excedido", http.StatusTooManyRequests)
				return
			} else {
				err = rate.DB.Incr(r.RemoteAddr).Err()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
	})
}
