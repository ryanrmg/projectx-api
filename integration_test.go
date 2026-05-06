package projectx

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

var (
	username = os.Getenv("PROJECTX_USERNAME")
	apiKey   = os.Getenv("PROJECTX_API_KEY")
	client = NewProjectXClient(
		"https://api.topstepx.com/api",
		"rtc.topstepx.com/hubs/",
		os.Getenv("PROJECTX_USERNAME"),
		os.Getenv("PROJECTX_API_KEY"),
	)
)

func TestIntegration_Orders(t *testing.T) {

	// Skip if creds not provided
	if username == "" || apiKey == "" {
		t.Skip("Skipping integration test (no credentials provided)")
	}

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

	// Skip if creds not provided
	if username == "" || apiKey == "" {
		t.Skip("Skipping integration test (no credentials provided)")
	}

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

	// Skip if creds not provided
	if username == "" || apiKey == "" {
		t.Skip("Skipping integration test (no credentials provided)")
	}

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

func TestRealtime_Trades(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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

	timeout := time.After(5 * time.Second)
	received := false

	for !received {
		select {

		case msg := <-trades:
			fmt.Println("TRADE:", msg.String())

			t.Log("SUCCESS — received realtime trade!")
			received = true

		// TIMEOUT
		case <-timeout:
			t.Skip("timeout waiting for realtime trade message, (market likely closed)")
			return
		}
	}
}

func TestRealtime_Quotes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// connect
	if err := client.Realtime.Connect(ctx); err != nil {
		t.Fatalf("connect error: %v", err)
	}

	// subscribe
	contract := "CON.F.US.MNQ.M26"
	if err := client.Realtime.SubscribeContractQuotes(contract); err != nil {
		t.Fatalf("subscribe error: %v", err)
	}

	t.Log("subscribed — waiting for trade message...")

	quotes := client.Realtime.QuotesStream()

	timeout := time.After(5 * time.Second)
	received := false

	for !received {
		select {

		case msg := <-quotes:
			fmt.Println("QUOTE:", msg.String())

			t.Log("SUCCESS — received realtime trade!")
			received = true

		// TIMEOUT
		case <-timeout:
			t.Skip("timeout waiting for realtime trade message, (market likely closed)")
			return
		}
	}
}

func TestRealtime_Depth(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// connect
	if err := client.Realtime.Connect(ctx); err != nil {
		t.Fatalf("connect error: %v", err)
	}

	// subscribe
	contract := "CON.F.US.MNQ.M26"
	if err := client.Realtime.SubscribeContractMarketDepth(contract); err != nil {
		t.Fatalf("subscribe error: %v", err)
	}

	t.Log("subscribed — waiting for trade message...")

	depth := client.Realtime.DepthStream()

	timeout := time.After(5 * time.Second)
	received := false

	for !received {
		select {

		case msg := <-depth:
			fmt.Println("DEPTH:", msg.String())

			t.Log("SUCCESS — received realtime trade!")
			received = true

		// TIMEOUT
		case <-timeout:
			t.Skip("timeout waiting for realtime trade message, (market likely closed)")
			return
		}
	}
}