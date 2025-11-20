package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Vladimirmoscow84/Event_Booker/internal/model"
)

// CreateEvent создает новую запись в таблице  events
func (p *Postgres) CreateEvent(ctx context.Context, event *model.Event) (int, error) {
	query := `
			INSERT INTO events
				(title, date, tatal_seats, available_seats, booking_ttl, requires_payment,created_at)
			VALUES
				($1,$2,$3,$4,$5,$6,NOW())	
			RETURNING id;
		`
	row := p.DB.QueryRowxContext(ctx, query, event.Title, event.Date, event.TotalSeats, event.AvailableSeats, int(event.BookingTTL.Minutes()), event.RequiresPayment, event.CreatedAt)

	var id int
	err := row.Scan(&id)
	if err != nil {
		log.Printf("[postgres] erorr adding event to DB: %v", err)
		return 0, fmt.Errorf("[postgres] erorr adding event to DB: %w", err)
	}
	event.ID = id
	return id, nil
}

// GetEvent получает запись из таблицы events по id
func (p *Postgres) GetEvent(ctx context.Context, id int) (*model.Event, error) {
	var event model.Event

	query := `
			SELECT 
				id, title, date, tatal_seats, available_seats, booking_ttl, requires_payment,created_at
			FROM events
			WHERE id=$1
	`
	err := p.DB.GetContext(ctx, &event, query, id)
	if err != nil {
		log.Printf("[postgres] error getting event from events: %v", err)
		return nil, fmt.Errorf("[postgres] error getting event from events: %w", err)
	}
	event.BookingTTL = event.BookingTTL * time.Minute
	return &event, nil
}
