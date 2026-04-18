package projectx

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestAccountsSearch_LoginThenSearch(t *testing.T) {
	callCount := 0

	handler := func(w http.ResponseWriter, r *http.Request) {
		callCount++

		switch callCount {

		case 1:
			if r.URL.Path != "/Auth/loginKey" {
				t.Fatalf("expected login call, got %s", r.URL.Path)
			}

			resp := loginResponse{
				Token:     "test-token",
				ExpiresIn: 3600,
			}
			json.NewEncoder(w).Encode(resp)

		case 2:
			if r.URL.Path != "/Account/search" {
				t.Fatalf("expected account search, got %s", r.URL.Path)
			}

			if r.Header.Get("Authorization") != "Bearer test-token" {
				t.Fatalf("missing bearer token")
			}

			acct := []Account{
				{
					Id:        1,
					Name:      "Sim",
					Balance:   10000,
					CanTrade:  true,
					IsVisible: true,
				},
			}

			resp := AccountResponse{
				Accounts:     acct,
				Success:      true,
				ErrorCode:    0,
				ErrorMessage: "",
			}
			json.NewEncoder(w).Encode(resp)

		default:
			t.Fatalf("unexpected extra call")
		}
	}

	client, closeServer := newTestClient(handler)
	defer closeServer()

	accounts, err := client.Accounts.Search(
		context.Background(),
		AccountSearchRequest{OnlyActiveAccounts: true},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(accounts) != 1 {
		t.Fatalf("expected 1 account, got %d", len(accounts))
	}

	if accounts[0].Id != 1 {
		t.Fatalf("wrong account parsed")
	}
}

func TestAccountsSearch_UsesCachedToken(t *testing.T) {
	callCount := 0

	handler := func(w http.ResponseWriter, r *http.Request) {
		callCount++

		if r.URL.Path != "/Account/search" {
			t.Fatalf("expected only search call, got %s", r.URL.Path)
		}

		json.NewEncoder(w).Encode(AccountResponse{})
	}

	client, closeServer := newTestClient(handler)
	defer closeServer()

	// Pre-seed token (simulate prior login)
	client.token = "cached-token"
	client.expiresAt = time.Now().Add(time.Hour)

	_, err := client.Accounts.Search(
		context.Background(),
		AccountSearchRequest{OnlyActiveAccounts: true},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount != 1 {
		t.Fatalf("expected 1 request, got %d", callCount)
	}
}
