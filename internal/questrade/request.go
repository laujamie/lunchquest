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
)

// makeQuestradeRequest makes an authenticated request to the Questrade API.
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

	fmt.Printf("%s\n", fullURL)

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
