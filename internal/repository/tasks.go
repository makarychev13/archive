package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PgTasksRepository struct {
	pool *pgxpool.Pool
}

func NewTasksRepository(pool *pgxpool.Pool) *PgTasksRepository {
	return &PgTasksRepository{pool}
}

func (r *PgTasksRepository) Save(telegramID int64, name string, date time.Time) (int64, error) {
	sql :=
		`INSERT INTO "tasks"
		 ("day_id", "name", "start")
		 SELECT "id", $1, $2
		 FROM "days"
		 WHERE "telegram_id" = $3 AND "end" IS NULL
		 RETURNING "id"`

	var taskID int64
	err := r.pool.QueryRow(context.Background(), sql, name, date, telegramID).Scan(&taskID)

	return taskID, err
}

func (r *PgTasksRepository) Complete(taskID int64, date time.Time) (string, error) {
	sql :=
		`UPDATE "tasks"
		 SET "end" = $1
		 WHERE "id" = $2
		 RETURNING "name"`

	var name string
	err := r.pool.QueryRow(context.Background(), sql, date, taskID).Scan(&name)

	return name, err
}

func (r *PgTasksRepository) Remove(taskID int64) error {
	sql :=
		`DELETE FROM "tasks"
		 WHERE "id" = $1`

	_, err := r.pool.Exec(context.Background(), sql, taskID)

	return err
}
