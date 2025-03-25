package questrade

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ActivityType represents the type of account activity
type ActivityType string

const (
	ActivityDeposit         ActivityType = "Deposit"
	ActivityDividend        ActivityType = "Dividend"
	ActivityFee             ActivityType = "Fee"
	ActivityInterest        ActivityType = "Interest"
	ActivityRebate          ActivityType = "Rebate"
	ActivityTrade           ActivityType = "Trade"
	ActivityWithdrawal      ActivityType = "Withdrawal"
	ActivityTransfer        ActivityType = "Transfer"
	ActivityCorporateAction ActivityType = "CorporateAction"
	ActivityOther           ActivityType = "Other"
)

// Activity represents a single account activity
type Activity struct {
	TradeDate       time.Time    `json:"tradeDate"`
	TransactionDate time.Time    `json:"transactionDate"`
	SettlementDate  time.Time    `json:"settlementDate"`
	Action          string       `json:"action"`
	Symbol          string       `json:"symbol"`
	SymbolID        int64        `json:"symbolId"`
	Description     string       `json:"description"`
	Currency        string       `json:"currency"`
	Quantity        float64      `json:"quantity"`
	Price           float64      `json:"price"`
	GrossAmount     float64      `json:"grossAmount"`
	Commission      float64      `json:"commission"`
	NetAmount       float64      `json:"netAmount"`
	Type            ActivityType `json:"type"`
}

// GetActivitiesResponse represents the response from the activities endpoint
type GetActivitiesResponse struct {
	Activities []*Activity `json:"activities"`
}

// GetActivities retrieves the activities for a specific account within a date range
func GetActivities(
	ctx context.Context,
	accountNumber string,
	startTime, endTime time.Time,
) (*GetActivitiesResponse, error) {
	// Format endpoint with query parameters
	endpoint := fmt.Sprintf(
		"v1/accounts/%s/activities?startTime=%s&endTime=%s",
		accountNumber,
		startTime.Format("2006-01-02T15:04:05-07:00"),
		endTime.Format("2006-01-02T15:04:05-07:00"),
	)

	resp, err := makeQuestradeRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get account activities: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var activitiesResponse GetActivitiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&activitiesResponse); err != nil {
		return nil, fmt.Errorf("failed to parse activities response: %w", err)
	}

	return &activitiesResponse, nil
}
