package storage

import (
	"context"

	"github.com/Vladimirmoscow84/Event_Booker/internal/model"
	"github.com/Vladimirmoscow84/Event_Booker/internal/storage/postgres"
)

type eventStorage interface {
	CreateEvent(ctx context.Context, event *model.Event) (int, error)
	GetEvent(ctx context.Context, id int) (*model.Event, error)
	ListEvents(ctx context.Context) ([]*model.Event, error)
	UpdateEvent(ctx context.Context, event *model.Event) error
}

type bookingStorage interface {
	CreateBooking(ctx context.Context, booking *model.Booking) (int, error)
	GetBooking(ctx context.Context, id int) (*model.Booking, error)
	GetExpiredBookings(ctx context.Context) ([]*model.Booking, error)
	UpdateBookingStatus(ctx context.Context, id int, status model.BookingStatus) error
	DeleteBooking(ctx context.Context, id int) error
}

type userStorage interface {
	CreateUser(ctx context.Context, user *model.User) (int, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type Storage struct {
	*postgres.Postgres
}

func New(postgres *postgres.Postgres) *Storage {
	return &Storage{
		Postgres: postgres,
	}
}
