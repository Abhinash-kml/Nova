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

type FacebookProvider struct {
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
	logger       *zap.Logger
}

func NewFacebookProvider(ctx context.Context, l *zap.Logger, clientID, clientSecret, redirectUri string, scopes []string) (*FacebookProvider, error) {
	provider, err := oidc.NewProvider(ctx, "https://www.facebook.com")
	if err != nil {
		return nil, err
	}

	return &FacebookProvider{
		logger: l,
		oauth2Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     endpoints.Facebook,
			RedirectURL:  redirectUri,
			Scopes:       scopes,
		},
		verifier: provider.Verifier(
			&oidc.Config{ClientID: clientID},
		),
	}, nil
}

func (p FacebookProvider) Name() string {
	return "facebook"
}

func (p FacebookProvider) AuthURL() string {
	return p.oauth2Config.AuthCodeURL("meow", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func (p *FacebookProvider) AuthenticateWithCode(ctx context.Context, code string) (auth.TokenResponse, error) {
	ctx, span := tracer.Start(ctx, "auth.provider.facebook.authenticatewithcode")
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

func (p *FacebookProvider) GetProfile(ctx context.Context, token auth.TokenResponse) (auth.UnifiedUserProfile, error) {
	ctx, span := tracer.Start(ctx, "auth.provider.facebook.getprofile")
	defer span.End()

	// Make request to provider endpoint to get profile info
	// Main fields = id, name, email, avatar_url
	client := p.oauth2Config.Client(ctx, &token.Token)
	res, err := client.Get("https://graph.facebook.com/v22.0/me?fields=id,name,first_name,last_name,email,picture")
	if err != nil {
		return auth.UnifiedUserProfile{}, err
	}
	defer res.Body.Close()

	facebookProfile := auth.FacebookProfile{}
	json.NewDecoder(res.Body).Decode(&facebookProfile)

	p.logger.Info("Facebook Profile", zap.String("id", facebookProfile.ID),
		zap.String("name", facebookProfile.Name),
		zap.String("first name", facebookProfile.FirstName),
		zap.String("last name", facebookProfile.LastName),
		zap.String("email", facebookProfile.Email),
		zap.String("avatar url", facebookProfile.Picture.Data.URL))

	profile := auth.NormalizeFacebook(ctx, facebookProfile)
	return profile, nil
}

func (p *FacebookProvider) verifyIdToken(ctx context.Context, idToken string) (*auth.GoogleProfile, error) {
	// Verify offline
	response, err := p.verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("invalid facebook id token signature: %w", err)
	}

	var claims auth.GoogleProfile
	if err := response.Claims(&claims); err != nil {
		return nil, err
	}

	return &claims, nil
}
