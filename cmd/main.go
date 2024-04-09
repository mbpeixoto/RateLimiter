package main

import (
	"net/http"
	"ratelimiter/handlers"
	"time"

	"log"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

func main() {
	redisdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // use default Addr
		Password: "",               // no password set
		DB:       0,                // use default DB
	})

	// Conexão com o Redis,Limite de requisições
	ratelimit := handlers.RateLimit{DB: redisdb, LimiteRequisicoes: 10, Tempo: time.Second * 60}
	log.Println("Limite de requisições:",ratelimit.LimiteRequisicoes, "Tempo:", ratelimit.Tempo)

	r := mux.NewRouter()
	r.Use(ratelimit.RateLimitMiddleware)
	r.HandleFunc("/", handlers.HomeServer).Methods("GET")
	http.ListenAndServe(":8080", r)
}
