package operations

import (
	"context"
	"database/sql"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateLocation(ctx context.Context, loc *Location) error {
	query := `
		INSERT INTO locations (name, address, created_at)
		VALUES($1, $2, $3)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query, loc.Name, loc.Address, loc.CreatedAt).Scan(&loc.ID)
}

func (r *PostgresRepository) CreateBay(ctx context.Context, bay *Bay) error {
	query := `
		INSERT INTO bay (location_id, name, is_active)
		VALUES($1, $2, $3)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query, bay.LocationID, bay.Name, bay.IsActive).Scan(&bay.ID)
}

func (r *PostgresRepository) CreateService(ctx context.Context, ser *Service) error {
	query := `
		INSERT INTO services (name, duration_minutes, price_cents)
		VALUES ($1,$2,$3)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query, ser.Name, ser.DurationMinutes, ser.PriceCents).Scan(&ser.ID)
}
