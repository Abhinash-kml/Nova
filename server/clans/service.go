package clans

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/abhinash-kml/nova/server/cache"
	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
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

	GetAllMembers(ctx context.Context, req GetRequest) ([]ClanMember, error)

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
	config *config.ClansConfig
	cache  *cache.TieredCache
}

func NewLocalClansService(repo ClansRepository, c *config.ClansConfig, tc *cache.TieredCache, l *zap.Logger) *LocalClansService {
	return &LocalClansService{
		repo:   repo,
		config: c,
		cache:  tc,
		logger: l,
	}
}

func (s *LocalClansService) GetById(ctx context.Context, id uuid.UUID) (Clan, error) {
	ctx, span := tracer.Start(ctx, "clans.service.getbyid")
	defer span.End()

	key := s.config.KeyPrefix + id.String()

	// Get from cache or load into it
	bytes, err := s.cache.GetOrLoad(ctx, key, func(ctx context.Context) ([]byte, error) {
		clan, err := s.repo.GetById(ctx, id)
		if err != nil {
			return nil, err
		}

		return json.Marshal(clan)
	})
	if err != nil {
		return Clan{}, err
	}

	var clan Clan
	err = json.Unmarshal(bytes, &clan)
	if err != nil {
		return Clan{}, fmt.Errorf("unmarshalling clan: %w", err)
	}

	return clan, nil
}

func (s *LocalClansService) GetByName(ctx context.Context, name string) (Clan, error) {
	ctx, span := tracer.Start(ctx, "clans.service.getbyname")
	defer span.End()

	key := s.config.KeyPrefix + name

	// Get from cache or load into it
	bytes, err := s.cache.GetOrLoad(ctx, key, func(ctx context.Context) ([]byte, error) {
		clan, err := s.repo.GetByName(ctx, name)
		if err != nil {
			return nil, err
		}

		return json.Marshal(clan)
	})
	if err != nil {
		return Clan{}, err
	}

	var clan Clan
	err = json.Unmarshal(bytes, &clan)
	if err != nil {
		return Clan{}, fmt.Errorf("unmarshalling clan: %w", err)
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

	// Invalidate old record from cache, next get call with repopulate it
	key := s.config.KeyPrefix + dto.Id
	go s.cache.Delete(ctx, key)

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
	key := s.config.KeyPrefix + dto.Id
	go s.cache.Delete(ctx, key)

	return updatedClan, nil
}

func (s *LocalClansService) GetAllMembers(ctx context.Context, req GetRequest) ([]ClanMember, error) {
	ctx, span := tracer.Start(ctx, "clans.service.getallmembers")
	defer span.End()

	members, err := s.repo.GetAllMembers(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all clan members")
	}

	return members, nil
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
