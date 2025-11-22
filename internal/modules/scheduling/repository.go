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
