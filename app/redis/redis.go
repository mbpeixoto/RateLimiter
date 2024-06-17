package redisdb

import (
	"time"

	"github.com/go-redis/redis"
)

// RedisRateLimiter para gerenciar tentativas de login
type RedisRateLimiter struct {
	Client *redis.Client
}


func (r *RedisRateLimiter) Connect() error {

	redisdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // use default Addr
		Password: "",               // no password set
		DB:       0,                // use default DB
	})
	r.Client = redisdb
	return nil
}

func (r *RedisRateLimiter) ContarRequisicoes( key string, tempo time.Duration) (int64, error) {
	counter, err := r.Client.Get(key).Int64()

	if err == redis.Nil {

		err = r.Client.Set(key, 1, tempo).Err()
		if err != nil {
			return 0, err
		}
		counter = 1
	} else if err != nil {
		return 0, err
	}
	return counter, nil
}

func (r *RedisRateLimiter) Close() error {
	return r.Client.Close()
}


func (r *RedisRateLimiter) Incrementar(key string) error {
	return r.Client.Incr(key).Err()
}

func (r *RedisRateLimiter) Bloquear(key string, tempo time.Duration) error {
	return r.Client.Set(key+":blocked", 1, tempo).Err()
}

func (r *RedisRateLimiter) EstaBloqueado(key string) (bool, error) {
    resultado, err := r.Client.Get(key+":blocked").Result()
    if err == redis.Nil {
        return false, nil
    }
    if err != nil {
        return false, err
    }
    return resultado == "1", nil
}