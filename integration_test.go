package projectx

import (
	"context"
	"encoding/json"
	"fmt"
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
		"https://rtc.topstepx.com/hubs/",
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
		"https://rtc.topstepx.com/hubs/",
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
		"https://rtc.topstepx.com/hubs/",
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

func TestRealtime_LiveIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := NewProjectXClient(
		"https://api.topstepx.com/api",
		"rtc.topstepx.com/hubs",
		os.Getenv("PROJECTX_USERNAME"),
		os.Getenv("PROJECTX_API_KEY"),
	)

	// connect
	if err := client.Realtime.Connect(ctx); err != nil {
		t.Fatalf("connect error: %v", err)
	}

	// subscribe
	contract := "CON.F.US.MNQ.M26"
	if err := client.Realtime.SubscribeContractTrades(contract); err != nil {
		t.Fatalf("subscribe error: %v", err)
	}

	t.Log("subscribed — waiting for trade message...")

	trades := client.Realtime.TradesStream()

	timeout := time.After(10 * time.Second)
	received := false

	for !received {
		select {

		case msg := <-trades:
			fmt.Println("TRADE RAW:", string(msg))

			// Optional: validate structure
			var envelope struct {
				Type      int             `json:"type"`
				Target    string          `json:"target"`
				Arguments json.RawMessage `json:"arguments"`
			}

			if err := json.Unmarshal(msg, &envelope); err != nil {
				t.Fatalf("received invalid JSON: %v", err)
			}

			if envelope.Target != "GatewayTrade" {
				t.Fatalf("unexpected target: %s", envelope.Target)
			}

			t.Log("SUCCESS — received realtime trade!")
			received = true

		// TIMEOUT
		case <-timeout:
			t.Fatal("timeout waiting for realtime trade message")
		}
	}
}
