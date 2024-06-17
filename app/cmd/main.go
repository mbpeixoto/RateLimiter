package main

import (
	"net/http"

	configs "ratelimiter/config"
	"ratelimiter/handlers"
	redisdb "ratelimiter/redis"
	
	"github.com/gorilla/mux"
)

func main() {

	rateLimitConfigs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	redisRateLimiter := &redisdb.RedisRateLimiter{}


	r := mux.NewRouter()
	r.Use(handlers.RateLimitMiddleware(redisRateLimiter, *rateLimitConfigs))
	r.HandleFunc("/ratelimit", handlers.HomeServer).Methods("GET")
	http.ListenAndServe(":8080", r)
}
