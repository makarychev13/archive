package state

import "sync"

//MemoryStorage хранит текущий стейт пользователя в памяти. При перезапуске бота вся информация удаляется.
type MemoryStorage struct {
	mu sync.Mutex
	states map[int64]string
}

//NewMemoryStorage создаёт не singleton-объект MemoryStorage. Его нужно переиспользовать во всех обработчиках стейтов.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{states: map[int64]string{}}
}

func (s *MemoryStorage) Current(id int64) (Name, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	res := s.states[id]

	return res, nil
}

func (s *MemoryStorage) Set(id int64, state Name) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.states[id] = state

	return nil
}

func (s *MemoryStorage) Clear(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.states, id)

	return nil
}
