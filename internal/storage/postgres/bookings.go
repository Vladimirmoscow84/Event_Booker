package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Vladimirmoscow84/Event_Booker/internal/model"
)

// CreateBooking создает запись в таблице bookings
func (p *Postgres) CreateBooking(ctx context.Context, booking *model.Booking) (int, error) {
	query := `
			INSERT INTO bookings
				(event_id, user_id, status, created_at, expires_at)
			VALUES
				($1, $2, $3, NOW(), $4)
			RETURNING id;
	`
	row := p.DB.QueryRowContext(ctx, query, booking.EventID, booking.UserID, booking.Status, booking.ExpiresAt)

	var id int
	err := row.Scan(&id)
	if err != nil {
		log.Printf("[postgres] error creating booking: %v", err)
		return 0, fmt.Errorf("[postgres] error creating booking: %w", err)
	}
	booking.ID = id
	return id, nil
}

// GetBooking получает запись из таблицы bookings по id
func (p *Postgres) GetBooking(ctx context.Context, id int) (*model.Booking, error) {
	query := `
		SELECT id, event_id, user_id, status, created_at, expires_at
		FROM bookings
		WHERE id=$1
	`
	row := p.DB.QueryRowContext(ctx, query, id)

	booking := &model.Booking{}
	err := row.Scan(
		&booking.ID,
		&booking.EventID,
		&booking.UserID,
		&booking.Status,
		&booking.CreatedAt,
		&booking.ExpiresAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("[postgres] error getting booking from bookings: %v", err)
		return nil, fmt.Errorf("[postgres] error getting booking from bookings: %w", err)
	}

	return booking, nil
}

// GetExpiredBookings  возвращает просроченные бронирования из таблицы bookings (срок брони вышел)
func (p *Postgres) GetExpiredBookings(ctx context.Context) ([]*model.Booking, error) {
	var bookings []*model.Booking

	query := `
			SELECT 
				id, event_id, user_id, status, created_at, expires_at
			FROM bookings
			WHERE status='pending' AND expires_at <= NOW()
	`
	err := p.DB.SelectContext(ctx, &bookings, query)
	if err != nil {
		log.Printf("[postgres] error fetching expired bookings: %v", err)
		return nil, fmt.Errorf("[postgres] error fetching expired bookings: %w", err)
	}
	return bookings, nil
}

// UpdateStatusBooking менят статус брони (pending, confimed,  canceled)
func (p *Postgres) UpdateBookingStatus(ctx context.Context, id int, status model.BookingStatus) error {

	query := `
		UPDATE bookings
		SET status=$1
		WHERE id=$2
		`
	res, err := p.DB.ExecContext(ctx, query, status, id)
	if err != nil {
		log.Printf("[postgres] error updating booking status: %v", err)
		return fmt.Errorf("[postgres] error updating booking status: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		log.Printf("[postgres] update booking status: no rows affected for id=%d", id)

	}

	return nil
}

// DeletBooking удаляет бронь
func (p *Postgres) DeleteBooking(ctx context.Context, id int) error {
	query := `
			DELETE 
			FROM bookings
			WHERE id=$1
	`
	res, err := p.DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("[postgres] error deleting booking: %v", err)
		return fmt.Errorf("[postgres] error deleting booking: %w", err)

	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		log.Printf("[postgres] delete booking: no rows affected for id=%d", id)

	}
	return nil
}
