package entitlements

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"go.uber.org/zap"
)

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

	ent, err := s.repo.Create(ctx, req)
	if err != nil {
		return Entitlement{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to create entitlement")
	}

	return ent, nil
}

func (s *LocalEntitlementService) Delete(ctx context.Context, req DeleteRequest) (DeleteResponse, error) {
	ctx, span := tracer.Start(ctx, "entitlements.service.Delete")
	defer span.End()

	id, err := s.repo.Delete(ctx, req)
	if err != nil {
		return id, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to delete entitlement")
	}

	return id, nil
}

func (s *LocalEntitlementService) CheckIfActive(ctx context.Context, req CheckIfActiveRequest) (bool, error) {
	ctx, span := tracer.Start(ctx, "entitlements.service.CheckIfActive")
	defer span.End()

	isActive, err := s.repo.CheckIfActive(ctx, req)
	if err != nil {
		return false, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to check if entitlement was active")
	}

	return isActive, nil
}

func (s *LocalEntitlementService) GetAllOfUser(ctx context.Context, req GetAllOfUserRequest) ([]Entitlement, error) {
	ctx, span := tracer.Start(ctx, "entitlements.service.GetAllOfuser")
	defer span.End()

	entitlements, err := s.repo.GetAllOfUser(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all entitlements of user")
	}

	return entitlements, nil
}
