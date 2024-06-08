package service

import (
	"context"
	"time"
)

func (s *Srv) setRCache(key string, data []byte, duration time.Duration) error {
	return s.rClient.Set(context.Background(),
		key,
		data,
		duration).Err()
}
