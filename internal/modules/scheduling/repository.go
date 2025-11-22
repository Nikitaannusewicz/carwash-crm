package scheduling

import (
	"context"
	"database/sql"
	"time"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateShift(ctx context.Context, shift *Shift) error {
	query := `
		INSERT INTO shifts (user_id, location_id, start_time, end_time)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query, shift.UserID, shift.LocationID, shift.StartTime, shift.EndTime).Scan(&shift.ID)
}

func (r *PostgresRepository) CheckOverLap(ctx context.Context, userID int64, start, end time.Time) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM shifts
			WHERE user_id = $1
			AND start_time < $3
			AND end_time > $2
		)
	`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID, start, end).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) CreateBooking(ctx context.Context, booking *Booking) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `
		INSERT INTO bookings (location_id, bay_id, service_id, customer_id, start_time, end_time, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		Returning id
	`
	err = tx.QueryRowContext(
		ctx, query,
		booking.LocationID, booking.BayID, booking.ServiceID, booking.CustomerID,
		booking.StartTime, booking.EndTime, booking.Status).Scan(&booking.ID)

	if err != nil {
		return err
	}

	workerQuery := `INSERT INTO booking_workers (booking_id, user_id) VALUES ($1, $2)`

	stmt, err := tx.PrepareContext(ctx, workerQuery)
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, workerID := range booking.AssignedWorkerIDs {
		if _, err := stmt.ExecContext(ctx, booking.ID, workerID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) CheckBayAvailability(ctx context.Context, bayID int64, start, end time.Time) (bool, error) {
	query := `
		SELECT NOT EXISTS (
			SELECT 1 FROM bookings
			WHERE bay_id = $1
			AND status != 'canceled'
			AND start_time < $3
			AND end_time > $2
		)
	`
	var isFree bool
	err := r.db.QueryRowContext(ctx, query, bayID, start, end).Scan(&isFree)
	return isFree, err
}
