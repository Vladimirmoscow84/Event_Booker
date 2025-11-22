package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Vladimirmoscow84/Event_Booker/internal/model"
)

// CreateUser создает запись в таблице Users
func (p *Postgres) CreateUser(ctx context.Context, user *model.User) (int, error) {
	// Проверка, существует ли пользователь с таким email
	existing, err := p.GetUserByEmail(ctx, user.Email)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("[postgres] error checking existing user: %v", err)
		return 0, fmt.Errorf("[postgres] error checking existing user: %w", err)
	}
	if existing != nil {
		return existing.ID, fmt.Errorf("user with email %s already exists", user.Email)
	}
	query := `
			INSERT INTO users
				(email, created_at)
			VALUES
				($1, NOW())
			RETURNING id;
	`
	row := p.DB.QueryRowContext(ctx, query, user.Email)

	var id int
	err = row.Scan(&id)
	if err != nil {
		log.Printf("[postgres] error creating user: %v", err)
		return 0, fmt.Errorf("[postgres] error creating user: %w", err)
	}
	user.ID = id
	return id, nil
}

// GetUserByID возвращает пользователя по id
func (p *Postgres) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User

	query := `
			SELECT id, email, created_at
			FROM users
			WHERE id=$1
	`
	err := p.DB.GetContext(ctx, &user, query, id)
	if err != nil {
		log.Printf("[postgres] error getting user by ID: %v", err)
		return nil, fmt.Errorf("[postgres] error getting user by ID: %w", err)
	}
	return &user, nil
}

// GetUserByEmail возвращает юзера по email
func (p *Postgres) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
	SELECT 
		id, email, created_at 
	FROM users 
	WHERE email=$1
	`
	row := p.DB.QueryRowContext(ctx, query, email)

	user := &model.User{}
	err := row.Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("[postgres] error getting user by email: %v", err)
		return nil, fmt.Errorf("[postgres] error getting user by email: %w", err)
	}
	return user, nil
}
