package service

import (
	"context"
	"log"
	"time"
)

// ProcessExpiredBookings получает все просроченные брони и отменяет их
// и освобождает место в событии
func (s *Service) ProcessExpiredBookings(ctx context.Context) {
	expired, err := s.GetExpiredBookings(ctx)
	if err != nil {
		log.Printf("[worker] error fetching expired bookings: %v", err)
		return
	}

	if len(expired) == 0 {
		log.Println("[worker] no expired bookings found")
		return
	}

	for _, b := range expired {
		log.Printf("[worker] canceling expired booking %d for event %d", b.ID, b.EventID)
		err := s.CancelBooking(ctx, b.ID)
		if err != nil {
			log.Printf("[worker] error canceling booking %d: %v", b.ID, err)
			continue
		}

		event, err := s.GetEvent(ctx, b.EventID)
		if err != nil {
			log.Printf("[worker] error fetching event %d: %v", b.EventID, err)
			continue
		}

		event.AvailableSeats++
		err = s.storage.UpdateEvent(ctx, event)
		if err != nil {
			log.Printf("[worker] error updating available seats for event %d: %v", b.EventID, err)
		}
	}
}

// StartBookingWorker запускает фоновый воркер, который
// регулярно проверяет просроченные брони и отменяет их.

func (s *Service) StartBookingWorker(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("[worker] stopping booking worker")
				return
			case <-ticker.C:
				s.ProcessExpiredBookings(ctx)
			}
		}
	}()
}
