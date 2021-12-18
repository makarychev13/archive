package state

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
)

func Test_Set_NewEntity_Success(t *testing.T) {
	var (
		state    = gofakeit.LetterN(10)
		entityId = gofakeit.Int64()
	)

	storage := NewMemoryStorage()

	err := storage.Set(entityId, state)

	if err != nil {
		t.Fatalf("Set(%v, \"%v\") вернул ошибку: \"%v\"", entityId, state, err)
	}
	if currState, _ := storage.states[entityId]; currState != state {
		t.Fatalf("Set(%v, \"%v\") отработал правильно, но текущее значение стейта равно \"%v\"", entityId, state, currState)
	}
}

func Test_Set_ExistEntity_Success(t *testing.T) {
	var (
		prevState = gofakeit.LetterN(10)
		newState  = gofakeit.LetterN(10)
		entityId  = gofakeit.Int64()
	)

	storage := NewMemoryStorage()

	err1 := storage.Set(entityId, prevState)
	err2 := storage.Set(entityId, newState)

	if err1 != nil {
		t.Fatalf("Set(%v, \"%v\") вернул ошибку: \"%v\"", entityId, prevState, err1)
	}
	if err2 != nil {
		t.Fatalf("Set(%v, \"%v\") вернул ошибку: \"%v\"", entityId, prevState, err2)
	}
	if currState, _ := storage.states[entityId]; currState != newState {
		t.Fatalf("Set(%v, \"%v\") отработал правильно, но текущее значение стейта равно \"%v\"", entityId, newState, currState)
	}
}

func Test_Clear_ExistEntity_Success(t *testing.T) {
	var (
		state    = gofakeit.LetterN(10)
		entityId = gofakeit.Int64()
	)

	storage := NewMemoryStorage()

	storage.states[entityId] = state
	err := storage.Clear(entityId)

	if err != nil {
		t.Fatalf("Clear(%v) вернул ошибку: %v", entityId, err)
	}
	if currState, ok := storage.states[entityId]; ok {
		t.Fatalf("Clear(%v) отработал правильно, но текущее значение стейта непустое: \"%v\"", entityId, currState)
	}
}

func Test_Clear_NotExistEntity_Success(t *testing.T) {
	entityId := gofakeit.Int64()

	storage := NewMemoryStorage()

	err := storage.Clear(entityId)

	if err != nil {
		t.Fatalf("Clear(%v) вернул ошибку: %v", entityId, err)
	}
	if state, ok := storage.states[entityId]; ok {
		t.Fatalf("Clear(%v) отработал правильно, но текущее значение стейта непустое: \"%v\"", entityId, state)
	}
}

func Test_Current_NotExistEntity_EmptyState(t *testing.T) {
	entityId := gofakeit.Int64()

	storage := NewMemoryStorage()

	state, err := storage.Current(entityId)

	if err != nil {
		t.Fatalf("Current(%v) вернул ошибку: %v", entityId, err)
	}
	if state != "" {
		t.Fatalf("Current(%v) отработал правильно, но текущее значение стейта непустое: \"%v\"", entityId, state)
	}
}

func Test_Current_ExistEntity_Success(t *testing.T) {
	var (
		entityId = gofakeit.Int64()
		state    = gofakeit.LetterN(10)
	)

	storage := NewMemoryStorage()
	storage.states[entityId] = state

	currState, err := storage.Current(entityId)

	if err != nil {
		t.Fatalf("Current(%v) вернул ошибку: %v", entityId, err)
	}
	if err != nil {
		t.Fatalf("Current(%v) отработал правильно, но текущее значенпе стейта равно: \"%v\"", entityId, currState)
	}
}
