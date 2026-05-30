package auth

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/oauth2"
)

var (
	ErrProviderNotExist      = errors.New("provider doesn't exist")
	ErrProviderNotConfigured = errors.New("provider is not configured")
	ErrFailedToExchangeCode  = errors.New("failed to exchange auth code for token")
)

type UnifiedProfile struct {
	Provider  string
	Id        string
	Name      string
	Email     string
	AvatalUrl string
}

type UnifiedToken struct {
	Type string
	Raw  string
}

type Provider interface {
	Name() string
	ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error)
	GetProfile(ctx context.Context, token *oauth2.Token) (*UnifiedProfile, error)
	ValidateAndFetch(ctx context.Context, token UnifiedToken) (*UnifiedProfile, error)
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

func (s *SocialAuthEngine) CompleteAuthWithCode(ctx context.Context, providerName string, code string) (*UnifiedProfile, error) {
	s.mu.Lock()
	provider, exists := s.providers[providerName]
	s.mu.Unlock()

	if !exists {
		return nil, ErrProviderNotExist
	}

	// Exchange routing payload code for base tokens
	token, err := provider.ExchangeCode(ctx, code)
	if err != nil {
		return nil, ErrFailedToExchangeCode
	}

	// Delegate the profile gathering to the plugin (OIDC or OAuth2 server)
	return provider.GetProfile(ctx, token)
}

func (s *SocialAuthEngine) CompleteAuthWithToken(ctx context.Context, providerName string, token UnifiedToken) (*UnifiedProfile, error) {
	s.mu.Lock()
	provider, exists := s.providers[providerName]
	s.mu.Unlock()

	if !exists {
		return nil, ErrProviderNotExist
	}

	return provider.ValidateAndFetch(ctx, token)
}
