package users

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/abhinash-kml/nova/server/auth"
	"github.com/abhinash-kml/nova/server/cache"
	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	// General operations
	Add(ctx context.Context, req CreateRequest) (User, error)
	GetAll(ctx context.Context, cursor uuid.UUID, limit int, sortBy, sortDirection string) ([]User, error)
	GetAllByAttribute(ctx context.Context, attribute string) ([]User, error)
	GetById(ctx context.Context, id uuid.UUID) (User, error)
	GetByName(ctx context.Context, name string) (User, error)
	GetBySocialProfileId(ctx context.Context, provider, socialUserId string) (User, error)
	Update(ctx context.Context, req UpdateRequest) (User, error)
	Replace(ctx context.Context, req ReplaceRequest) (User, error)
	Delete(ctx context.Context, req DeleteRequest) (uuid.UUID, error)

	// Auth operations
	Login(ctx context.Context, loginRequest auth.LoginRequest) (auth.LoginResponse, error)
	Refresh(ctx context.Context, req auth.TokenRefreshRequest) (auth.TokenRefreshResponse, error)

	// Operations
	FindOrCreateUserFromSocialProfile(ctx context.Context, profile auth.UnifiedUserProfile) (uuid.UUID, bool, error)
	CheckIfUserIsBanned(ctx context.Context, userId uuid.UUID) (bool, error)

	Ban(ctx context.Context, req BanRequest) (BannedUserResponse, error)
	UnBan(ctx context.Context, req UnbanRequest) error

	// Bulk operations
	BulkAdd(ctx context.Context, req BulkCreateRequest) ([]common.BulkOpResult, error)
	BulkModify(ctx context.Context, req BulkModifyRequest) ([]common.BulkOpResult, error)
	BulkDelete(ctx context.Context, req BulkDeleteRequest) ([]common.BulkOpResult, error)

	CheckUsername(ctx context.Context, req CheckUsernameRequest) (bool, error)
	GetAllUsernames(ctx context.Context, req UserId) ([]string, error)
}

type LocalUsersService struct {
	repo   UsersRepository
	logger *zap.Logger
	cache  *cache.TieredCache
	config *config.UsersConfig
}

func NewLocalUsersService(repository UsersRepository, cache *cache.TieredCache, c *config.UsersConfig, l *zap.Logger) *LocalUsersService {
	return &LocalUsersService{
		repo:   repository,
		cache:  cache,
		logger: l,
		config: c,
	}
}

func (s *LocalUsersService) Add(ctx context.Context, req CreateRequest) (User, error) {
	ctx, span := tracer.Start(ctx, "users.service.add")
	defer span.End()

	createdUser, err := s.repo.Add(ctx, req)
	if err != nil {
		return User{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to create user")
	}

	return createdUser, nil
}

func (s *LocalUsersService) GetAll(ctx context.Context, cursor uuid.UUID, limit int, sortBy, sortDirection string) ([]User, error) {
	ctx, span := tracer.Start(ctx, "users.service.getall")
	defer span.End()

	users, err := s.repo.GetAll(ctx, cursor, limit, sortBy, sortDirection)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all users")
	}

	return users, nil
}

func (s *LocalUsersService) GetAllByAttribute(ctx context.Context, attribute string) ([]User, error) {
	ctx, span := tracer.Start(ctx, "users.service.getallbyattribute")
	defer span.End()

	users, err := s.repo.GetAllByAttribute(ctx, attribute)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all users by attribute")
	}

	return users, nil
}

