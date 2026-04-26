package projectx

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestTradeSearch_LoginThenSearch(t *testing.T) {
	callCount := 0

	handler := func(w http.ResponseWriter, r *http.Request) {
		callCount++

		if r.URL.Path != "/Trade/search" {
			t.Fatalf("expected only search call, got %s", r.URL.Path)
		}

		json.NewEncoder(w).Encode(AccountResponse{})
	}

	client, closeServer := newTestClient(handler)
	defer closeServer()

	// Pre-seed token (simulate prior login)
	client.token = "cached-token"
	client.expiresAt = time.Now().Add(time.Hour)

	_, err := client.Trades.Search(
		context.Background(),
		TradeSearchRequest{
			AccountId:      123456789,
			EndTimestamp:   time.Now().UTC().Format(time.RFC3339),
			StartTimestamp: time.Now().UTC().Add(-time.Duration(24*5) * time.Hour).Format(time.RFC3339),
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount != 1 {
		t.Fatalf("expected 1 request, got %d", callCount)
	}
}
