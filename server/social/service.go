package social

import (
	"context"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	AddFriend(ctx context.Context, dto AddFriendDTO) error
	RemoveFriend(ctx context.Context, dto RemoveFriendDTO) error
	AcceptFriendRequest(ctx context.Context, dto AcceptFriendRequestDTO) error
	RejectFriendRequest(ctx context.Context, dto RejectFriendRequestDTO) error
	GetAllIncomingFriendRequests(ctx context.Context, dto GetIncomingRequestsOfUserDTO) ([]RequestDTO, error)
	GetAllOutgoingFriendRequests(ctx context.Context, dto GetOutgoingRequestsOfUserDTO) ([]RequestDTO, error)
	BlockUser(ctx context.Context, dto BlockUserDTO) error
	UnblockUser(ctx context.Context, dto UnblockUserDTO) error
	GetAllFriends(ctx context.Context, dto GetAllFriendsDTO) ([]uuid.UUID, error)
	GetAllBlocked(ctx context.Context, dto GetAllBlockedDTO) ([]uuid.UUID, error)
	GetMutualFriends(ctx context.Context, dto GetMutualFriendsDTO) ([]uuid.UUID, error)
}

type LocalSocialService struct {
	repo   Repository
	config *config.Config
	logger *zap.Logger
}

func NewLocalSocialService(repo Repository, c *config.Config, l *zap.Logger) *LocalSocialService {
	return &LocalSocialService{
		repo:   repo,
		config: c,
		logger: l,
	}
}

func (s *LocalSocialService) AddFriend(ctx context.Context, dto AddFriendDTO) error {
	ctx, span := tracer.Start(ctx, "social.service.addfriend")
	defer span.End()

	return s.repo.AddFriend(ctx, dto)
}

func (s *LocalSocialService) RemoveFriend(ctx context.Context, dto RemoveFriendDTO) error {
	ctx, span := tracer.Start(ctx, "social.service.removefriend")
	defer span.End()

	return s.repo.RemoveFriend(ctx, dto)
}

func (s *LocalSocialService) AcceptFriendRequest(ctx context.Context, dto AcceptFriendRequestDTO) error {
	ctx, span := tracer.Start(ctx, "social.service.acceptfriendrequest")
	defer span.End()

	return s.repo.AcceptFriendRequest(ctx, dto)
}

func (s *LocalSocialService) RejectFriendRequest(ctx context.Context, dto RejectFriendRequestDTO) error {
	ctx, span := tracer.Start(ctx, "social.service.rejectfriendrequest")
	defer span.End()

	return s.repo.RejectFriendRequest(ctx, dto)
}

func (s *LocalSocialService) GetAllIncomingFriendRequests(ctx context.Context, dto GetIncomingRequestsOfUserDTO) ([]RequestDTO, error) {
	ctx, span := tracer.Start(ctx, "social.service.getallincomingfriendrequests")
	defer span.End()

	return s.repo.GetAllIncomingFriendRequests(ctx, dto)
}

func (s *LocalSocialService) GetAllOutgoingFriendRequests(ctx context.Context, dto GetOutgoingRequestsOfUserDTO) ([]RequestDTO, error) {
	ctx, span := tracer.Start(ctx, "social.service.getalloutgoingfriendrequests")
	defer span.End()

	return s.repo.GetAllOutgoingFriendRequests(ctx, dto)
}

func (s *LocalSocialService) BlockUser(ctx context.Context, dto BlockUserDTO) error {
	ctx, span := tracer.Start(ctx, "social.service.blockuser")
	defer span.End()

	return s.repo.BlockUser(ctx, dto)
}

func (s *LocalSocialService) UnblockUser(ctx context.Context, dto UnblockUserDTO) error {
	ctx, span := tracer.Start(ctx, "social.service.unblockuser")
	defer span.End()

	return s.repo.UnblockUser(ctx, dto)
}

func (s *LocalSocialService) GetAllFriends(ctx context.Context, dto GetAllFriendsDTO) ([]uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "social.service.getallfriends")
	defer span.End()

	return s.repo.GetAllFriends(ctx, dto)
}

func (s *LocalSocialService) GetAllBlocked(ctx context.Context, dto GetAllBlockedDTO) ([]uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "social.service.getallblocked")
	defer span.End()

	return s.repo.GetAllBlocked(ctx, dto)
}

func (s *LocalSocialService) GetMutualFriends(ctx context.Context, dto GetMutualFriendsDTO) ([]uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "social.service.getmutualfriends")
	defer span.End()

	return s.repo.GetMutualFriends(ctx, dto)
}
