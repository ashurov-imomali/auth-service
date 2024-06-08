package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

func (s *Srv) setRCache(key string, data []byte, duration time.Duration) error {
	return s.rClient.Set(context.Background(),
		key,
		data,
		duration).Err()
}

func (s *Srv) getRCache(key string, data interface{}) (bool, error) {
	cmd := s.rClient.Get(context.Background(), key)
	bytes, err := cmd.Bytes()
	if err != nil {
		return false, err
	}
	return errors.Is(err, redis.Nil), json.Unmarshal(bytes, data)
}
