// Copyright (c) HashiCorp, Inc.

package oidc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Client struct {
	// Token endpoint for the OIDC provider.
	tokenEndpoint string
	// The client ID used to authenticate with the OIDC provider.
	clientID string
	// The client secret used to authenticate with the OIDC provider.
	clientSecret string
	// The email address used to authenticate with the OIDC provider.
	email string

	// The access token used to access the Disruptive REST API.
	token *Token
}

type Token struct {
	// The access token used to access the Disruptive REST API.
	accessToken string
	tokenType   string
	expiry      time.Time
	mu          sync.RWMutex
}

func (t *Token) get() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.accessToken
}

func (t *Token) set(token, tokenType string, expiry time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.accessToken = token
	t.expiry = expiry
	t.tokenType = tokenType
}

func (t *Token) isExpired() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.expiry.Before(time.Now())
}

type Config struct {
	// Token endpoint for the OIDC provider.
	TokenEndpoint string
	// The client ID used to authenticate with the OIDC provider.
	ClientID string
	// The client secret used to authenticate with the OIDC provider.
	ClientSecret string
	// The email address used to authenticate with the OIDC provider.
	Email string
}

func NewClient(cfg Config) *Client {
	return &Client{
		tokenEndpoint: cfg.TokenEndpoint,
		clientID:      cfg.ClientID,
		clientSecret:  cfg.ClientSecret,
		email:         cfg.Email,
		token:         &Token{},
	}
}

func (c *Client) createJWT() (string, error) {
	// Construct the JWT header.
	jwtHeader := map[string]interface{}{
		"alg": "HS256",
		"kid": c.clientID,
	}

	// Construct the JWT payload.
	now := time.Now()
	jwtPayload := &jwt.RegisteredClaims{
		Issuer:    c.email,
		Audience:  jwt.ClaimStrings{c.tokenEndpoint},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
	}

	// Sign and encode JWT with the secret.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtPayload)
	token.Header = jwtHeader
	encodedJwt, err := token.SignedString([]byte(c.clientSecret))
	if err != nil {
		return "", err
	}

	return encodedJwt, nil
}

func (c *Client) GetToken(ctx context.Context) (*AuthResponse, error) {
	// Check if we already have a valid cached token, if so, return it.
	if cachedToken := c.token.get(); cachedToken != "" && !c.token.isExpired() {
		tflog.Debug(ctx, "using cached token")
		return &AuthResponse{
			AccessToken: cachedToken,
			TokenType:   c.token.tokenType,
			ExpiresIn:   int(time.Until(c.token.expiry).Seconds()),
		}, nil
	}

	jwt, err := c.createJWT()
	if err != nil {
		return nil, fmt.Errorf("oidc: failed to create JWT: %w", err)
	}

	// Prepare HTTP POST request data.
	reqData := url.Values{
		"assertion":  {jwt},
		"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
	}.Encode()

	ctx = tflog.SetField(ctx, "jwt", jwt)
	ctx = tflog.SetField(ctx, "token_endpoint", c.tokenEndpoint)
	ctx = tflog.SetField(ctx, "method", http.MethodPost)
	ctx = tflog.SetField(ctx, "body", reqData)

	tflog.Debug(ctx, "sending request to OIDC provider")

	// Create the request to exchange the JWT for an access token.
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.tokenEndpoint, strings.NewReader(reqData))
	if err != nil {
		return nil, fmt.Errorf("oidc: failed to create request: %w", err)
	}

	// Set Content-Type header to specify that our body is Form-URL Encoded.
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Exchange the JWT for an access token. Set a 3 second
	// timeout in case the server can't be reached.
	httpClient := &http.Client{Timeout: time.Second * 3}
	response, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("oidc: failed to send request: %w", err)
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("oidc: failed to read response body: %w, status: %d", err, response.StatusCode)
	}
	ctx = tflog.SetField(ctx, "body", string(bodyBytes))
	ctx = tflog.SetField(ctx, "status_code", response.StatusCode)
	if response.StatusCode != http.StatusOK {
		tflog.Debug(ctx, "received non-200 status code from DT API")
		return nil, fmt.Errorf("HTTP error: %d: %s", response.StatusCode, string(bodyBytes))
	}

	// Decode the response body to an AuthResponse.
	var authResponse *AuthResponse
	err = json.Unmarshal(bodyBytes, &authResponse)
	if err != nil {
		tflog.Debug(ctx, "failed to unmarshal response body")
		return nil, fmt.Errorf("oidc: failed to unmarshal response body: %w", err)
	}
	// Set the token in the client cache.
	expiry := time.Now().Add(time.Duration(authResponse.ExpiresIn) * time.Second)
	c.token.set(authResponse.AccessToken, authResponse.TokenType, expiry)

	return authResponse, nil

}

type AuthResponse struct {
	// The access token used to access the Disruptive REST API.
	AccessToken string `json:"access_token"`
	// The type of token this is. Will typically be "Bearer".
	TokenType string `json:"token_type"`
	// How many seconds until the token expires. Will typically be 3600.
	ExpiresIn int `json:"expires_in"`
}
