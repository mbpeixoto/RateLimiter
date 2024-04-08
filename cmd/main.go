package main

import (
	"net/http"
	"ratelimiter/handlers"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

func main(){
	redisdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // use default Addr
		Password: "",               // no password set
		DB:       0,                // use default DB
	})

	ratelimit := handlers.RateLimit{DB: redisdb}

	r := mux.NewRouter()
	r.Use(ratelimit.RateLimitMiddleware)
	r.HandleFunc("/", handlers.ServeHTTP).Methods("GET")
	http.ListenAndServe(":8080", r)
}