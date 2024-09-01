package repository

import (
	"bronirovanie/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type BookingPostgres struct {
	db *pgx.Conn
}

func NewBookingPostgres(db *pgx.Conn) *BookingPostgres {
	return &BookingPostgres{db: db}
}

func (r *BookingPostgres) Create(c context.Context, reservation *models.CreateReservation) error {
	/*
		var count int

			//get count of conflicting reservations
			queryGet := fmt.Sprintf(`SELECT COUNT(*) FROM %s
									WHERE room_id = $1
									AND ((start_time < $3 AND end_time > $2)
									OR (start_time >= $2 AND start_time < $3)`, reservationsTable)
			row := r.db.QueryRow(c, queryGet, reservation.RoomId, reservation.StartTime, reservation.EndTime)

			if err := row.Scan(&count); err != nil {
				return err
			}

			if count > 0 {
				return fmt.Errorf("Conflict")
			}
	*/

	//insertion
	queryInsert := fmt.Sprintf("INSERT INTO %s (room_id, start_time, end_time) values($1, $2, $3)", reservationsTable)
	_, err := r.db.Exec(c, queryInsert, reservation.RoomId, reservation.StartTime, reservation.EndTime)
	if err != nil {
		return err
	}
	return nil
}

func (r *BookingPostgres) GetAll(c context.Context, roomId int) ([]models.Reservation, error) {
	var reservs []models.Reservation
	query := fmt.Sprintf("select * from %s where room_id = $1 order by start_time", reservationsTable)
	rows, err := r.db.Query(c, query, roomId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var reserv models.Reservation
		err := rows.Scan(&reserv.Id, &reserv.RoomId, &reserv.StartTime, &reserv.EndTime)
		if err != nil {
			return reservs, err
		}
		reservs = append(reservs, reserv)
	}

	return reservs, nil
}
