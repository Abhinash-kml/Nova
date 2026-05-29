package auth

import (
	"context"
	"sync"

	"golang.org/x/oauth2"
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
	ValidateAndFetch(ctx context.Context, token UnifiedToken)
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
	return nil, nil
}

func (s *SocialAuthEngine) CompleteAuthWithToken(token UnifiedToken) (*UnifiedProfile, error) {
	return nil, nil
}
