package scheduling

import "time"

type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"   // Created, not yet staffed
	StatusConfirmed BookingStatus = "confirmed" // Staffed and ready
	StatusCompleted BookingStatus = "completed" // Wash finished
	StatusCancelled BookingStatus = "cancelled" // User or Admin canceled
)

type Booking struct {
	ID         int64 `json:"id"`
	LocationID int64 `json:"location_id"`
	BayID      int64 `json:"bay_id"`
	CustomerID int64 `json:"customer_id"`
	ServiceID  int64 `json:"service_id"`

	StartTime time.Time `json:"start_time"`
	// EndTime = StartTime + Service.Duration
	EndTime time.Time `json:"end_time"`

	Status BookingStatus `json:"status"`

	// AssignedWorkerIDs tracks which employees are working on this car
	// to manage employee availability
	AssignedWorkerIDs []int64 `json:"assigned_workers_ids"`
}

type Shift struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	LocationID int64     `json:"location_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
}
