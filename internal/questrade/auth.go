package questrade

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/laujamie/lunchquest/internal/constants"
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
		if err := json.Unmarshal(b, &storedAuthToken); err == nil && storedAuthToken.RefreshToken == refreshToken && storedAuthToken.Expiry.After(time.Now().Add(time.Second*180)) {
			return nil, nil
		}
	}

	// Get auth token
	qc := oauth2.NewClient(ctx, nil)

	req, err := http.NewRequest(http.MethodGet, questradeLoginEndpoint, nil)
	if err != nil {
		log.Printf("request failed: %v", err)
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res authenticateResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	token := &customOauth2Token{
		Token: oauth2.Token{
			AccessToken:  res.AccessToken,
			TokenType:    res.TokenType,
			RefreshToken: res.RefreshToken,
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

	err = keyring.Set(constants.SERVICE_NAME, "authenticate-response", string(btoken))
	if err != nil {
		return nil, err
	}

	return token, nil
}

// GetStoredAuthToken retrieves the stored authentication token from the keyring
func GetStoredAuthToken() (*customOauth2Token, error) {
	storedAuth, err := keyring.Get(constants.SERVICE_NAME, "authenticate-response")
	if err != nil {
		return nil, err
	}

	var token customOauth2Token
	if err := json.Unmarshal([]byte(storedAuth), &token); err != nil {
		return nil, err
	}

	return &token, nil
}
