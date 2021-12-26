package state

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var (
	redisOpts   redis.Options
	redisClient *redis.Client
)

func init() {
	conf := "../../test.env"
	if err := godotenv.Load(conf); err != nil {
		log.Fatalf("Не удалось загрузить конфиг '%v': %v", conf, err)
	}

	redisOpts = redis.Options{
		Addr: fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	}

	redisClient = redis.NewClient(&redisOpts)
}

func TestRedisStorage_SetNewEntity(t *testing.T) {
	var (
		state    = gofakeit.LetterN(10)
		entityId = gofakeit.Int64()
	)

	storage := NewRedisStorage(redisOpts)

	err := storage.Set(entityId, state)

	if err != nil {
		t.Fatalf("Set(%v, \"%v\") вернул ошибку: \"%v\"", entityId, state, err)
	}
	if currState, err := redisClient.Get(context.Background(), toRedisKey(entityId)).Result(); currState != state {
		t.Fatalf("Set(%v, \"%v\") отработал правильно, но текущее значение стейта равно \"%v\" (err: %v)", entityId, state, currState, err)
	}
}

func TestRedisStorage_SetExistEntity(t *testing.T) {
	var (
		prevState = gofakeit.LetterN(10)
		newState  = gofakeit.LetterN(10)
		entityId  = gofakeit.Int64()
	)

	storage := NewRedisStorage(redisOpts)

	err1 := storage.Set(entityId, prevState)
	err2 := storage.Set(entityId, newState)

	if err1 != nil {
		t.Fatalf("Set(%v, \"%v\") вернул ошибку: \"%v\"", entityId, prevState, err1)
	}
	if err2 != nil {
		t.Fatalf("Set(%v, \"%v\") вернул ошибку: \"%v\"", entityId, prevState, err2)
	}
	if currState, err := redisClient.Get(context.Background(), toRedisKey(entityId)).Result(); currState != newState {
		t.Fatalf("Set(%v, \"%v\") отработал правильно, но текущее значение стейта равно \"%v\" (err: %v)", entityId, newState, currState, err)
	}
}

func TestRedisStorage_ClearExistEntity(t *testing.T) {
	var (
		state    = gofakeit.LetterN(10)
		entityId = gofakeit.Int64()
	)

	storage := NewRedisStorage(redisOpts)

	errSet := redisClient.Set(context.Background(), toRedisKey(entityId), state, 0).Err()
	errClear := storage.Clear(entityId)
	currState, errGet := redisClient.Get(context.Background(), toRedisKey(entityId)).Result()

	if errSet != nil {
		t.Fatalf("Не удалось подготовить данные для теста: %v", errSet)
	}
	if errClear != nil {
		t.Fatalf("Clear(%v) вернул ошибку: %v", entityId, errClear)
	}
	if errGet != nil && !errors.Is(errGet, redis.Nil) {
		t.Fatalf("Не удалось получить текущий стейт из Redis: %v", errGet)
	}
	if currState != "" {
		t.Fatalf("Clear(%v) отработал правильно, но текущее значение стейта непустое: %v", entityId, currState)
	}
}

func TestRedisStorage_ClearNotExistEntity(t *testing.T) {
	entityId := gofakeit.Int64()

	storage := NewRedisStorage(redisOpts)

	errClear := storage.Clear(entityId)
	currState, errGet := redisClient.Get(context.Background(), toRedisKey(entityId)).Result()

	if errClear != nil {
		t.Fatalf("Clear(%v) вернул ошибку: %v", entityId, errClear)
	}
	if errGet != nil && !errors.Is(errGet, redis.Nil) {
		t.Fatalf("Не удалось получить текущий стейт из Redis: %v", errGet)
	}
	if currState != "" {
		t.Fatalf("Clear(%v) отработал правильно, но текущее значение стейта непустое: %v", entityId, currState)
	}
}

func TestRedisStorage_CurrentNotExistEntity(t *testing.T) {
	entityId := gofakeit.Int64()

	storage := NewRedisStorage(redisOpts)

	state, err := storage.Current(entityId)

	if err != nil {
		t.Fatalf("Current(%v) вернул ошибку: %v", entityId, err)
	}
	if state != "" {
		t.Fatalf("Current(%v) отработал правильно, но текущее значение стейта непустое: \"%v\"", entityId, state)
	}
}

func TestRedisStorage_CurrentExistEntity(t *testing.T) {
	var (
		entityId = gofakeit.Int64()
		state    = gofakeit.LetterN(10)
	)

	storage := NewRedisStorage(redisOpts)

	errSet := redisClient.Set(context.Background(), toRedisKey(entityId), state, 0).Err()
	currState, errGet := storage.Current(entityId)

	if errSet != nil {
		t.Fatalf("Не удалось подготовить данные для теста: %v", errSet)
	}
	if errGet != nil {
		t.Fatalf("Current(%v) вернул ошибку: %v", entityId, errGet)
	}
	if currState != state {
		t.Fatalf("Current(%v) отработал правильно, но текущее значенпе стейта равно: \"%v\"", entityId, currState)
	}
}
