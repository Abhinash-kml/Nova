package providers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/abhinash-kml/nova/server/auth"
	"github.com/coreos/go-oidc/v3/oidc"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

type GoogleProvider struct {
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
	logger       *zap.Logger
}

func NewGoogleProvider(ctx context.Context, l *zap.Logger, clientID, clientSecret, redirectUri string, scopes []string) (*GoogleProvider, error) {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, err
	}

	return &GoogleProvider{
		logger: l,
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

func (p GoogleProvider) Name() string {
	return "google"
}

func (p GoogleProvider) AuthURL() string {
	return p.oauth2Config.AuthCodeURL("meow", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func (p *GoogleProvider) AuthenticateWithCode(ctx context.Context, code string) (auth.TokenResponse, error) {
	ctx, span := tracer.Start(ctx, "auth.provider.google.authenticatewithcode")
	defer span.End()

	t, err := p.oauth2Config.Exchange(ctx, code, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	if err != nil {
		return auth.TokenResponse{}, err
	}

	response := auth.TokenResponse{
		Token:   *t,
		IdToken: t.Extra("id_token").(string),
		Scopes:  auth.ExtractScopesStr(t),
	}

	return response, nil
}

func (p *GoogleProvider) GetProfile(ctx context.Context, token auth.TokenResponse) (auth.UnifiedUserProfile, error) {
	ctx, span := tracer.Start(ctx, "auth.provider.google.getprofile")
	defer span.End()

	// Make request to provider endpoint to get profile info
	// Main fields = id, name, email, avatar_url
	client := p.oauth2Config.Client(ctx, &token.Token)
	res, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return auth.UnifiedUserProfile{}, err
	}
	defer res.Body.Close()

	googleProfile := auth.GoogleProfile{}
	json.NewDecoder(res.Body).Decode(&googleProfile)

	p.logger.Debug("Google Profile", zap.String("sub", googleProfile.Sub),
		zap.String("name", googleProfile.Name),
		zap.String("given_name", googleProfile.GivenName),
		zap.String("family_name", googleProfile.FamilyName),
		zap.String("email", googleProfile.Email),
		zap.Bool("email_verified", googleProfile.EmailVerified),
		zap.String("locale", googleProfile.Locale))

	profile := auth.NormalizeGoogle(ctx, googleProfile)
	return profile, nil
}

func (p *GoogleProvider) verifyIdToken(ctx context.Context, idToken string) (*auth.GoogleProfile, error) {
	// Verify offline
	response, err := p.verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("invalid google id token signature: %w", err)
	}

	var claims auth.GoogleProfile
	if err := response.Claims(&claims); err != nil {
		return nil, err
	}

	return &claims, nil
}
