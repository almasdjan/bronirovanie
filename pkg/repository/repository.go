package repository

import (
	"bronirovanie/models"
	"context"

	"github.com/jackc/pgx/v5"
)

type Booking interface {
	Create(c context.Context, rezervation *models.CreateReservation) error
	GetAll(c context.Context, roomId int) ([]models.Reservation, error)
}

type Repository struct {
	Booking
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{
		Booking: NewBookingPostgres(db),
	}
}
