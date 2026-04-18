package projectx

import (
	"net/http"
	"sync"
	"time"
)

type ProjectXClient struct {
	httpClient *http.Client
	baseUrl    string
	username   string
	apiKey     string

	mu        sync.RWMutex
	token     string
	expiresAt time.Time

	Accounts *AccountService
	// Markets  *MarketService
}

func NewProjectXClient(baseUrl, username, apiKey string) *ProjectXClient {
	c := &ProjectXClient{
		baseUrl:  baseUrl,
		username: username,
		apiKey:   apiKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
	c.Accounts = &AccountService{client: c}
	// c.Markets = &MarketService{client: c}
	return c
}
