package scheduling

import (
	"context"
	"errors"
	"time"

	"github.com/Nikitaannusewicz/carwash-crm/internal/modules/identity"
)

type Service struct {
	repo Repository
}

type Repository interface {
	CreateShift(ctx context.Context, shift *Shift) error
	CheckOverLap(ctx context.Context, userID int64, start, end time.Time) (bool, error)
	CreateBooking(ctx context.Context, booking *Booking) error
	CheckBayAvailability(ctx context.Context, bayID int64, start, end time.Time) (bool, error)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

type CreateShiftRequest struct {
	UserID     int64     `json:"user_id"`
	LocationID int64     `json:"location_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
}

type CreateBookingRequest struct {
	LocationID int64     `json:"location_id"`
	BayID      int64     `json:"bay_id"`
	ServiceID  int64     `json:"service_id"`
	CustomerID int64     `json:"customer_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	WorkerIDs  []int64   `json:"worker_ids"`
}

func (s *Service) CreateShift(ctx context.Context, req CreateShiftRequest, requesterID int64, requesterRole identity.Role) (*Shift, error) {
	// Authorization: Admins and managers can create shifts for anyone, workers can only create their own
	if requesterRole == identity.RoleWorker && req.UserID != requesterID {
		return nil, errors.New("unauthorized: workers can create only their own shifts")
	}

	if req.StartTime.After(req.EndTime) {
		return nil, errors.New("start time must be before end time")
	}

	if req.StartTime.Before(time.Now()) {
		return nil, errors.New("shift cannot be in the past")
	}

	exists, err := s.repo.CheckOverLap(ctx, req.UserID, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("worker already has a shift overlapping this time")
	}

	shift := &Shift{
		UserID:     req.UserID,
		LocationID: req.LocationID,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
	}

	if err := s.repo.CreateShift(ctx, shift); err != nil {
		return nil, err
	}

	return shift, nil
}

func (s *Service) CreateBooking(ctx context.Context, req CreateBookingRequest, role identity.Role) (*Booking, error) {
	durationMinutes := 60 // Temporary implementation

	if role != identity.RoleWorker {
		return nil, errors.New("Unauthorized")
	}

	if durationMinutes <= 0 {
		return nil, errors.New("duration must be positive")
	}

	if len(req.WorkerIDs) <= 0 {
		return nil, errors.New("at least one worker must be assigned to a booking")
	}

	endTime := req.StartTime.Add(time.Duration(durationMinutes) * time.Minute)

	isFree, err := s.repo.CheckBayAvailability(ctx, req.BayID, req.StartTime, endTime)
	if err != nil {
		return nil, err
	}
	if !isFree {
		return nil, errors.New("bay is occupied during this time")
	}

	booking := &Booking{
		LocationID:        req.LocationID,
		BayID:             req.BayID,
		ServiceID:         req.ServiceID,
		CustomerID:        req.CustomerID,
		StartTime:         req.StartTime,
		EndTime:           endTime,
		Status:            StatusPending,
		AssignedWorkerIDs: req.WorkerIDs,
	}

	if err := s.repo.CreateBooking(ctx, booking); err != nil {
		return nil, err
	}

	return booking, nil
}
