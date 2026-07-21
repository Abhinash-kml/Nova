package users

import (
	"context"
	"errors"
	"time"

	"github.com/abhinash-kml/nova/server/auth"
	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type Service interface {
	// General operations
	Add(ctx context.Context, dto CreateDTO) (User, error)
	GetAll(ctx context.Context, cursor int, limit int) ([]User, error)
	GetAllByAttribute(ctx context.Context, attribute string) ([]User, error)
	GetById(ctx context.Context, id uuid.UUID) (User, error)
	GetByName(ctx context.Context, name string) (User, error)
	GetBySocialProfileId(ctx context.Context, provider, socialUserId string) (User, error)
	Update(ctx context.Context, dto UpdateDTO) (User, error)
	Replace(ctx context.Context, dto ReplaceDTO) (User, error)
	Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error)

	// Auth operations
	Login(ctx context.Context, loginRequest auth.LoginRequest) (auth.LoginResponse, error)
	Refresh(ctx context.Context, dto auth.TokenRefreshRequest) (auth.TokenRefreshResponse, error)

	// Operations
	FindOrCreateUserFromSocialProfile(ctx context.Context, profile auth.UnifiedUserProfile) (uuid.UUID, bool, error)
	CheckIfUserIsBanned(ctx context.Context, userId uuid.UUID) bool

	// Bulk operations
	BulkAdd(ctx context.Context, dto BulkCreateDTO) ([]common.BulkOpResult, error)
	BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error)
	BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error)
}

type LocalUsersService struct {
	repo   UsersRepository
	logger *zap.Logger
	cache  *redis.Client
	config *config.Config
}

func NewLocalUsersService(repository UsersRepository, r *redis.Client, l *zap.Logger, c *config.Config) *LocalUsersService {
	return &LocalUsersService{
		repo:   repository,
		cache:  r,
		logger: l,
		config: c,
	}
}

func (s *LocalUsersService) Add(ctx context.Context, user CreateDTO) (User, error) {
	ctx, span := tracer.Start(ctx, "users.service.add")
	defer span.End()

	return s.repo.Add(ctx, user)
}

func (s *LocalUsersService) GetAll(ctx context.Context, cursor, count int) ([]User, error) {
	ctx, span := tracer.Start(ctx, "users.service.getall")
	defer span.End()

	return s.repo.GetAll(ctx, cursor, count)
}

func (s *LocalUsersService) GetAllByAttribute(ctx context.Context, attribute string) ([]User, error) {
	ctx, span := tracer.Start(ctx, "users.service.getallbyattribute")
	defer span.End()

	return s.repo.GetAllByAttribute(ctx, attribute)
}

func (s *LocalUsersService) GetById(ctx context.Context, id uuid.UUID) (User, error) {
	ctx, span := tracer.Start(ctx, "users.service.getbyid")
	defer span.End()

	key := UserPrefix + id.String()

	// 1. Try cache
	var user User
	err := s.cache.Get(ctx, key).Scan(&user)
	if err == nil && len(user.Username) != 0 {
		return user, nil
	}

	// If Redis failed for infra reason, log but continue
	if err != nil && err != redis.Nil {
		s.logger.Warn("cache error", zap.Error(err))
	}

	// 2. Fallback to repo
	user, err = s.repo.GetById(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return User{}, common.ErrResourceNotFound
	}

	// 3. Populate cache asynchronously
	go func(u User, key string) {
		bgCtx := context.WithoutCancel(ctx)
		_, err := s.cache.Set(bgCtx, key, &u, 0).Result()
		if err != nil {
			s.logger.Error("failed to populate cache", zap.Error(err))
		}
	}(user, key)

	return user, nil
}

func (s *LocalUsersService) GetByName(ctx context.Context, name string) (User, error) {
	ctx, span := tracer.Start(ctx, "users.service.getbyname")
	defer span.End()

	// Get from cache

	// Get from repository
	return s.repo.GetByName(ctx, name)
}

func (s *LocalUsersService) Update(ctx context.Context, dto UpdateDTO) (User, error) {
	ctx, span := tracer.Start(ctx, "users.service.update")
	defer span.End()

	// Update repository first
	user, err := s.repo.Update(ctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return User{}, err
	}

	// Invalidate old record from cache, next get call with repopulate it
	go func() {
		bgCtx := context.WithoutCancel(ctx)
		key := UserPrefix + dto.Id
		err := s.cache.Del(bgCtx, key).Err()
		if err != nil {
			s.logger.Error("Failed to delete user from cache", zap.Error(err))
		}
	}()

	return user, nil
}

func (s *LocalUsersService) Replace(ctx context.Context, dto ReplaceDTO) (User, error) {
	ctx, span := tracer.Start(ctx, "users.service.replace")
	defer span.End()

	// Update repository first
	user, err := s.repo.Replace(ctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return User{}, err
	}

	// Invalidate old record from cache, next get call with repopulate it
	go func() {
		bgCtx := context.WithoutCancel(ctx)
		key := UserPrefix + dto.Id
		err := s.cache.Del(bgCtx, key).Err()
		if err != nil {
			s.logger.Error("Failed to delete user from cache", zap.Error(err))
		}
	}()

	return user, nil
}

