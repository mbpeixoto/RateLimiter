## Modo de uso:
1. Configure os limiter do ratelimier nas variáveis de ambiente no arquivo docker-compose. Por padrão foi usado 10 req/s para ip e 100 req/s para token.

2. Rodar o comando:
```
docker-compose up -d
```
3. A partir dai a aplicação estará disponível para teste em:
`http://localhost:8080/ratelimit`

4. Use o arquivo da pasta /test para fazer requisições http que testam o sitema. Em main.go altere se necessário os limites de requisição conforme usado nas variáveis de ambiente do docker-compose. Após isso rode:
```
go run /test/main.go
```
