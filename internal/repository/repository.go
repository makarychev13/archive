package repository

import (
	"time"

	"github.com/makarychev13/archive/internal/domain"
)

type Days interface {
	Save(telegramID int64, date time.Time) error
	CompleteDay(telegramID int64, end time.Time) error
}

type Tasks interface {
	Save(telegramID int64, task domain.Task) (domain.TaskID, error)
	Complete(id domain.TaskID, date time.Time) (*domain.Task, error)
	Remove(id domain.TaskID) error
}
