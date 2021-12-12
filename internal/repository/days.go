package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	alreadyExistsCode = "23505"
)

var (
	ErrAlreadyExists = errors.New("Не удалось добавить запись из-за уникального констрейта")
)

type PgDaysRepository struct {
	pool *pgxpool.Pool
}

func NewDaysRepository(p *pgxpool.Pool) PgDaysRepository {
	return PgDaysRepository{p}
}

func (r PgDaysRepository) New(telegramID int64, date time.Time) error {
	sql :=
		`INSERT INTO "days"
		("telegram_id", "date", "start")
		VALUES($1, $2, $3)`

	_, err := r.pool.Exec(context.Background(), sql, telegramID, date)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == alreadyExistsCode {
			return ErrAlreadyExists
		}

		return err
	}

	return nil
}

func (r PgDaysRepository) ActiveDay(telegramID int64) (*DayID, error) {
	sql := `SELECT FROM "days" WHERE "end" IS NULL AND "telegram_id" = $1`

	var dayID int64
	if err := r.pool.QueryRow(context.Background(), sql, telegramID).Scan(&dayID); err != nil {
		return nil, err
	}

	return &dayID, nil
}