func (s *LocalUsersService) GetById(ctx context.Context, id uuid.UUID) (User, error) {
	ctx, span := tracer.Start(ctx, "users.service.getbyid")
	defer span.End()

	key := s.config.KeyPrefix + id.String()

	// Get from cache
	raw, err := s.cache.GetOrLoadWithTTL(ctx, key, s.config.CacheDuration, func(ctx context.Context) ([]byte, error) {
		user, err := s.repo.GetById(ctx, id)
		if err != nil {
			return nil, err
		}

		return json.Marshal(user)
	})
	if err != nil {
		return User{}, err
	}

	var user User
	err = json.Unmarshal(raw, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *LocalUsersService) GetByName(ctx context.Context, name string) (User, error) {
	ctx, span := tracer.Start(ctx, "users.service.getbyname")
	defer span.End()

	key := s.config.KeyPrefix + name

	// Get from cache or load into it from db
	raw, err := s.cache.GetOrLoadWithTTL(ctx, key, s.config.CacheDuration, func(ctx context.Context) ([]byte, error) {
		user, err := s.repo.GetByName(ctx, name)
		if err != nil {
			return nil, err
		}

		return json.Marshal(user)
	})
	if err != nil {
		return User{}, err
	}

	var user User
	err = json.Unmarshal(raw, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *LocalUsersService) Update(ctx context.Context, req UpdateRequest) (User, error) {
	ctx, span := tracer.Start(ctx, "users.service.update")
	defer span.End()

	// Update repository first
	user, err := s.repo.Update(ctx, req)
	if err != nil {
		return User{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to update user")
	}

	// Invalidate old record from cache, next get call with repopulate it
	key := s.config.KeyPrefix + req.Id
	go s.cache.Delete(ctx, key)

	return user, nil
}

func (s *LocalUsersService) Replace(ctx context.Context, req ReplaceRequest) (User, error) {
	ctx, span := tracer.Start(ctx, "users.service.replace")
	defer span.End()

	// Update repository first
	user, err := s.repo.Replace(ctx, req)
	if err != nil {
		return User{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to replace user")
	}

	// Invalidate old record from cache, next get call with repopulate it
	key := s.config.KeyPrefix + req.Id
	go s.cache.Delete(ctx, key)

	return user, nil
}

func (s *LocalUsersService) Delete(ctx context.Context, req DeleteRequest) (uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "users.service.delete")
	defer span.End()

	// Delete in repo
	deletedId, err := s.repo.Delete(ctx, req)
	if err != nil {
		return uuid.Nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to delete user")
	}

	// Delete from cache
	key := s.config.KeyPrefix + req.Id
	go s.cache.Delete(ctx, key)

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
		isBanned, err := s.CheckIfUserIsBanned(ctx, userId)
		if err != nil {
			return auth.LoginResponse{}, err
		}
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
		ExpiresIn:    time.Now().Add(time.Duration(config.GetInstance().AuthToken.AccessToken.ExpiresIn)), // Fixed cast typo
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

func (s *LocalUsersService) Refresh(ctx context.Context, req auth.TokenRefreshRequest) (auth.TokenRefreshResponse, error) {
	tokenPair, err := auth.GetJwtService().RefreshTokens(ctx, req.RefreshToken, "player")
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
			ExpiresIn:    time.Now().Add(time.Duration(config.GetInstance().AuthToken.AccessToken.ExpiresIn)),
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

func (s *LocalUsersService) CheckIfUserIsBanned(ctx context.Context, userId uuid.UUID) (bool, error) {
	ctx, span := tracer.Start(ctx, "users.service.checkifuserisbanned")
	defer span.End()

	return s.repo.CheckIfUserIsBanned(ctx, userId)
}

func (s *LocalUsersService) BulkAdd(ctx context.Context, req BulkCreateRequest) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "users.service.bulkadd")
	defer span.End()

	results, err := s.repo.BulkAdd(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to bulk add users")
	}

	return results, nil
}

func (s *LocalUsersService) BulkModify(ctx context.Context, req BulkModifyRequest) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "users.service.bulkmodify")
	defer span.End()

	results, err := s.repo.BulkModify(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to bulk modify users")
	}

	return results, nil
}

func (s *LocalUsersService) BulkDelete(ctx context.Context, req BulkDeleteRequest) ([]common.BulkOpResult, error) {
	ctx, span := tracer.Start(ctx, "users.service.bulkdelete")
	defer span.End()

	results, err := s.repo.BulkDelete(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to bulk delete users")
	}

	return results, nil
}

func (s *LocalUsersService) Ban(ctx context.Context, req BanRequest) (BannedUserResponse, error) {
	ctx, span := tracer.Start(ctx, "users.service.ban")
	defer span.End()

	bannedUser, err := s.repo.Ban(ctx, req)
	if err != nil {
		return BannedUserResponse{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to ban user")
	}

	return bannedUser, nil
}

func (s *LocalUsersService) UnBan(ctx context.Context, req UnbanRequest) error {
	ctx, span := tracer.Start(ctx, "users.service.unban")
	defer span.End()

	err := s.repo.UnBan(ctx, req)
	if err != nil {
		return common.TranslatePostgresError(err, s.logger).WithMessage("Failed to unban user")
	}

	return nil
}

func (s *LocalUsersService) CheckUsername(ctx context.Context, req CheckUsernameRequest) (bool, error) {
	ctx, span := tracer.Start(ctx, "users.service.checkusername")
	defer span.End()

	valid, err := s.repo.CheckUsername(ctx, req)
	if err != nil {
		return false, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to check username")
	}

	return valid, nil
}

func (s *LocalUsersService) GetAllUsernames(ctx context.Context, req UserId) ([]string, error) {
	ctx, span := tracer.Start(ctx, "users.service.getallusername")
	defer span.End()

	usernames, err := s.repo.GetAllUsernames(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get all usernames")
	}

	return usernames, nil
}
