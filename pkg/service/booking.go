package service

import (
	"bronirovanie/models"
	"bronirovanie/pkg/repository"
	"context"
	"fmt"
	"strconv"
	"sync"
)

type BookingService struct {
	repo   repository.Booking
	muMap  map[string]*sync.Mutex
	muLock sync.Mutex
}

func NewBookingService(repo repository.Booking) *BookingService {
	return &BookingService{
		repo:  repo,
		muMap: make(map[string]*sync.Mutex),
	}
}

func (s *BookingService) getMutex(roomID string) *sync.Mutex {
	s.muLock.Lock()
	defer s.muLock.Unlock()

	if _, exists := s.muMap[roomID]; !exists {
		s.muMap[roomID] = &sync.Mutex{}
	}

	return s.muMap[roomID]
}

func (s *BookingService) Create(c context.Context, reservation *models.CreateReservation) error {
	if reservation.StartTime.After(reservation.EndTime) {
		return fmt.Errorf("start time must be before end time")
	}
	mu := s.getMutex(strconv.Itoa(reservation.RoomId))
	mu.Lock()
	defer mu.Unlock()

	existingReservations, err := s.repo.GetAll(c, reservation.RoomId)
	if err != nil {
		return err
	}

	for _, existing := range existingReservations {
		if (reservation.StartTime.Before(existing.EndTime) && reservation.StartTime.After(existing.StartTime)) ||
			(reservation.EndTime.After(existing.StartTime) && reservation.EndTime.Before(existing.EndTime)) ||
			(reservation.StartTime.Equal(existing.StartTime) || reservation.EndTime.Equal(existing.EndTime)) {
			return fmt.Errorf("Conflict")
		}
	}

	err = s.repo.Create(c, reservation)
	if err != nil {
		return err
	}
	return nil
}

func (s *BookingService) GetAll(c context.Context, roomId int) ([]models.Reservation, error) {
	reservs, err := s.repo.GetAll(c, roomId)
	if err != nil {
		return reservs, err
	}
	return reservs, nil
}
