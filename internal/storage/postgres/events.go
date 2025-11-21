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
				(title, date, total_seats, available_seats, booking_ttl, requires_payment,created_at)
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

// ListEvents получает все записи из таблицы events
func (p *Postgres) ListEvents(ctx context.Context) ([]*model.Event, error) {
	var events []*model.Event

	query := `
		SELECT 
			id, title, date, total_seats, available_seats, booking_ttl, requires_payment, created_at
		FROM events
		ORDER BY date ASC
	`
	err := p.DB.SelectContext(ctx, &events, query)
	if err != nil {
		log.Printf("[postgres] error listing events: %v", err)
		return nil, fmt.Errorf("[postgres] error listing events: %w", err)
	}

	for _, event := range events {
		event.BookingTTL = event.BookingTTL * time.Minute
	}

	return events, nil
}

// UpdateEvent обновляет запись в таблице events
func (p *Postgres) UpdateEvent(ctx context.Context, event *model.Event) error {
	query := `
			UPDATE events
			SET 
				title=$1, date=$2, total_seats=$3, available_seats=$4,booking_ttl=$5, requires_payment=$6
			WHERE id=$7
	`
	res, err := p.DB.ExecContext(ctx, query, event.Title, event.Date, event.TotalSeats, event.AvailableSeats, int(event.BookingTTL.Minutes()), event.RequiresPayment, event.ID)

	if err != nil {
		log.Printf("[postgres] error updating event: %v", err)
		return fmt.Errorf("[postgres] error updating event: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		log.Printf("[postgres] update event: no rows affected for id=%d", event.ID)
	}
	return nil
}
