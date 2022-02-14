package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/makarychev13/archive/internal/domain"
)

type TasksPg struct {
	pool *pgxpool.Pool
}

func NewPgTasks(pool *pgxpool.Pool) *TasksPg {
	return &TasksPg{pool}
}

func (r *TasksPg) Save(telegramID int64, task domain.Task) (domain.TaskID, error) {
	sql :=
		`INSERT INTO "tasks"
		 ("day_id", "name", "start")
		 SELECT "id", $1, $2
		 FROM "days"
		 WHERE "telegram_id" = $3 AND "end" IS NULL
		 RETURNING "id"`

	var taskID int64
	err := r.pool.QueryRow(context.Background(), sql, task.Name, task.Start, telegramID).Scan(&taskID)

	return taskID, err
}

func (r *TasksPg) Complete(id domain.TaskID, date time.Time) (*domain.Task, error) {
	sql :=
		`UPDATE "tasks"
		 SET "end" = $1
		 WHERE "id" = $2
		 RETURNING "name", "start", "end"`

	var name string
	var start, end time.Time
	err := r.pool.QueryRow(context.Background(), sql, date, id).Scan(&name, &start, &end)

	task := domain.NewFinishedTask(name, start, end)

	return &task, err
}

func (r *TasksPg) Remove(id domain.TaskID) error {
	sql :=
		`DELETE FROM "tasks"
		 WHERE "id" = $1`

	_, err := r.pool.Exec(context.Background(), sql, id)

	return err
}
