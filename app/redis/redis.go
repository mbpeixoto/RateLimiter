package redisdb

import (
	"time"

	"github.com/go-redis/redis"
)

// RedisClient para gerenciar tentativas de login
type RedisClient struct {
	Client *redis.Client
}


func (r *RedisClient) ConnectRedis() *RedisClient {

	redisdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // use default Addr
		Password: "",               // no password set
		DB:       0,                // use default DB
	})
	return &RedisClient{
		Client: redisdb,
	}
}

func (r *RedisClient) ContarRequisicoes( key string, tempo time.Duration) (int64, error) {
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

func (r *RedisClient) CloseRedis() error {
	return r.Client.Close()
}

// Middleware de limitação de tentativas
const MAX_ATTEMPTS = 5

func (r *RedisClient) IncrementLoginAttempts(userName string) (int, error) {
	attempts, err := r.Client.Incr(userName).Result()
	if err != nil {
		return 0, err
	}
	if attempts == 1 {
		r.Client.Expire(userName, time.Minute*15).Err()
	}
	return int(attempts), nil
}

func (r *RedisClient) ResetLoginAttempts(userName string) error {
	return r.Client.Del(userName).Err()
}

func (r *RedisClient) GetLoginAttempts(userName string) (int, error) {
	attempts, err := r.Client.Get(userName).Int()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return attempts, nil
}