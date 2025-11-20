package operations

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Nikitaannusewicz/carwash-crm/internal/modules/identity"
)

type Repository interface {
	CreateLocation(ctx context.Context, loc *Location) error
	CreateBay(ctx context.Context, bay *Bay) error
}

type OperationsService struct {
	repo Repository
}

func NewService(repo Repository) *OperationsService {
	return &OperationsService{repo: repo}
}

type CreateBayRequest struct {
	LocationID string `json:"location_id"`
	Name       string `json:"name"`
	IsActive   bool   `json:"is_active"`
}

type CreateLocationRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (s *OperationsService) CreateBay(ctx context.Context, req CreateBayRequest, requesterRole identity.Role) (*Bay, error) {
	if requesterRole != identity.RoleAdmin {
		return nil, errors.New("unauthorized: only admins can create bays")
	}

	if req.LocationID == "" || req.Name == "" {
		return nil, errors.New("LocationID and Name are required")
	}

	b := &Bay{
		LocationID: req.LocationID,
		Name:       req.Name,
		IsActive:   req.IsActive,
	}

	if err := s.repo.CreateBay(ctx, b); err != nil {
		return nil, err
	}

	return b, nil
}

func (s *OperationsService) CreateLocation(ctx context.Context, req CreateLocationRequest, requesterRole identity.Role) (*Location, error) {
	if requesterRole != identity.RoleAdmin {
		return nil, errors.New("unauthorized: only admins can create locations")
	}

	if req.Name == "" || req.Address == "" {
		fmt.Printf("name: %v, address: %v:", req.Name, req.Address)
		return nil, errors.New("name and adress are required")
	}

	loc := &Location{
		Name:      req.Name,
		Address:   req.Address,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateLocation(ctx, loc); err != nil {
		return nil, err
	}

	return loc, nil
}
