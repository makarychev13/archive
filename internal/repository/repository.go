package repository

import (
	"time"
)

type Days interface {
	Save(telegramID int64, date time.Time) error
	CompleteDay(telegramID int64, date, end time.Time) error
}

type Tasks interface {
	New(name string, date time.Time) error
}
