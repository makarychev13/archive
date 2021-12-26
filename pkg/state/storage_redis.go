package state

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-redis/redis/v8"
)

//RedisStorage хранит текущий стейт в redis.
type RedisStorage struct {
	redis *redis.Client
}

//NewRedisStorage создаёт экземпляр RedisStorage.
func NewRedisStorage(opts redis.Options) *RedisStorage {
	return &RedisStorage{redis.NewClient(&opts)}
}

func (s *RedisStorage) Current(id int64) (Name, error) {
	val, err := s.redis.Get(context.Background(), toRedisKey(id)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}

	return val, nil
}

func (s *RedisStorage) Set(id int64, name Name) error {
	if err := s.redis.Set(context.Background(), toRedisKey(id), name, 0).Err(); err != nil {
		return err
	}

	return nil
}

func (s *RedisStorage) Clear(id int64) error {
	if err := s.redis.Del(context.Background(), toRedisKey(id)).Err(); err != nil {
		return err
	}

	return nil
}

func toRedisKey(id int64) string {
	return strconv.FormatInt(id, 10)
}
