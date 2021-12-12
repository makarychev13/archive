package repository

import "time"

type DayID = int64

type Days interface {
	New(telegramID int64, date time.Time) error
	ActiveDay(telegramID int64) (*DayID, error)
}

type Tasks interface {
	New(name string, date time.Time) error
}