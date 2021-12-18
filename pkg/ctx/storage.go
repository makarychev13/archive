package ctx

type Storage interface {
	Get(telegramID int64) (interface{}, error)
	Set(telegramID int64, ctx interface{}) error
}