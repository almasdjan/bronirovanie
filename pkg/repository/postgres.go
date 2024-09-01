package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

const (
	reservationsTable = "reservations"
)

func NewPostgresDB(databaseURL string) (*pgx.Conn, error) {

	db, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		logrus.Print(err.Error())
		return nil, err
	}

	return db, nil
}
