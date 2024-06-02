package main

import (
	"log"
	"net/http"
)

// altere os valores para testar o rate limiter
const limiteRequisicoesIp = 10
const limiteRequisicoesToken = 100

func TestRateLimiterByIP() {

	client := &http.Client{}

	log.Print("Testando RateLimiter por IP")

	
	for i := 1; i <= (limiteRequisicoesIp+1); i++ {
		resp, err := client.Get("http://localhost:8080/ratelimit")
		if err != nil {
			log.Printf("Falha ao fazer a requisição #%d: %v", i, err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("Requisição #%d status %d", i, resp.StatusCode)
		}
	}
}

func TestRateLimiterByToken() {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "http://localhost:8080/ratelimit", nil)
	req.Header.Set("API_KEY", "token123")

	log.Print("Testando RateLimiter por Token")
	
	for i := 1; i <= (limiteRequisicoesToken+1); i++ {
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Falha ao fazer a requisição #%d: %v", i, err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("Requisição #%d status %d", i, resp.StatusCode)
		}
	}
}

func main() {
	TestRateLimiterByIP()
	TestRateLimiterByToken()
}
