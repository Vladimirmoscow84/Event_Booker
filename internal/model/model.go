package model

import "time"

// Бронируемое событие
type Event struct {
	ID              int           `json:"id" db:"id"`
	Title           string        `json:"title" db:"title"`
	Date            time.Time     `json:"date" db:"date"`
	TotalSeats      int           `json:"total_seats" db:"total_seats"`
	AvailableSeats  int           `json:"available_seats" db:"available_seats"`
	BookingTTL      time.Duration `json:"booking_ttl" db:"booking_ttl"`
	RequiresPayment bool          `json:"requires_payment" db:"requires_payment"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
}

// Статус брони
type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCanceled  BookingStatus = "canceled"
)

// Бронь
type Booking struct {
	ID        int           `json:"id" db:"id"`
	EventID   int           `json:"event_id" db:"event_id"`
	UserID    int           `json:"user_id" db:"user_id"`
	Status    BookingStatus `json:"status" db:"status"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	ExpiresAt time.Time     `json:"expires_at" db:"expires_at"`
}

// Данные пользователя
type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
