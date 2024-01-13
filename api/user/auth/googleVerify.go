package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type JWKS struct {
	Keys []JSONWebKey `json:"keys"`
}

// Struct for parsing JSON Web Key
type JSONWebKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func fetchJWKS(jwksURL string) (*JWKS, error) {
	response, err := http.Get(jwksURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var jwks JWKS
	err = json.NewDecoder(response.Body).Decode(&jwks)
	if err != nil {
		return nil, err
	}

	return &jwks, nil
}

func verifyIDToken(idToken string) (bool, error) {
	jwksURL := "https://www.googleapis.com/oauth2/v3/certs"

	// public key from Google
	jwks, err := fetchJWKS(jwksURL)
	if err != nil {
		return false, err
	}

	token, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		// find public key for token
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("Missing key ID (kid) in token header")
		}

		for _, key := range jwks.Keys {
			if key.Kid == kid {
				// pub key constr using N and E
				modulus, err := base64.RawURLEncoding.DecodeString(key.N)
				if err != nil {
					return nil, err
				}

				exponent, err := base64.RawURLEncoding.DecodeString(key.E)
				if err != nil {
					return nil, err
				}

				return &rsa.PublicKey{
					N: new(big.Int).SetBytes(modulus),
					E: int(new(big.Int).SetBytes(exponent).Int64()),
				}, nil
			}
		}
		return nil, errors.New("Unable to find matching key")
	})

	if err != nil {
		return false, err
	}

	return token.Valid, nil
}
