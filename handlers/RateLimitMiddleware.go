package handlers

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-redis/redis"
)

type RateLimit struct {
	DB                *redis.Client
	LimiteRequisicoes int
	Tempo             time.Duration
}

func (rate *RateLimit) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		ipSemPorta, _, _ := net.SplitHostPort(ip)

		counter, err := rate.DB.Get(ipSemPorta).Int64()
		if err == redis.Nil {

			err = rate.DB.Set(ipSemPorta, 1, rate.Tempo).Err()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			counter = 1

		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {

			log.Println(counter, ipSemPorta)
			if counter >= int64(rate.LimiteRequisicoes) {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			} else {
				err = rate.DB.Incr(ipSemPorta).Err()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			next.ServeHTTP(w, r)
		}
	})
}
