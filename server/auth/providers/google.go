package providers

import (
	"context"
	"fmt"

	"github.com/abhinash-kml/nova/server/auth"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Claims struct {
	Subject string `json:"sub"`
	Name    string `json:"name"`
	Email   string `json:"email"`
}

type Google struct {
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
}

func NewGoogle(ctx context.Context, clientID string) (*Google, error) {
	provider, err := oidc.NewProvider(ctx, "https://google.com")
	if err != nil {
		return nil, err
	}
	return &Google{
		verifier: provider.Verifier(
			&oidc.Config{ClientID: clientID},
		),
	}, nil
}

func (g Google) Name() string {
	return "google"
}

func (g *Google) ExchangeAuthCodeForToken(ctx context.Context, code string) (auth.TokenResponse, error) {
	t, err := g.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return auth.TokenResponse{}, err
	}

	response := auth.TokenResponse{
		Token:   *t,
		IdToken: t.Extra("id_token").(string),
		Scopes:  t.Extra("scope").(string),
	}

	return response, nil
}

func (g *Google) GetProfile(ctx context.Context, token auth.TokenResponse) (*auth.UnifiedProfile, error) {
	// Make request to provider endpoint to get profile info
	// Main fields = id, name, email, avatar_url

	profile := auth.AutheticatedUserProfile{}
	return profile, nil
}

func (g *Google) ValidateAndFetch(ctx context.Context, t auth.UnifiedToken) (*auth.UnifiedProfile, error) {
	if t.Type != "id" {
		return nil, fmt.Errorf("google login requires an ide token")
	}

	claims, err := g.verifyIdToken(ctx, t.Raw)
	if err != nil {
		return nil, fmt.Errorf("google id token verfication failed: %w", err)
	}

	return &auth.UnifiedProfile{
		Provider: g.Name(),
		Id:       claims.Subject,
		Name:     claims.Name,
		Email:    claims.Email,
	}, nil

}

func (g *Google) verifyIdToken(ctx context.Context, t string) (*Claims, error) {
	// Verify offline
	token, err := g.verifier.Verify(ctx, t)
	if err != nil {
		return nil, fmt.Errorf("invalid google id token signature: %w", err)
	}

	var claims Claims
	if err := token.Claims(&claims); err != nil {
		return nil, err
	}

	return &claims, nil
}
