

## Funcionamento
- O ratelimiter funciona criando uma chave valor baseada no token ou ip no redis. Vai incrementando a cada requisição e verificando o limite. Se o limite for atingido cria-se um chave valor de bloqueio. O limite de requisições por token  e ip, o tempo desse limite (expiração) e o tempo de bloqueio quando o limite é atingido pode ser definido nas variáveis de ambiente em /cmd/.env
Exemplo:
```
LIMITE_REQUISICOES_TOKEN=20
LIMITE_REQUISICOES_IP=10
TEMPO_EXPIRACAO=10s
TEMPO_BLOQUEIO=30s
```
Como as configurações do exemplo acima o limite por token é 20 req/10s e o limite de requisições por ip é 10 req/10s sendo que o tempo de bloqueio após atingir tais limites é de 30s

## Modo de uso:
1. Configure os limiter do ratelimier nas variáveis de ambiente no arquivo docker-compose. Por padrão foi usado 10 req/s para ip e 100 req/s para token.

2. Rodar o comando:
```
docker-compose up -d
```
3. A partir dai a aplicação estará disponível para teste em:
`http://localhost:8080/ratelimit`

4. Use o arquivo da pasta /cmd para fazer teste automatizado:
```
docker-compose run --rm test

```
5. Se preferir podemos testar via curl:
- Sem token:
```
curl -v -X GET localhost:8080/ratelimit

```
- Com token:
```
curl -v -X GET  -H "API_KEY: abc" localhost:8080/ratelimit

```

- Para enviar várias requisições adicone como parâmetro da url "?[1-n]", exemplo:
```
curl -v -X GET localhost:8080/ratelimit?[1-20]

```

## Strategy

Para aplicar o padrão Strategy, defini uma interface que abstrai as operações de rate limiting e implementar essa interface para diferentes armazenamentos.
- Interface de Estratégia:
```
type RateLimiter interface {
    Connect() error
    Close() error
    ContarRequisicoes(chave string, duracao time.Duration) (int64, error)
    Incrementar(chave string) error
}
```
- Estratégia Redis:
```
type RedisRateLimiter struct {
    Client redis.Client
}

func (r *RedisRateLimiter) Connect() error {
    // Código para conectar ao Redis
    return nil
}

func (r *RedisRateLimiter) Close() error {
    // Código para desconectar do Redis
    return nil
}

func (r *RedisRateLimiter) ContarRequisicoes(chave string, duracao time.Duration) (int64, error) {
    // Código para contar requisições no Redis
    return 0, nil
}

func (r *RedisRateLimiter) Incrementar(chave string) error {
    // Código para incrementar contador no Redis
    return nil
}

```