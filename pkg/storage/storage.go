package storage

type Storage interface {
	Current(id int64) (string, error)
	Set(id int64, state string) error
	Clear(id int64) error
}
