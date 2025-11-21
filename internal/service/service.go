package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Vladimirmoscow84/Event_Booker/internal/model"
	"github.com/Vladimirmoscow84/Event_Booker/internal/notifier"
	"github.com/Vladimirmoscow84/Event_Booker/internal/storage"
)

type eventService interface {
	CreateEvent(ctx context.Context, event *model.Event) (int, error)
	GetEvent(ctx context.Context, id int) (*model.Event, error)
	ListEvents(ctx context.Context) ([]*model.Event, error)
}

type bookingService interface {
	CreateBooking(ctx context.Context, eventID, userID int) (int, error)
	ConfirmBooking(ctx context.Context, bookingID int) error
	CancelBooking(ctx context.Context, bookingID int) error
	GetBooking(ctx context.Context, id int) (*model.Booking, error)
	GetExpiredBookings(ctx context.Context) ([]*model.Booking, error)
}

type userService interface {
	CreateUser(ctx context.Context, email string) (int, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type Service struct {
	storage  *storage.Storage
	notifier *notifier.Client
	ttl      time.Duration
}

// compile-time assertions
var _ eventService = (*Service)(nil)
var _ bookingService = (*Service)(nil)
var _ userService = (*Service)(nil)

func New(storage *storage.Storage, notifier *notifier.Client, ttl time.Duration) *Service {
	return &Service{
		storage:  storage,
		notifier: notifier,
		ttl:      ttl,
	}
}

// Реализация интерфейса eventService

// CreateEvent создает событие и возвращает его ID
func (s *Service) CreateEvent(ctx context.Context, event *model.Event) (int, error) {
	log.Printf("[service] CreateEvent: %s", event.Title)
	id, err := s.storage.CreateEvent(ctx, event)
	if err != nil {
		return 0, err
	}

	// email уведомление
	if s.notifier != nil {
		go s.notifier.Send(ctx,
			"Новое мероприятие создано",
			fmt.Sprintf("Создано мероприятие: %s (ID=%d)", event.Title, id),
		)
	}
	return id, nil

}

// GetEvent возвращает событие по ID
func (s *Service) GetEvent(ctx context.Context, id int) (*model.Event, error) {
	log.Printf("[service] GetEvent: id=%d", id)
	return s.storage.GetEvent(ctx, id)
}

// ListEvents возвращает список всех событий
func (s *Service) ListEvents(ctx context.Context) ([]*model.Event, error) {
	log.Println("[service] ListEvents")
	return s.storage.ListEvents(ctx)
}

//Реализация интерфейса bookongService

// CreateBooking создает бронь для пользователя на событие
func (s *Service) CreateBooking(ctx context.Context, eventID, userID int) (int, error) {
	log.Printf("[service] CreateBooking: event=%d user=%d", eventID, userID)

	event, err := s.storage.GetEvent(ctx, eventID)
	if err != nil {
		return 0, err
	}

	if event.AvailableSeats <= 0 {
		return 0, fmt.Errorf("no seats available")
	}

	booking := &model.Booking{
		EventID:   eventID,
		UserID:    userID,
		Status:    model.BookingStatusPending,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(event.BookingTTL), // время жизни брони
	}

	id, err := s.storage.CreateBooking(ctx, booking)
	if err != nil {
		return 0, err
	}

	event.AvailableSeats--
	err = s.storage.UpdateEvent(ctx, event)
	if err != nil {
		log.Printf("[service] warning: failed to update available seats for event %d: %v", eventID, err)
	}

	// email уведомление
	if s.notifier != nil {
		go s.notifier.Send(ctx,
			"Бронь создана",
			fmt.Sprintf("Создана бронь #%d на мероприятие: %s", id, event.Title),
		)
	}

	return id, nil
}

// ConfirmBooking подтверждает бронь при оплате
func (s *Service) ConfirmBooking(ctx context.Context, bookingID int) error {
	log.Printf("[service] ConfirmBooking: id=%d", bookingID)
	b, err := s.storage.GetBooking(ctx, bookingID)
	if err != nil {
		return err
	}
	if b.Status != model.BookingStatusPending {
		return fmt.Errorf("booking was payed or canceled")
	}

	err = s.storage.UpdateBookingStatus(ctx, bookingID, model.BookingStatusConfirmed)
	if err != nil {
		return err
	}

	// email уведомление
	if s.notifier != nil {
		go s.notifier.Send(ctx,
			"Бронь подтверждена",
			fmt.Sprintf("Ваша бронь #%d успешно подтверждена!", bookingID),
		)
	}

	return nil
}

// CancelBooking отменяет бронь
func (s *Service) CancelBooking(ctx context.Context, bookingID int) error {
	log.Printf("[service] CancelBooking: id=%d", bookingID)
	b, err := s.storage.GetBooking(ctx, bookingID)
	if err != nil {
		return err
	}

	if b.Status == model.BookingStatusCanceled {
		return nil
	}

	err = s.storage.UpdateBookingStatus(ctx, bookingID, model.BookingStatusCanceled)
	if err != nil {
		return err
	}

	// в случае отмены идет возврат места
	event, err := s.storage.GetEvent(ctx, b.EventID)
	if err == nil {
		event.AvailableSeats++
		_ = s.storage.UpdateEvent(ctx, event)
	}

	// email уведомление
	if s.notifier != nil {
		go s.notifier.Send(ctx,
			"Бронь отменена",
			fmt.Sprintf("Ваша бронь #%d была отменена.", bookingID),
		)
	}

	return nil
}

// GetBooking возвращает бронь по ID
func (s *Service) GetBooking(ctx context.Context, id int) (*model.Booking, error) {
	log.Printf("[service] GetBooking: id=%d", id)
	return s.storage.GetBooking(ctx, id)
}

// GetExpiredBookings возвращает все просроченные брони
func (s *Service) GetExpiredBookings(ctx context.Context) ([]*model.Booking, error) {
	log.Println("[service] GetExpiredBookings")
	return s.storage.GetExpiredBookings(ctx)
}

//Реализация интерфейса userService

// CreateUser создает нового пользователя
func (s *Service) CreateUser(ctx context.Context, email string) (int, error) {
	log.Printf("[service] CreateUser: %s", email)
	user := &model.User{
		Email:     email,
		CreatedAt: time.Now(),
	}
	return s.storage.CreateUser(ctx, user)
}

// GetUserByID возвращает пользователя по ID
func (s *Service) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	log.Printf("[service] GetUserByID: id=%d", id)
	return s.storage.GetUserByID(ctx, id)
}

// GetUserByEmail возвращает пользователя по email
func (s *Service) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	log.Printf("[service] GetUserByEmail: %s", email)
	return s.storage.GetUserByEmail(ctx, email)
}
