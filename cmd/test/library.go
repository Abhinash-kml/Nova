package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

var (
	Token UnifiedToken
)

type UnifiedToken struct {
	oauth2.Token
	IdToken string `json:"id_token"`
	Scope   string `json:"scope"`
}

type ExtraClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type IdTokenVerificationResponse struct {
	oidc.IDToken
	ExtraClaims
}

type UnifiedResult struct {
	TokenResponse        UnifiedToken                `json:"TOKEN_RESPONSE"`
	VerificationResponse IdTokenVerificationResponse `json:"ID_TOKEN_VERIFICATION_RESPONSE"`
}

func PerformOAuthGoogleLibrary() {
	config := oauth2.Config{
		ClientID:     CLIENT_ID,
		ClientSecret: CLIENT_SECRET,
		Endpoint:     endpoints.Google,
		Scopes: []string{
			"openid",
			"profile",
			"email",
			// "https://www.googleapis.com/auth/userinfo.profile",
			// "https://www.googleapis.com/auth/userinfo.email",
			// "https://www.googleapis.com/auth/user.gender.read",
			// "https://www.googleapis.com/auth/user.phonenumbers.read",
		},
		RedirectURL: REDIRECT_URI,
	}

	fmt.Println("authcode url:", config.AuthCodeURL("meow-meow", oauth2.AccessTypeOffline, oauth2.ApprovalForce))

	go ListenOnRedirectURICustom(&config)
}

func ListenOnRedirectURICustom(config *oauth2.Config) {
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		code := query.Get("code")
		_ = query.Get("state")

		if code == "" {
			panic("No code found")
		}

		token, err := config.Exchange(context.Background(), code, oauth2.AccessTypeOffline)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		Token = UnifiedToken{
			Token:   *token,
			IdToken: token.Extra("id_token").(string),
			Scope:   token.Extra("scope").(string),
		}

		provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
		if err != nil {
			fmt.Println("Failed to create oidc verifier")
		}

		oidcConfig := oidc.Config{
			ClientID: CLIENT_ID,
		}
		verifier := provider.Verifier(&oidcConfig)
		idToken, err := verifier.Verify(context.Background(), Token.IdToken)
		if err != nil {
			fmt.Println("Failed to verify id token")
			os.Exit(1)
		}

		var claims ExtraClaims
		idToken.Claims(&claims)

		response := UnifiedResult{
			TokenResponse: Token,
			VerificationResponse: IdTokenVerificationResponse{
				IDToken:     *idToken,
				ExtraClaims: claims,
			},
		}

		encoder := json.NewEncoder(w)
		encoder.SetIndent(" ", "  ")
		encoder.Encode(response)

	})

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
}
