package operations

import "time"

type Location struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

type Bay struct {
	ID         int64  `json:"id"`
	LocationID string `json:"location_id"`
	Name       string `json:"name"`
	IsActive   bool   `json:"is_active"`
}

type Service struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	DurationMinutes int    `json:"duration_minutes"`

	// AI Note: Price in cents to avoid floating point errors with currency
	PriceCents int64 `json:"price_cents"`
}
