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
)

type PgDaysRepository struct {
	pool *pgxpool.Pool
}

func NewDaysRepository(p *pgxpool.Pool) *PgDaysRepository {
	return &PgDaysRepository{p}
}

func (r *PgDaysRepository) Save(telegramID int64, date time.Time) error {
	sql :=
		`INSERT INTO "days"
		("telegram_id", "date", "start")
		VALUES($1, $2, $3)`

	_, err := r.pool.Exec(context.Background(), sql, telegramID, date, date)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == alreadyExistsCode {
			return ErrAlreadyExists
		}

		return err
	}

	return nil
}

func (r *PgDaysRepository) CompleteDay(telegramID int64, date, end time.Time) error {
	sql :=
		`UPDATE "days"
		 SET "end" = $1
		 FROM (SELECT "id" FROM "days" WHERE "date" = $2 AND "end" IS NULL AND "telegram_id" = $3) AS "cte"
		 WHERE "cte"."id" = "days"."id"
		 returning days.id`

	var dayID int
	err := r.pool.QueryRow(context.Background(), sql, end, date, telegramID).Scan(&dayID)
	if err == pgx.ErrNoRows {
		return ErrDayAlreadyCompleted
	}
	if err != nil {
		return err
	}

	return nil
}
