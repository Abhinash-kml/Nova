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
	AddFriend(ctx context.Context, req AddFriend_Request) error
	RemoveFriend(ctx context.Context, req RemoveFriend_Request) error
	AcceptFriendRequest(ctx context.Context, req AcceptFriendRequest_Request) error
	RejectFriendRequest(ctx context.Context, req RejectFriendRequest_Request) error
	GetAllIncomingFriendRequests(ctx context.Context, req GetIncomingRequestsOfUser_Request) ([]FriendRequest, error)
	GetAllOutgoingFriendRequests(ctx context.Context, req GetOutgoingRequestsOfUser_Request) ([]FriendRequest, error)
	BlockUser(ctx context.Context, req BlockUser_Request) error
	UnblockUser(ctx context.Context, req UnblockUser_Request) error
	GetAllFriends(ctx context.Context, req GetAllFriends_Request) ([]uuid.UUID, error)
	GetAllBlocked(ctx context.Context, req GetAllBlocked_Request) ([]uuid.UUID, error)
	GetMutualFriends(ctx context.Context, req GetMutualFriends_Request) ([]uuid.UUID, error)

	CheckIfConversationExists(ctx context.Context, conversationID uuid.UUID) (bool, error)
	CreateConversation(ctx context.Context, req CreateConversation_Request) (ConversationDetails_Response, error)
	GetConversationParticipants(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error)
	AddConversationParticipant(ctx context.Context, req AddConversationParticipant_Request) error
	RemoveConversationParticipant(ctx context.Context, req RemoveConversationParticipant_Request) error
	SendMessage(ctx context.Context, envelope realtime.Envelope) error
	CreateConversationSendMessage(ctx context.Context, message realtime.ChatMessage, senderID, targetID uuid.UUID) error
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

func (s *LocalSocialService) AddFriend(ctx context.Context, req AddFriend_Request) error {
	ctx, span := tracer.Start(ctx, "social.service.AddFriend")
	defer span.End()

	return common.TranslatePostgresError(
		s.repo.AddFriend(ctx, req),
		s.logger,
	).WithMessage("Failed to add Friend")
}

func (s *LocalSocialService) RemoveFriend(ctx context.Context, req RemoveFriend_Request) error {
	ctx, span := tracer.Start(ctx, "social.service.RemoveFriend")
	defer span.End()

	return common.TranslatePostgresError(
		s.repo.RemoveFriend(ctx, req),
		s.logger,
	).WithMessage("Failed to remove friend")
}

func (s *LocalSocialService) AcceptFriendRequest(ctx context.Context, req AcceptFriendRequest_Request) error {
	ctx, span := tracer.Start(ctx, "social.service.AcceptFriendRequest")
	defer span.End()

	return common.TranslatePostgresError(
		s.repo.AcceptFriendRequest(ctx, req),
		s.logger,
	).WithMessage("Failed to accept friend request")
}

func (s *LocalSocialService) RejectFriendRequest(ctx context.Context, req RejectFriendRequest_Request) error {
	ctx, span := tracer.Start(ctx, "social.service.RejectFriendRequest")
	defer span.End()

	return common.TranslatePostgresError(
		s.repo.RejectFriendRequest(ctx, req),
		s.logger,
	).WithMessage("Failed to reject friend request")
}

func (s *LocalSocialService) GetAllIncomingFriendRequests(ctx context.Context, req GetIncomingRequestsOfUser_Request) ([]FriendRequest, error) {
	ctx, span := tracer.Start(ctx, "social.service.GetAllIncomingFriendRequests")
	defer span.End()

	requests, err := s.repo.GetAllIncomingFriendRequests(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all incoming friend requests")
	}

	return requests, nil
}

func (s *LocalSocialService) GetAllOutgoingFriendRequests(ctx context.Context, req GetOutgoingRequestsOfUser_Request) ([]FriendRequest, error) {
	ctx, span := tracer.Start(ctx, "social.service.GetAllOutgoingFriendRequests")
	defer span.End()

	requests, err := s.repo.GetAllOutgoingFriendRequests(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all outgoing friend requests")
	}

	return requests, nil
}

func (s *LocalSocialService) BlockUser(ctx context.Context, req BlockUser_Request) error {
	ctx, span := tracer.Start(ctx, "social.service.BlockUser")
	defer span.End()

	return common.TranslatePostgresError(
		s.repo.BlockUser(ctx, req),
		s.logger,
	).WithMessage("Failed to block user")
}

func (s *LocalSocialService) UnblockUser(ctx context.Context, req UnblockUser_Request) error {
	ctx, span := tracer.Start(ctx, "social.service.UnblockUser")
	defer span.End()

	return common.TranslatePostgresError(
		s.repo.UnblockUser(ctx, req),
		s.logger,
	).WithMessage("Failed to unblock user")
}

func (s *LocalSocialService) GetAllFriends(ctx context.Context, req GetAllFriends_Request) ([]uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "social.service.GetAllFriends")
	defer span.End()

	friends, err := s.repo.GetAllFriends(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all friends")
	}

	return friends, nil
}

func (s *LocalSocialService) GetAllBlocked(ctx context.Context, req GetAllBlocked_Request) ([]uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "social.service.GetAllBlocked")
	defer span.End()

	blockedUsers, err := s.repo.GetAllBlocked(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all blocked users")
	}

	return blockedUsers, nil
}

func (s *LocalSocialService) GetMutualFriends(ctx context.Context, req GetMutualFriends_Request) ([]uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "social.service.GetMutualFriends")
	defer span.End()

	mutualFriends, err := s.repo.GetMutualFriends(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get mutual friends")
	}

	return mutualFriends, nil
}

func (s *LocalSocialService) CheckIfConversationExists(ctx context.Context, conversationID uuid.UUID) (bool, error) {
	ctx, span := tracer.Start(ctx, "social.service.CheckIfConversationExists")
	defer span.End()

	return s.repo.CheckIfConversationExists(ctx, conversationID)
}

func (s *LocalSocialService) CreateConversation(ctx context.Context, req CreateConversation_Request) (ConversationDetails_Response, error) {
	ctx, span := tracer.Start(ctx, "social.service.CreateConversation")
	defer span.End()

	return s.repo.CreateConversation(ctx, req)
}

func (s LocalSocialService) GetConversationParticipants(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "social.service.GetConversationParticipants")
	defer span.End()

	return s.repo.GetConversationParticipants(ctx, conversationID)
}

func (s *LocalSocialService) AddConversationParticipant(ctx context.Context, req AddConversationParticipant_Request) error {
	ctx, span := tracer.Start(ctx, "social.service.AddConversationParticipant")
	defer span.End()

	return s.repo.AddConversationParticipant(ctx, req)
}

func (s *LocalSocialService) RemoveConversationParticipant(ctx context.Context, req RemoveConversationParticipant_Request) error {
	ctx, span := tracer.Start(ctx, "social.service.RemoveConversationParticipant")
	defer span.End()

	return s.repo.RemoveConversationParticipant(ctx, req)
}

func (s *LocalSocialService) SendMessage(ctx context.Context, envelope realtime.Envelope) error {
	ctx, span := tracer.Start(ctx, "social.service.SendMessage")
	defer span.End()

	return s.repo.SendMessage(ctx, envelope)
}

func (s *LocalSocialService) CreateConversationSendMessage(ctx context.Context, message realtime.ChatMessage, senderID, targetID uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "social.service.CreateConversationSendMessage")
	defer span.End()

	return s.repo.CreateConversationSendMessage(ctx, message, senderID, targetID)
}
