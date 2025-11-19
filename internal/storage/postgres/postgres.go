package postgres

import "github.com/jmoiron/sqlx"

type Postgres struct {
	DB *sqlx.DB
}
