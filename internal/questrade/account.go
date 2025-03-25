package questrade

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type QuestradeAccount struct {
	Type              string `json:"type"`
	Number            string `json:"number"`
	Status            string `json:"status"`
	IsPrimary         bool   `json:"isPrimary"`
	IsBilling         bool   `json:"isBilling"`
	ClientAccountType string `json:"ClientAccountType"`
}

type GetAccountsResponse struct {
	Accounts []*QuestradeAccount `json:"accounts"`
	UserId   int64               `json:"userId"`
}

func GetAccounts(ctx context.Context) (*GetAccountsResponse, error) {
	resp, err := makeQuestradeRequest(ctx, http.MethodGet, "v1/accounts", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response body
	var accountsResponse GetAccountsResponse
	if err := json.NewDecoder(resp.Body).Decode(&accountsResponse); err != nil {
		return nil, fmt.Errorf("failed to parse accounts response: %w", err)
	}

	return &accountsResponse, nil
}
