package identity

import (
	"context"
	"time"
)

type Role string

const (
	RoleAdmin   Role = "admin"   // Can manage all location
	RoleManager Role = "manager" // Can manage one location
	RoleWorker  Role = "worker"  // Assigned to wash card
)

type User struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Role         Role   `json:"role"`

	// AI Note: LocationID is nullable because Admins might not belong to a specific location
	// In Go, we often use a pointer to int for nullable integers.
	LocationID *int64 `json:"location_id,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
}
