package main

import (
	"log"
	"net/http"
	"os"
	"ratelimiter/handlers"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	limiteRequisicoesToken, err := strconv.Atoi(os.Getenv("LIMITE_REQUISICOES_TOKEN"))
    if err != nil {
        log.Println("Erro ao converter LIMITE_REQUISICOES_IP")
    }

	limiteRequisicoesIp, err := strconv.Atoi(os.Getenv("LIMITE_REQUISICOES_IP"))
	if err != nil {
		log.Println("Erro ao converter LIMITE_REQUISICOES_IP")
	}

    tempo, err := strconv.Atoi(os.Getenv("TEMPO"))
    if err != nil {
        log.Println("Erro ao converter TEMPO")
    }

    rateLimit := handlers.RateLimit{LimiteRequisicoesToken: limiteRequisicoesToken, LimiteRequisicoesIP: limiteRequisicoesIp,Tempo: time.Second * time.Duration(tempo)}


	r := mux.NewRouter()
	r.Use(handlers.RateLimiRequesttMiddleware(rateLimit))
	r.HandleFunc("/teste", handlers.HomeServer).Methods("GET")
	http.ListenAndServe(":8080", r)
}
