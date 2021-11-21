package storage

type InMemory struct {
	states map[int64]string
}

func NewInMemory() *InMemory {
	return &InMemory{map[int64]string{}}
}

func (s *InMemory) Current(id int64) (string, error) {
	return s.states[id], nil
}

func (s *InMemory) Set(id int64, state string) error {
	s.states[id] = state
	return nil
}

func (s *InMemory) Clear(id int64) error {
	delete(s.states, id)
	return nil
}
