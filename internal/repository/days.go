package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	alreadyExistsCode = "23505"
)

var (
	ErrAlreadyExists       = errors.New("не удалось добавить запись из-за уникального констрейта")
	ErrDayAlreadyCompleted = errors.New("день уже завершён")
	ErrAnotherDayStarted   = errors.New("начат уже другой день")
)

type DaysPg struct {
	pool *pgxpool.Pool
}

func NewDaysRepository(p *pgxpool.Pool) *DaysPg {
	return &DaysPg{p}
}

func (r *DaysPg) Save(telegramID int64, date time.Time) error {
	sql :=
		`INSERT INTO "days"
		("telegram_id", "date", "start")
		VALUES($1, $2, $3)`

	_, err := r.pool.Exec(context.Background(), sql, telegramID, date, date)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == alreadyExistsCode {
		return ErrAlreadyExists
	}

	return err
}

func (r *DaysPg) CompleteDay(telegramID int64, end time.Time) error {
	sql :=
		`UPDATE "days"
		 SET "end" = $1
		 FROM (SELECT "id" FROM "days" WHERE "end" IS NULL AND "telegram_id" = $2) AS "cte"
		 WHERE "cte"."id" = "days"."id"
		 RETURNING "days"."id"`

	var dayID int
	err := r.pool.QueryRow(context.Background(), sql, end, telegramID).Scan(&dayID)
	if errors.Is(pgx.ErrNoRows, err) {
		return ErrDayAlreadyCompleted
	}

	return err
}
