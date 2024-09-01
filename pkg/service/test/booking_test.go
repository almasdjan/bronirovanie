package test

import (
	"bronirovanie/models"
	"bronirovanie/pkg/repository"
	"bronirovanie/pkg/service"
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	db             *pgx.Conn
	repo           *repository.Repository
	bookingService *service.BookingService
)

func setupTestDB() (*pgx.Conn, error) {
	db, err := repository.NewPostgresDB(viper.GetString("environment.TEST_DATABASE_URL"))
	if err != nil {
		logrus.Print(viper.GetString("environment.TEST_DATABASE_URL"))
		return nil, err
	}
	return db, nil
}

func createTables() error {
	stmt := `
        CREATE TABLE IF NOT EXISTS reservations (
            id SERIAL PRIMARY KEY,
            room_id INTEGER NOT NULL,
            start_time TIMESTAMP NOT NULL,
            end_time TIMESTAMP NOT NULL
        );
    `
	_, err := db.Exec(context.Background(), stmt)
	return err
}

func TestMain(m *testing.M) {
	var err error
	db, err := setupTestDB()
	if err != nil {
		logrus.Fatalf("Failed to set up test DB: %v", err)
	}
	defer db.Close(context.Background())

	repo = repository.NewRepository(db)
	bookingService = service.NewBookingService(repo)

	err = createTables()
	if err != nil {
		logrus.Fatalf("Failed to create tables: %v", err)
	}

	code := m.Run()

	if err := clearTestDatabase(); err != nil {
		logrus.Fatalf("Failed to clear test database: %v", err)
	}

	os.Exit(code)
}

func TestCreateReservation_Success(t *testing.T) {

	reservation := &models.CreateReservation{
		RoomId:    1,
		StartTime: time.Now().Add(1 * time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
	}

	err := bookingService.Create(context.Background(), reservation)
	if err != nil && err.Error() != "Conflict" {
		t.Fatalf("unexpected error: %v", err)
	}

	reservations, err := bookingService.GetAll(context.Background(), reservation.RoomId)
	if err != nil {
		t.Fatalf("failed to fetch reservations: %v", err)
	}

	for _, r := range reservations {
		t.Logf("Reservation ID: %d, Room ID: %d, Start Time: %v, End Time: %v", r.Id, r.RoomId, r.StartTime, r.EndTime)
	}
}

func TestCreateReservationWithConflict_Success(t *testing.T) {

	reservation := &models.CreateReservation{
		RoomId:    1,
		StartTime: time.Now().Add(3 * time.Hour),
		EndTime:   time.Now().Add(4 * time.Hour),
	}
	reservation2 := &models.CreateReservation{
		RoomId:    1,
		StartTime: time.Now().Add(3 * time.Hour).Add(59 * time.Minute),
		EndTime:   time.Now().Add(4 * time.Hour).Add(30 * time.Minute),
	}

	err := bookingService.Create(context.Background(), reservation)
	if err != nil && err.Error() != "Conflict" {
		t.Fatalf("unexpected error: %v", err)
	}

	err = bookingService.Create(context.Background(), reservation2)
	if err != nil && err.Error() != "Conflict" {
		t.Fatalf("unexpected error: %v", err)
	}

	reservations, err := bookingService.GetAll(context.Background(), reservation.RoomId)
	if err != nil {
		t.Fatalf("failed to fetch reservations: %v", err)
	}

	for _, r := range reservations {
		t.Logf("Reservation ID: %d, Room ID: %d, Start Time: %v, End Time: %v", r.Id, r.RoomId, r.StartTime, r.EndTime)
	}
}

func TestCreateReservation_ConcurrentRequests(t *testing.T) {

	var wg sync.WaitGroup
	errors := make(chan error, 2)

	reservation := &models.CreateReservation{
		RoomId:    1,
		StartTime: time.Now().Add(5 * time.Hour),
		EndTime:   time.Now().Add(6 * time.Hour),
	}

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := bookingService.Create(context.Background(), reservation)
			errors <- err
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil && err.Error() != "Conflict" {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	reservations, err := bookingService.GetAll(context.Background(), reservation.RoomId)
	if err != nil {
		t.Fatalf("failed to fetch reservations: %v", err)
	}

	for _, r := range reservations {
		t.Logf("Reservation ID: %d, Room ID: %d, Start Time: %v, End Time: %v", r.Id, r.RoomId, r.StartTime, r.EndTime)
	}
}

func clearTestDatabase() error {
	_, err := db.Exec(context.Background(), `TRUNCATE TABLE reservations`)
	return err
}
