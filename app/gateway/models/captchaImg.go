package models

import (
	"context"
	"star/app/storage/redis"
	"time"
)

type CaptchaStore struct {
}

func (s *CaptchaStore) Set(id string, value string) error {
	return redis.Client.Set(context.Background(), id, value, 1*time.Minute).Err()
}

func (s *CaptchaStore) Get(id string, clear bool) string {
	captcha := redis.Client.Get(context.Background(), id).String()
	if clear {
		redis.Client.Del(context.Background(), id)
	}
	return captcha
}

func (s *CaptchaStore) Verify(id, answer string, clear bool) bool {
	captcha := s.Get(id, clear)
	return captcha == answer
}
