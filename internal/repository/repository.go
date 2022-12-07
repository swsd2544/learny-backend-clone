package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repositories struct {
	Users UserRepository
}

func New(db *pgxpool.Pool) Repositories {
	return Repositories{
		Users: UserRepository{db: db},
	}
}
