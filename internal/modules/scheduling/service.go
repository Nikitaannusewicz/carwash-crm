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
