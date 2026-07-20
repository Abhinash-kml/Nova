package auth

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrLoginFailed           = errors.New("login failed")
	ErrProviderNotExist      = errors.New("provider doesn't exist")
	ErrProviderNotConfigured = errors.New("provider is not configured")
	ErrFailedToExchangeCode  = errors.New("failed to exchange auth code for token")
)

var (
	eonce  sync.Once
	engine *SocialAuthEngine
)

type Provider interface {
	// Get name of provider
	Name() string

	// Exchange auth code for token
	AuthenticateWithCode(ctx context.Context, code string) (TokenResponse, error)

	// Get unified profile from provider
	GetProfile(ctx context.Context, token TokenResponse) (UnifiedUserProfile, error)

	// Get generated Auth URl
	AuthURL() string
}

type SocialAuthEngine struct {
	providers map[string]Provider
	mu        sync.Mutex
}

func NewSocialAuthEngine() *SocialAuthEngine {
	return &SocialAuthEngine{
		providers: make(map[string]Provider),
	}
}

func (s *SocialAuthEngine) Register(p Provider) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.providers[p.Name()] = p
}

func (s *SocialAuthEngine) CompleteAuthentication(ctx context.Context, providerName string, code string) (UnifiedUserProfile, error) {
	provider, exists := s.providers[providerName]

	if !exists {
		return UnifiedUserProfile{}, ErrProviderNotExist
	}

	// Exchange routing payload code for base tokens
	tokenResponse, err := provider.AuthenticateWithCode(ctx, code)
	if err != nil {
		return UnifiedUserProfile{}, ErrFailedToExchangeCode
	}

	// Delegate the profile gathering to the plugin (OIDC or OAuth2 server)
	return provider.GetProfile(ctx, tokenResponse)
}

// Should only be called once
func SetDefaultAuthEngine(e *SocialAuthEngine) {
	eonce.Do(func() {
		engine = e
	})
}

// Returns the global auth engine, sets up one if its nil
func GetSocialAuthEngine() *SocialAuthEngine {
	if engine == nil {
		engine = &SocialAuthEngine{
			providers: make(map[string]Provider, 10),
		}
	}

	return engine
}
