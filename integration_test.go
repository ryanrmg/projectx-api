package projectx

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestIntegration_Orders(t *testing.T) {
	username := os.Getenv("PROJECTX_USERNAME")
	apiKey := os.Getenv("PROJECTX_API_KEY")

	// Skip if creds not provided
	if username == "" || apiKey == "" {
		t.Skip("Skipping integration test (no credentials provided)")
	}

	client := NewProjectXClient(
		"https://api.topstepx.com/api",
		username,
		apiKey,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	accounts, err := client.Accounts.Search(
		ctx,
		AccountSearchRequest{OnlyActiveAccounts: true},
	)

	if err != nil {
		t.Logf("did not get any accounts to search orders on")
	}

	if len(accounts) > 0 { // search orders on account
		orders, err := client.Orders.OrderSearch(
			ctx,
			OrderSearchRequest{
				AccountId:      accounts[0].Id,
				StartTimestamp: time.Now().UTC().Add(-time.Duration(24*3) * time.Hour).Format(time.RFC3339),
				EndTimestamp:   time.Now().UTC().Format(time.RFC3339),
			},
		)

		if err != nil {
			t.Logf("got no orders, check if you expect orders on this acct")
		}
		t.Logf("%v", orders)
	}

}

func TestIntegration_Market(t *testing.T) {
	username := os.Getenv("PROJECTX_USERNAME")
	apiKey := os.Getenv("PROJECTX_API_KEY")

	// Skip if creds not provided
	if username == "" || apiKey == "" {
		t.Skip("Skipping integration test (no credentials provided)")
	}

	client := NewProjectXClient(
		"https://api.topstepx.com/api",
		username,
		apiKey,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	bars, err := client.Markets.History(
		ctx,
		BarHistoryRequest{
			ContractId:        "CON.F.US.MNQ.M26",
			Live:              false,
			EndTime:           time.Now().UTC().Format(time.RFC3339),
			StartTime:         time.Now().UTC().Add(-time.Duration(24*5) * time.Hour).Format(time.RFC3339),
			Unit:              2,
			UnitNumber:        5,
			Limit:             100,
			IncludePartialBar: false,
		},
	)

	if err != nil {
		t.Fatalf("real API call failed: %v", err)
	}

	if len(bars) == 0 {
		t.Fatalf("expected at least one bar")
	}

	t.Logf("Got history for %d bars", len(bars))
}

func TestIntegration_AccountsSearch(t *testing.T) {
	username := os.Getenv("PROJECTX_USERNAME")
	apiKey := os.Getenv("PROJECTX_API_KEY")

	// Skip if creds not provided
	if username == "" || apiKey == "" {
		t.Skip("Skipping integration test (no credentials provided)")
	}

	client := NewProjectXClient(
		"https://api.topstepx.com/api",
		username,
		apiKey,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	accounts, err := client.Accounts.Search(
		ctx,
		AccountSearchRequest{OnlyActiveAccounts: true},
	)
	if err != nil {
		t.Fatalf("real API call failed: %v", err)
	}

	if len(accounts) == 0 {
		t.Fatalf("expected at least one account")
	}

	t.Logf("Found %d accounts", len(accounts))
}
