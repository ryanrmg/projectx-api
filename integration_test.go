package projectx

import (
	"context"
	"os"
	"testing"
	"time"
)

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
