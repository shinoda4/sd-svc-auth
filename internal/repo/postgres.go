package repo

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repo struct {
	db *sqlx.DB
}

func (r *Repo) Close() {
	err := r.db.Close()
	if err != nil {
		return
	}
}

func NewPostgres(dsn string) (*Repo, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &Repo{db: db}, nil
}
