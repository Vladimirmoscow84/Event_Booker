package postgres

import (
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	DB *sqlx.DB
}

func New(databaseURI string) (*Postgres, error) {
	db, err := sqlx.Connect("pgx", databaseURI)
	if err != nil {
		return nil, fmt.Errorf("[postgres] failed to connect to DB: %w", err)
	}

	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(30)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("[postgres] ping failed: %w ", err)
	}

	log.Println("[postgres] successfull connect to DB")
	return &Postgres{
		DB: db,
	}, nil
}

func (p *Postgres) Close() error {
	err := p.DB.Close()
	if err != nil {
		log.Printf("[postgres] failed close DB connection: %v", err)
		return err
	}
	log.Println("[postgres] connection closed")
	return nil
}
