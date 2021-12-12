package repository

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgTasksRepository struct {
	pool *pgxpool.Pool
}

func NewTasksRepository(pool *pgxpool.Pool) PgTasksRepository {
	return PgTasksRepository{pool}
}

//func (r PgTasksRepository) New(name string, date time.Time) error {
//
//}