func (s *LocalUsersService) Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "users.service.delete")
	defer span.End()

	// Delete in repo
	deletedId, err := s.repo.Delete(ctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return uuid.Nil, err
	}

	// Delete from cache
	go func() {
		bgCtx := context.WithoutCancel(ctx)
		key := UserPrefix + dto.Id
		err := s.cache.Del(bgCtx, key).Err()
		if err != nil {
			s.logger.Error("Failed to delete user from cache", zap.Error(err))
		}
	}()

	return deletedId, nil
}

// Auth operations
func (s *LocalUsersService) Login(ctx context.Context, loginRequest auth.LoginRequest) (auth.LoginResponse, error) {
	ctx, span := tracer.Start(ctx, "users.service.login")
	defer span.End()

	// 1. Exchange provider token for social profile validation
	profile, err := auth.GetSocialAuthEngine().CompleteAuthentication(ctx, loginRequest.Provider, loginRequest.Code)
	if err != nil {
		return auth.LoginResponse{
			Status: "failed",
			Reason: err.Error(),
		}, errors.Join(auth.ErrLoginFailed, err)
	}

	s.logger.Debug("Profile", zap.String("provider", profile.Provider),
		zap.String("display_name", profile.DisplayName),
		zap.String("userid", profile.UserId),
		zap.String("email", profile.Email),
		zap.String("avatar_url", profile.AvatarUrl))

	// 2. Atomic Find user if exists else create a new one on spot
	userId, createdRightNow, err := s.FindOrCreateUserFromSocialProfile(ctx, profile)
	if err != nil {
		return auth.LoginResponse{
			Status: "failed",
			Reason: err.Error(),
		}, errors.Join(auth.ErrLoginFailed, err)
	}

	// 3. User exists -> Check if account is banned
	if !createdRightNow {
		isBanned := s.CheckIfUserIsBanned(ctx, userId)
		if isBanned {
			return auth.LoginResponse{
				Status: "failed",
				Reason: "banned",
				Meta: map[string]string{
					"on":      time.Now().String(),
					"unlocks": time.Now().Add(time.Hour * 72).String(),
					"reason":  "cheating",
					"by":      "system",
				},
			}, nil
		}
	}

	// 4. Generate system access tokens
	tokenPair, err := auth.GetJwtService().GenerateTokenPair(ctx, userId.String(), "player")
	if err != nil {
		// Fix: Join the error so internal telemetry logs the root cause
		s.logger.Error("Token generation failed", zap.Error(err))
		return auth.LoginResponse{
			Status: "failed",
			Reason: err.Error(),
		}, errors.Join(auth.ErrLoginFailed, err)
	}

	tokens := auth.SuccessfulResponse{
		AccessToken:  tokenPair.AccessToken,
		TokenType:    "Bearer",
		ExpiresIn:    time.Now().Add(time.Duration(s.config.AuthToken.AccessToken.ExpiresIn)), // Fixed cast typo
		RefreshToken: tokenPair.RefreshToken,
		Scope:        "", //s.config.AuthToken.AccessToken.Scopes
	}

	response := auth.LoginResponse{
		Status: "success",
		Tokens: &tokens,
	}
	if createdRightNow {
		response.Status = "created"
	}

	return response, nil
}

func (s *LocalUsersService) Refresh(ctx context.Context, dto auth.TokenRefreshRequest) (auth.TokenRefreshResponse, error) {
	tokenPair, err := auth.GetJwtService().RefreshTokens(ctx, dto.RefreshToken, "player")
	if err != nil {
		s.logger.Error("Failed to generate token pair for refreshing tokens", zap.Error(err))

		return auth.TokenRefreshResponse{
			Failed: &auth.FailedRefreshResponse{
				Error:            err.Error(),
				ErrorDescription: err.Error(),
			},
		}, err
	}

	return auth.TokenRefreshResponse{
		Successful: &auth.SuccessfulResponse{
			AccessToken:  tokenPair.AccessToken,
			TokenType:    "Bearer",
			ExpiresIn:    time.Now().Add(time.Duration(s.config.AuthToken.AccessToken.ExpiresIn)),
			RefreshToken: tokenPair.RefreshToken,
			// Scope: ,
		},
	}, nil
}

func (s *LocalUsersService) GetBySocialProfileId(ctx context.Context, provider, socialUserId string) (User, error) {
	ctx, span := tracer.Start(ctx, "users.service.getbysocialprofileid")
	defer span.End()

	return s.repo.GetBySocialProfileId(ctx, provider, socialUserId)
}

func (s *LocalUsersService) FindOrCreateUserFromSocialProfile(ctx context.Context, profile auth.UnifiedUserProfile) (uuid.UUID, bool, error) {
	ctx, span := tracer.Start(ctx, "users.service.findorcreateuserfromsocialprofile")
	defer span.End()

	return s.repo.FindOrCreateUserFromSocialProfile(ctx, profile)
}

func (s *LocalUsersService) CheckIfUserIsBanned(ctx context.Context, userId uuid.UUID) bool {
	ctx, span := tracer.Start(ctx, "users.service.checkifuserisbanned")
	defer span.End()

	return s.repo.CheckIfUserIsBanned(ctx, userId)
}

func (s *LocalUsersService) BulkAdd(ctx context.Context, dto BulkCreateDTO) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "users.service.bulkadd")
	defer span.End()

	return s.repo.BulkAdd(ctx, dto)
}

func (s *LocalUsersService) BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "users.service.bulkmodify")
	defer span.End()

	return s.repo.BulkModify(ctx, dto)
}

func (s *LocalUsersService) BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "users.service.bulkdelete")
	defer span.End()

	return s.repo.BulkDelete(ctx, dto)
}
