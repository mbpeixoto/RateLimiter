package handlers

import (
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type RateLimitConfig struct {
	DB                     *redis.Client
	LimiteRequisicoesToken string        `mapstructure:"LIMITE_REQUISICOES_TOKEN"`
	LimiteRequisicoesIP    string        `mapstructure:"LIMITE_REQUISICOES_IP"`
	TempoExpiracao         time.Duration `mapstructure:"TEMPO_EXPIRACAO"`
	TempoBloqueio          time.Duration `mapstructure:"TEMPO_BLOQUEIO"`
}

type RateLimiter interface {
	Connect() error
	Close() error
	ContarRequisicoes(chave string, duracao time.Duration) (int64, error)
	Incrementar(chave string) error
	Bloquear(chave string, duracao time.Duration) error
	EstaBloqueado(chave string) (bool, error)
}

func RateLimitMiddleware(rateLimiter RateLimiter, rateLimitConfig RateLimitConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Conecta ao Redis
			err := rateLimiter.Connect()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Erro ao conectar ao Redis:", err)
				return
			}
			defer rateLimiter.Close()

			token := r.Header.Get("API_KEY")
			chave := ""

			if token != "" {
				chave = token
			} else {
				ip := r.RemoteAddr
				ipSemPorta, _, _ := net.SplitHostPort(ip)
				chave = ipSemPorta
			}

			// Verifica se está bloqueado
			bloqueado, err := rateLimiter.EstaBloqueado(chave)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Erro ao verificar bloqueio:", err)
				return
			}
			if bloqueado {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("Cliente bloqueado devido a excesso de requisições"))
				return
			}

			// Conta requisições
			contador, err := rateLimiter.ContarRequisicoes(chave, rateLimitConfig.TempoExpiracao)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Erro ao contar requisições:", err)
				return
			}

			// Verifica o limite de requisições
			limiteRequisicoes := rateLimitConfig.LimiteRequisicoesIP
			if token != "" {
				limiteRequisicoes = rateLimitConfig.LimiteRequisicoesToken
			}

		
			num, err := strconv.Atoi(limiteRequisicoes)
			if err != nil {
				log.Println("Erro ao converter string para int:", err)
				return
			}

			log.Println("Contador de requisições para chave", chave, ":", contador)
			log.Print("Limite de requisições para chave", chave, ":", limiteRequisicoes)

			if contador > int64(num) {
				log.Println("Limite de requisições excedido, aplicando bloqueio para chave:", chave)
				err = rateLimiter.Bloquear(chave, rateLimitConfig.TempoBloqueio)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Println("Erro ao aplicar bloqueio:", err)
				} else {
					w.WriteHeader(http.StatusTooManyRequests)
					w.Write([]byte("Limite de requisições excedido, cliente bloqueado"))
				}
				return
			} else {
				err = rateLimiter.Incrementar(chave)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Println("Erro ao incrementar contador de requisições:", err)
					return
				}
			}

		

			// Se o limite não foi excedido, prossiga para o próximo handler
			next.ServeHTTP(w, r)
		})
	}
}
