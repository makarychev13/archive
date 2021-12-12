package storage

import "sync"

type InMemory struct {
	mu sync.Mutex
	states map[int64]string
}

func NewInMemory() *InMemory {
	return &InMemory{states: map[int64]string{}}
}

func (s *InMemory) Current(id int64) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	res := s.states[id]

	return res, nil
}

func (s *InMemory) Set(id int64, state string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.states[id] = state

	return nil
}

func (s *InMemory) Clear(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.states, id)

	return nil
}
