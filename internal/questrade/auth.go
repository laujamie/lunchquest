package questrade

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

type authenticateResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int32  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	ApiServer    string `json:"api_server"`
}

type customOauth2Token struct {
	oauth2.Token
	ApiServer string `json:"api_server,omitempty"`
}

const questradeLoginEndpoint = "https://login.questrade.com/oauth2/token"

func Authenticate(ctx context.Context, refreshToken string) (*customOauth2Token, error) {
	// Check if we previously already have a valid access token
	storedAuth, err := keyring.Get("lunchquest-cli", "authenticate-response")
	if err == nil {
		b := []byte(storedAuth)
		var storedAuthToken customOauth2Token
		if err := json.Unmarshal(b, &storedAuthToken); err == nil && storedAuthToken.Expiry.After(time.Now()) {
			return nil, nil
		}
	}

	// Get auth token
	qc := oauth2.NewClient(ctx, nil)

	req, err := http.NewRequest(http.MethodGet, questradeLoginEndpoint, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("refresh_token", refreshToken)
	q.Add("grant_type", "refresh_token")

	req.URL.RawQuery = q.Encode()

	req.Header.Add(
		"User-Agent",
		viper.GetString("UserAgent"),
	)

	resp, err := qc.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var res authenticateResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	token := &customOauth2Token{
		Token: oauth2.Token{
			AccessToken: res.AccessToken,
			TokenType:   res.TokenType,
		},
		ApiServer: res.ApiServer,
	}

	if secs := res.ExpiresIn; secs > 0 {
		token.Token.Expiry = time.Now().Add(time.Duration(secs) * time.Second)
	}

	btoken, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}

	err = keyring.Set("lunchquest-cli", "authenticate-response", string(btoken))
	if err != nil {
		return nil, err
	}

	return token, nil
}

// MakeQuestradeRequest makes an authenticated request to the Questrade API.
// Returns an error if no valid authentication credentials are found.
func makeQuestradeRequest(ctx context.Context, method, endpoint string, body io.Reader) (*http.Response, error) {
	// Try to get the token from keyring
	storedAuth, err := keyring.Get("lunchquest-cli", "authenticate-response")
	if err != nil {
		return nil, fmt.Errorf("no authentication credentials found: %w", err)
	}

	// Parse the stored token
	token := &customOauth2Token{}
	if err := json.Unmarshal([]byte(storedAuth), token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stored token: %w", err)
	}

	// Check if token is expired or will expire within 1 minute
	if time.Until(token.Expiry) < time.Minute {
		// Attempt to refresh the token
		newToken, err := Authenticate(ctx, viper.GetString("refreshToken"))
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}

		// Check if the refreshed token is still expired
		if newToken.Expiry.Before(time.Now()) {
			return nil, fmt.Errorf("refreshed token is still expired")
		}

		token = newToken
	}

	// Create client
	client := &http.Client{}

	// Build the full URL
	fullURL := token.ApiServer + endpoint

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))
	req.Header.Add("User-Agent", viper.GetString("UserAgent"))
	req.Header.Add("Content-Type", "application/json")

	// Make the request
	return client.Do(req)
}
