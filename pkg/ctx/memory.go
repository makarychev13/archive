package ctx

import "sync"

type MemoryStorage struct {
	mu   sync.Mutex
	data map[int64]interface{}
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{data: map[int64]interface{}{}}
}

func (m *MemoryStorage) Get(telegramID int64) (interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.data[telegramID], nil
}

func (m *MemoryStorage) Set(telegramID int64, data interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[telegramID] = data

	return nil
}
