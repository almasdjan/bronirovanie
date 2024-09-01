package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func CreateTable(db *pgx.Conn) error {
	stmt := `
		CREATE TABLE IF NOT EXISTS reservations (
			id SERIAL PRIMARY KEY ,
			room_id INTEGER not null,
			start_time TIMESTAMP NOT NULL,
    		end_time TIMESTAMP NOT NULL
		);


    `

	_, err := db.Exec(context.Background(), stmt)
	if err != nil {
		return err
	}

	return nil
}
