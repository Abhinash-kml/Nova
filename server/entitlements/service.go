package entitlements

import (
	"context"

	"github.com/abhinash-kml/nova/server/apperr"
	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

//go:generate mockery
type Service interface {
	Create(ctx context.Context, req CreateRequest) (Entitlement, error)
	Delete(ctx context.Context, req DeleteRequest) (DeleteResponse, error)
	CheckIfActive(ctx context.Context, req CheckIfActiveRequest) (bool, error)

	GetAllOfUser(ctx context.Context, req GetAllOfUserRequest) ([]Entitlement, error)
}

type LocalEntitlementService struct {
	repo   Repository
	logger *zap.Logger
}

func NewLocalEntitlementService(r Repository, l *zap.Logger) *LocalEntitlementService {
	return &LocalEntitlementService{
		repo:   r,
		logger: l,
	}
}

func (s *LocalEntitlementService) Create(ctx context.Context, req CreateRequest) (Entitlement, error) {
	ctx, span := tracer.Start(ctx, "entitlements.service.Create")
	defer span.End()

	if req.Source == "" {
		return Entitlement{}, apperr.InvalidInput(nil, "cannot create entitlement with no source")
	}
	if !req.IsActive {
		return Entitlement{}, apperr.InvalidInput(nil, "cannot create disabled entitlement")
	}

	ent, err := s.repo.Create(ctx, req)
	if err != nil {
		return Entitlement{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to create entitlement")
	}

	return ent, nil
}

func (s *LocalEntitlementService) Delete(ctx context.Context, req DeleteRequest) (DeleteResponse, error) {
	ctx, span := tracer.Start(ctx, "entitlements.service.Delete")
	defer span.End()

	if uuid.Validate(req.Id) != nil {
		return DeleteResponse{}, apperr.InvalidInput(nil, "cannot delete entitlement with empty id")
	}

	id, err := s.repo.Delete(ctx, req)
	if err != nil {
		return id, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to delete entitlement")
	}

	return id, nil
}

func (s *LocalEntitlementService) CheckIfActive(ctx context.Context, req CheckIfActiveRequest) (bool, error) {
	ctx, span := tracer.Start(ctx, "entitlements.service.CheckIfActive")
	defer span.End()

	if uuid.Validate(req.Id) != nil {
		return false, apperr.InvalidInput(nil, "cannot check with invalid entitlement id")
	}

	isActive, err := s.repo.CheckIfActive(ctx, req)
	if err != nil {
		return false, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to check if entitlement was active")
	}

	return isActive, nil
}

func (s *LocalEntitlementService) GetAllOfUser(ctx context.Context, req GetAllOfUserRequest) ([]Entitlement, error) {
	ctx, span := tracer.Start(ctx, "entitlements.service.GetAllOfuser")
	defer span.End()

	if uuid.Validate(req.UserID) != nil {
		return nil, apperr.InvalidInput(nil, "cannot get all entitlements of user with invalid uuid")
	}

	entitlements, err := s.repo.GetAllOfUser(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all entitlements of user")
	}

	return entitlements, nil
}
