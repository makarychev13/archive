package repository

import (
	"time"
)

type Days interface {
	Save(telegramID int64, date time.Time) error
	CompleteDay(telegramID int64, end time.Time) error
}

type Tasks interface {
	Save(telegramID int64, name string, date time.Time) (int64, error)
	Complete(taskID int64, date time.Time) (string, error)
	Remove(taskID int64) error
}
