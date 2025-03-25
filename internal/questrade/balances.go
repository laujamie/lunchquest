package questrade

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Balance represents a single currency balance for an account
type Balance struct {
	Currency          string  `json:"currency"`
	Cash              float64 `json:"cash"`
	MarketValue       float64 `json:"marketValue"`
	TotalEquity       float64 `json:"totalEquity"`
	BuyingPower       float64 `json:"buyingPower"`
	MaintenanceExcess float64 `json:"maintenanceExcess"`
	IsRealTime        bool    `json:"isRealTime"`
}

// GetBalancesResponse represents the response from the balances endpoint
type GetBalancesResponse struct {
	PerCurrencyBalances    []*Balance `json:"perCurrencyBalances"`
	CombinedBalances       []*Balance `json:"combinedBalances"`
	SodPerCurrencyBalances []*Balance `json:"sodPerCurrencyBalances"`
	SodCombinedBalanced    []*Balance `json:"sodCombinedBalances"`
}

// GetBalances retrieves the balances for a specific account
func GetBalances(ctx context.Context, accountNumber string) (*GetBalancesResponse, error) {
	endpoint := fmt.Sprintf("v1/accounts/%s/balances", accountNumber)

	resp, err := makeQuestradeRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get account balances: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var balancesResponse GetBalancesResponse
	if err := json.NewDecoder(resp.Body).Decode(&balancesResponse); err != nil {
		return nil, fmt.Errorf("failed to parse balances response: %w", err)
	}

	return &balancesResponse, nil
}
