package providers

import (
	"context"
	"fmt"

	"github.com/abhinash-kml/nova/server/auth"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type GoogleProvider struct {
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
}

func NewGoogleProvider(ctx context.Context, clientID string) (*GoogleProvider, error) {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, err
	}

	return &GoogleProvider{
		verifier: provider.Verifier(
			&oidc.Config{ClientID: clientID},
		),
	}, nil
}

func (g GoogleProvider) Name() string {
	return "google"
}

func (g *GoogleProvider) Authenticate(ctx context.Context, code string) (auth.TokenResponse, error) {
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

func (g *GoogleProvider) GetProfile(ctx context.Context, token auth.TokenResponse) (auth.UnifiedUserProfile, error) {
	// Make request to provider endpoint to get profile info
	// Main fields = id, name, email, avatar_url

	profile := auth.UnifiedUserProfile{}
	return profile, nil
}

func (g *GoogleProvider) verifyIdToken(ctx context.Context, idToken string) (*auth.GoogleProfile, error) {
	// Verify offline
	response, err := g.verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("invalid google id token signature: %w", err)
	}

	var claims auth.GoogleProfile
	if err := response.Claims(&claims); err != nil {
		return nil, err
	}

	return &claims, nil
}
