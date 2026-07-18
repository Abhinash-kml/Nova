package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var GOT_CODE string

func PerformOAuthGoogleCustom() {
	loginURI := ConstructLoginURI(CLIENT_ID, SCOPE)
	fmt.Println(loginURI)

	go ListenOnRedirectURI()
}

func ConstructLoginURI(clientId, scope string) string {
	base := "https://accounts.google.com/o/oauth2/v2/auth"
	u, err := url.Parse(base)
	if err != nil {
		panic("Failed to parse base url")
	}

	q := u.Query()
	q.Set("prompt", "select_account consent")
	q.Set("response_type", "code")
	q.Set("scope", scope)
	q.Set("access_type", "offline")
	q.Set("include_granted_scopes", "true")
	q.Set("state", "meow-meow")
	q.Set("redirect_uri", "http://localhost:8000/token")
	q.Set("client_id", "853460492110-oj3lstd9g819rupfduv3588rm308a1bn.apps.googleusercontent.com")
	u.RawQuery = q.Encode()

	return u.String()
}

func ConstructTokenExchangeURI(clientId, code string) string {
	base := "https://oauth2.googleapis.com/token"
	u, err := url.Parse(base)
	if err != nil {
		panic("Failed to parse token exchange uri")
	}

	// q := u.Query()
	// q.Set("client_id", clientId)
	// q.Set("code", code)
	// q.Set("grant_type", "authorization_code")
	// q.Set("redirect_uri", REDIRECT_URI)
	// u.RawQuery = q.Encode()

	return u.String()
}

func ListenOnRedirectURI() {
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		code, ok := query["code"]
		if !ok {
			fmt.Println("No code found")
			return
		} else {
			fmt.Println("Code:", code[0])
			GOT_CODE = code[0]
		}

		w.Write(ExchangeCodeForToken(CLIENT_ID, GOT_CODE))
	})

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		panic("Failed to start api server")
	}
}

func ExchangeCodeForToken(clientId, code string) []byte {
	exchangeURI := ConstructTokenExchangeURI(clientId, code)
	fmt.Println("Exchange URI:", exchangeURI)

	v := url.Values{}
	v.Set("client_id", clientId)
	v.Set("client_secret", CLIENT_SECRET)
	v.Set("code", code)
	v.Set("grant_type", "authorization_code")
	v.Set("redirect_uri", REDIRECT_URI)

	client := http.DefaultClient
	res, err := client.Post(exchangeURI, "application/x-www-form-urlencoded", strings.NewReader(v.Encode()))
	if err != nil {
		panic("Failed to send request to exchange code for token")
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic("Failed to read body of response")
	}

	fmt.Println("Token --")
	fmt.Println(string(body))

	return body
}
