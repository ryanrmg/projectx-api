package projectx

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestMarketHistory_LoginThenSearch(t *testing.T) {
	callCount := 0

	handler := func(w http.ResponseWriter, r *http.Request) {
		callCount++

		if r.URL.Path != "/History/retrieveBars" {
			t.Fatalf("expected only search call, got %s", r.URL.Path)
		}

		json.NewEncoder(w).Encode(AccountResponse{})
	}

	client, closeServer := newTestClient(handler)
	defer closeServer()

	// Pre-seed token (simulate prior login)
	client.token = "cached-token"
	client.expiresAt = time.Now().Add(time.Hour)

	_, err := client.Markets.History(
		context.Background(),
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
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount != 1 {
		t.Fatalf("expected 1 request, got %d", callCount)
	}
}

func TestAvailableContracts_Live(t *testing.T) {
	callCount := 0

	handler := func(w http.ResponseWriter, r *http.Request) {
		callCount++

		if r.URL.Path != "/Contract/available" {
			t.Fatalf("expected only search call, got %s", r.URL.Path)
		}

		json.NewEncoder(w).Encode(AccountResponse{})
	}

	client, closeServer := newTestClient(handler)
	defer closeServer()

	// Pre-seed token (simulate prior login)
	client.token = "cached-token"
	client.expiresAt = time.Now().Add(time.Hour)

	_, err := client.Markets.AvailableContracts(
		context.Background(),
		AvailableContractRequest{Live: true},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount != 1 {
		t.Fatalf("expected 1 request, got %d", callCount)
	}
}
