package models

import "time"

type Reservation struct {
	Id        int       `json:"id" db:"id"`
	RoomId    int       `json:"room_id" db:"room_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}

type CreateReservation struct {
	RoomId    int       `json:"room_id" db:"room_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}
