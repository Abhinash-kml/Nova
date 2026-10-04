package social

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/abhinash-kml/nova/server/realtime"
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

	CheckIfConversationExists(ctx context.Context) (bool, error)
	CreateConversation(ctx context.Context, dto CreateConversationDTO) (ConversationDetailsDTO, error)
	AddConversationParticipant(ctx context.Context, dto AddConversationParticipantDTO) error
	RemoveConversationParticipant(ctx context.Context, dto RemoveConversationParticipantDTO) error
	SendMessage(ctx context.Context, envelope realtime.Envelope) error
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

	return common.TranslatePostgresError(
		s.repo.AddFriend(ctx, dto),
		s.logger,
	).WithMessage("Failed to add Friend")
}

func (s *LocalSocialService) RemoveFriend(ctx context.Context, dto RemoveFriendDTO) error {
	ctx, span := tracer.Start(ctx, "social.service.removefriend")
	defer span.End()

	return common.TranslatePostgresError(
		s.repo.RemoveFriend(ctx, dto),
		s.logger,
	).WithMessage("Failed to remove friend")
}

func (s *LocalSocialService) AcceptFriendRequest(ctx context.Context, dto AcceptFriendRequestDTO) error {
	ctx, span := tracer.Start(ctx, "social.service.acceptfriendrequest")
	defer span.End()

	return common.TranslatePostgresError(
		s.repo.AcceptFriendRequest(ctx, dto),
		s.logger,
	).WithMessage("Failed to accept friend request")
}

func (s *LocalSocialService) RejectFriendRequest(ctx context.Context, dto RejectFriendRequestDTO) error {
	ctx, span := tracer.Start(ctx, "social.service.rejectfriendrequest")
	defer span.End()

	return common.TranslatePostgresError(
		s.repo.RejectFriendRequest(ctx, dto),
		s.logger,
	).WithMessage("Failed to reject friend request")
}

func (s *LocalSocialService) GetAllIncomingFriendRequests(ctx context.Context, dto GetIncomingRequestsOfUserDTO) ([]RequestDTO, error) {
	ctx, span := tracer.Start(ctx, "social.service.getallincomingfriendrequests")
	defer span.End()

	requests, err := s.repo.GetAllIncomingFriendRequests(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all incoming friend requests")
	}

	return requests, nil
}

func (s *LocalSocialService) GetAllOutgoingFriendRequests(ctx context.Context, dto GetOutgoingRequestsOfUserDTO) ([]RequestDTO, error) {
	ctx, span := tracer.Start(ctx, "social.service.getalloutgoingfriendrequests")
	defer span.End()

	requests, err := s.repo.GetAllOutgoingFriendRequests(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all outgoing friend requests")
	}

	return requests, nil
}

func (s *LocalSocialService) BlockUser(ctx context.Context, dto BlockUserDTO) error {
	ctx, span := tracer.Start(ctx, "social.service.blockuser")
	defer span.End()

	return common.TranslatePostgresError(
		s.repo.BlockUser(ctx, dto),
		s.logger,
	).WithMessage("Failed to block user")
}

func (s *LocalSocialService) UnblockUser(ctx context.Context, dto UnblockUserDTO) error {
	ctx, span := tracer.Start(ctx, "social.service.unblockuser")
	defer span.End()

	return common.TranslatePostgresError(
		s.repo.UnblockUser(ctx, dto),
		s.logger,
	).WithMessage("Failed to unblock user")
}

func (s *LocalSocialService) GetAllFriends(ctx context.Context, dto GetAllFriendsDTO) ([]uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "social.service.getallfriends")
	defer span.End()

	friends, err := s.repo.GetAllFriends(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all friends")
	}

	return friends, nil
}

func (s *LocalSocialService) GetAllBlocked(ctx context.Context, dto GetAllBlockedDTO) ([]uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "social.service.getallblocked")
	defer span.End()

	blockedUsers, err := s.repo.GetAllBlocked(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all blocked users")
	}

	return blockedUsers, nil
}

func (s *LocalSocialService) GetMutualFriends(ctx context.Context, dto GetMutualFriendsDTO) ([]uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "social.service.getmutualfriends")
	defer span.End()

	mutualFriends, err := s.repo.GetMutualFriends(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get mutual friends")
	}

	return mutualFriends, nil
}

func (s *LocalSocialService) CheckIfConversationExists(ctx context.Context) (bool, error) {
	return s.repo.CheckIfConversationExists(ctx)
}

func (s *LocalSocialService) CreateConversation(ctx context.Context, dto CreateConversationDTO) (ConversationDetailsDTO, error) {
	return s.repo.CreateConversation(ctx, dto)
}

func (s *LocalSocialService) AddConversationParticipant(ctx context.Context, dto AddConversationParticipantDTO) error {
	return s.repo.AddConversationParticipant(ctx, dto)
}

func (s *LocalSocialService) RemoveConversationParticipant(ctx context.Context, dto RemoveConversationParticipantDTO) error {
	return s.repo.RemoveConversationParticipant(ctx, dto)
}

func (s *LocalSocialService) SendMessage(ctx context.Context, envelope realtime.Envelope) error {
	return s.repo.SendMessage(ctx, envelope)
}
