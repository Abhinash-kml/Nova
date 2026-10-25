package clans

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Service interface {
	// General operations
	Add(ctx context.Context, dto CreateRequest) (Clan, error)
	GetById(ctx context.Context, id uuid.UUID) (Clan, error)
	GetByName(ctx context.Context, name string) (Clan, error)
	GetAll(ctx context.Context, cursor uuid.UUID, limit int) ([]Clan, error)
	Update(ctx context.Context, req UpdateRequest) (Clan, error)
	Delete(ctx context.Context, req DeleteRequest) (DeleteResponse, error)

	AddMember(ctx context.Context, req ActionRequest) (ActionResponse, error)
	RemoveMember(ctx context.Context, req ActionRequest) (ActionResponse, error)
	ListAllOutogingInvites(ctx context.Context, req ActionRequest) ([]OutgoingInvite, error)
	InviteMember(ctx context.Context, req ActionRequest) (ActionResponse, error)
	AcceptInvite(ctx context.Context, req ActionRequest) (ActionResponse, error)
	RejectInvite(ctx context.Context, req ActionRequest) (ActionResponse, error)
	ListAllJoinRequests(ctx context.Context, req ActionRequest) ([]uuid.UUID, error)
	SendJoinRequest(ctx context.Context, req ActionRequest) (ActionResponse, error)
	AcceptJoinRequest(ctx context.Context, req ActionRequest) (ActionResponse, error)
	RejectJoinRequest(ctx context.Context, req ActionRequest) (ActionResponse, error)
}

type LocalClansService struct {
	repo   ClansRepository
	logger *zap.Logger
	config *config.Config
	cache  *redis.Client
}

func NewLocalClansService(repo ClansRepository, c *config.Config, r *redis.Client, l *zap.Logger) *LocalClansService {
	return &LocalClansService{
		repo:   repo,
		config: c,
		cache:  r,
		logger: l,
	}
}

func (s *LocalClansService) GetById(ctx context.Context, id uuid.UUID) (Clan, error) {
	ctx, span := tracer.Start(ctx, "clans.service.getbyid")
	defer span.End()

	key := ClanPrefix + id.String()

	// 1. Try cache
	var clan Clan
	err := s.cache.Get(ctx, key).Scan(&clan)
	if err == nil && len(clan.Name) != 0 {
		return clan, nil
	}

	// If Redis failed for infra reason, log but continue
	if err != nil && err != redis.Nil {
		s.logger.Warn("cache error", zap.Error(err))
	}

	// 2. Fallback to repo
	clan, err = s.repo.GetById(ctx, id)
	if err != nil {
		return Clan{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get clan")
	}

	// 3. Populate cache asynchronously (safe version)
	go func(c Clan, key string) {
		bgCtx := context.WithoutCancel(ctx)

		_, err := s.cache.Set(bgCtx, key, &c, 0).Result()

		if err != nil {
			s.logger.Error("failed to populate cache", zap.Error(err))
		}
	}(clan, key)

	return clan, nil
}

func (s *LocalClansService) GetByName(ctx context.Context, name string) (Clan, error) {
	ctx, span := tracer.Start(ctx, "clans.service.getbyname")
	defer span.End()

	// Caching logic

	clan, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return Clan{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get clan by name")
	}

	return clan, nil
}

func (s *LocalClansService) GetAll(ctx context.Context, cursor uuid.UUID, limit int) ([]Clan, error) {
	ctx, span := tracer.Start(ctx, "clans.service.getall")
	defer span.End()

	clans, err := s.repo.GetAll(ctx, cursor, limit)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all clans")
	}

	return clans, nil
}

