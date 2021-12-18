package state

//Storage описывает контракт для хранения и изменения текущего стейта пользователя.
type Storage interface {
	Current(id int64) (Name, error)
	Set(id int64, name Name) error
	Clear(id int64) error
}
