package pkg

import "github.com/redis/go-redis/v9"

func NewRedisClient(conf *Redis) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     conf.Uri,
		Username: conf.Username,
		Password: conf.Password,
		DB:       conf.Db,
	})
}
