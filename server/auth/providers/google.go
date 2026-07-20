package providers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/abhinash-kml/nova/server/auth"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

type GoogleProvider struct {
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
}

func NewGoogleProvider(ctx context.Context, clientID, clientSecret, redirectUri string, scopes []string) (*GoogleProvider, error) {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, err
	}

	return &GoogleProvider{
		oauth2Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     endpoints.Google,
			RedirectURL:  redirectUri,
			Scopes:       scopes,
		},
		verifier: provider.Verifier(
			&oidc.Config{ClientID: clientID},
		),
	}, nil
}

func (g GoogleProvider) Name() string {
	return "google"
}

func (g GoogleProvider) AuthURL() string {
	return g.oauth2Config.AuthCodeURL("meow", oauth2.AccessTypeOffline)
}

func (g *GoogleProvider) AuthenticateWithCode(ctx context.Context, code string) (auth.TokenResponse, error) {
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
	client := g.oauth2Config.Client(ctx, &token.Token)
	res, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return auth.UnifiedUserProfile{}, err
	}
	defer res.Body.Close()

	googleProfile := auth.GoogleProfile{}
	json.NewDecoder(res.Body).Decode(&googleProfile)

	profile := auth.NormalizeGoogle(googleProfile)
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
