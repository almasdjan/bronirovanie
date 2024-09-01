package service

import (
	"bronirovanie/models"
	"bronirovanie/pkg/repository"
	"context"
)

type Booking interface {
	Create(c context.Context, reservation *models.CreateReservation) error
	GetAll(c context.Context, roomId int) ([]models.Reservation, error)
}

type Service struct {
	Booking
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Booking: NewBookingService(repos.Booking),
	}
}
