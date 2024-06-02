package handlers

import (
	"log"
	"net"
	"net/http"
	redisdb "ratelimiter/redis"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

type RateLimit struct {
	DB                     *redis.Client
	LimiteRequisicoesToken int
	LimiteRequisicoesIP    int
	Tempo                  time.Duration
}

func RateLimiRequesttMiddleware(rateLimit RateLimit) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			redis := redisdb.RedisClient{}
			redisCliente := redis.ConnectRedis()
			defer redisCliente.CloseRedis()

			token := r.Header.Get("API_KEY")
			if token != "" {
				counter, err := redisCliente.ContarRequisicoes(token, rateLimit.Tempo)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Println(err)
				} else {
					if counter > int64(rateLimit.LimiteRequisicoesToken) {
						w.WriteHeader(http.StatusTooManyRequests)
						w.Write([]byte("Limite de requisições excedido"))
						return
					} else {
						err = redisCliente.Client.Incr(token).Err()
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							log.Println(err)
							next.ServeHTTP(w, r)
						}
					}
				}
			} else {

				ip := r.RemoteAddr
				ipSemPorta, _, _ := net.SplitHostPort(ip)

				//rateLimit := RateLimit{LimiteRequisicoes: 10, Tempo: time.Second * 30}

				counter, err := redisCliente.ContarRequisicoes(ipSemPorta, rateLimit.Tempo)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Println(err)
				} else {
					if counter > int64(rateLimit.LimiteRequisicoesIP) {
						w.WriteHeader(http.StatusTooManyRequests)
						w.Write([]byte("Limite de requisições excedido"))
						return
					} else {
						err = redisCliente.Client.Incr(ipSemPorta).Err()
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							log.Println(err)
							return
						}
					}
					next.ServeHTTP(w, r)
				}
			}
		})
	}
}