func (s *LocalClansService) Add(ctx context.Context, dto CreateRequest) (Clan, error) {
	ctx, span := tracer.Start(ctx, "clans.service.add")
	defer span.End()

	clan, err := s.repo.Add(ctx, dto)
	if err != nil {
		return Clan{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to create clan")
	}

	return clan, nil
}

func (s *LocalClansService) Delete(ctx context.Context, dto DeleteRequest) (DeleteResponse, error) {
	ctx, span := tracer.Start(ctx, "clans.service.delete")
	defer span.End()

	// Delete in repo
	response, err := s.repo.Delete(ctx, dto)
	if err != nil {
		return DeleteResponse{DeletedClanID: uuid.Nil}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to delete clan")
	}

	// Delete from cache
	go func() {
		bgCtx := context.WithoutCancel(ctx)
		key := ClanPrefix + dto.Id
		err := s.cache.Del(bgCtx, key).Err()
		if err != nil {
			s.logger.Error("Failed to delete clan from cache", zap.Error(err))
		}
	}()

	return response, nil
}

func (s *LocalClansService) Update(ctx context.Context, dto UpdateRequest) (Clan, error) {
	ctx, span := tracer.Start(ctx, "clans.service.update")
	defer span.End()

	// Update repository first
	updatedClan, err := s.repo.Update(ctx, dto)
	if err != nil {
		return Clan{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to update clan")
	}

	// Invalidate old record from cache, next get call with repopulate it
	go func() {
		bgCtx := context.WithoutCancel(ctx)
		key := ClanPrefix + dto.Id
		err := s.cache.Del(bgCtx, key).Err()
		if err != nil {
			s.logger.Error("Failed to delete clan from cache", zap.Error(err))
		}
	}()

	return updatedClan, nil
}

func (s *LocalClansService) AddMember(ctx context.Context, req ActionRequest) (ActionResponse, error) {
	ctx, span := tracer.Start(ctx, "clans.service.addmember")
	defer span.End()

	res, err := s.repo.AddMember(ctx, req)
	if err != nil {
		return ActionResponse{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to add member")
	}

	return res, nil
}

func (s *LocalClansService) RemoveMember(ctx context.Context, req ActionRequest) (ActionResponse, error) {
	ctx, span := tracer.Start(ctx, "clans.service.removemember")
	defer span.End()

	res, err := s.repo.RemoveMember(ctx, req)
	if err != nil {
		return ActionResponse{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to remove member")
	}

	return res, nil
}

func (s *LocalClansService) ListAllOutogingInvites(ctx context.Context, req ActionRequest) ([]OutgoingInvite, error) {
	ctx, span := tracer.Start(ctx, "clans.service.listalloutgoinginvites")
	defer span.End()

	invites, err := s.repo.ListAllOutogingInvites(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to list all outgoing invites")
	}

	return invites, nil
}

func (s *LocalClansService) InviteMember(ctx context.Context, req ActionRequest) (ActionResponse, error) {
	ctx, span := tracer.Start(ctx, "clans.service.invitemember")
	defer span.End()

	res, err := s.repo.InviteMember(ctx, req)
	if err != nil {
		return ActionResponse{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to invite member")
	}

	return res, nil
}

func (s *LocalClansService) AcceptInvite(ctx context.Context, req ActionRequest) (ActionResponse, error) {
	ctx, span := tracer.Start(ctx, "clans.service.acceptinvite")
	defer span.End()

	res, err := s.repo.AcceptInvite(ctx, req)
	if err != nil {
		return ActionResponse{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to accept invite")
	}

	return res, nil
}

func (s *LocalClansService) RejectInvite(ctx context.Context, req ActionRequest) (ActionResponse, error) {
	ctx, span := tracer.Start(ctx, "clans.service.rejectinvite")
	defer span.End()

	res, err := s.repo.RejectInvite(ctx, req)
	if err != nil {
		return ActionResponse{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to reject invite")
	}

	return res, nil
}

func (s *LocalClansService) ListAllJoinRequests(ctx context.Context, req ActionRequest) ([]uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "clans.service.listalljoinrequests")
	defer span.End()

	response, err := s.repo.ListAllJoinRequests(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to list all join requests")
	}

	return response, nil
}

func (s *LocalClansService) SendJoinRequest(ctx context.Context, req ActionRequest) (ActionResponse, error) {
	ctx, span := tracer.Start(ctx, "clans.service.sendjoinrequest")
	defer span.End()

	res, err := s.repo.SendJoinRequest(ctx, req)
	if err != nil {
		return ActionResponse{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to send join request")
	}

	return res, nil
}

func (s *LocalClansService) AcceptJoinRequest(ctx context.Context, req ActionRequest) (ActionResponse, error) {
	ctx, span := tracer.Start(ctx, "clans.service.acceptjoinrequest")
	defer span.End()

	res, err := s.repo.AcceptJoinRequest(ctx, req)
	if err != nil {
		return ActionResponse{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to accept join request")
	}

	return res, nil
}

func (s *LocalClansService) RejectJoinRequest(ctx context.Context, req ActionRequest) (ActionResponse, error) {
	ctx, span := tracer.Start(ctx, "clans.service.rejectjoinrequest")
	defer span.End()

	res, err := s.repo.RejectJoinRequest(ctx, req)
	if err != nil {
		return ActionResponse{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to reject join request")
	}

	return res, nil
}
