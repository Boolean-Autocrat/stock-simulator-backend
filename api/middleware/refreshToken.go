package middleware

import (
	"context"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func refreshAccessToken(refreshToken *oauth2.Token) (*oauth2.Token, error) {
	// Create a new OAuth2 configuration with the same parameters as googleOauthConfig
	conf := &oauth2.Config{
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}

	tokenSource := conf.TokenSource(context.Background(), refreshToken)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	return newToken, nil
}
