package projectx

import (
	"net/http"
	"sync"
	"time"
)

type ProjectXClient struct {
	httpClient  *http.Client
	baseUrl     string
	marketWsUrl string
	userWsUrl   string
	username    string
	apiKey      string

	mu        sync.RWMutex
	token     string
	expiresAt time.Time

	Accounts *AccountService
	Markets  *MarketService
	Orders   *OrderService
	Trades   *TradeService

	Realtime *RealtimeService
}

func NewProjectXClient(baseUrl, wsUrl, username, apiKey string) *ProjectXClient {
	c := &ProjectXClient{
		baseUrl:     baseUrl,
		marketWsUrl: wsUrl + "/market",
		userWsUrl:   wsUrl + "/user",
		username:    username,
		apiKey:      apiKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
	c.Accounts = &AccountService{client: c}
	c.Markets = &MarketService{client: c}
	c.Orders = &OrderService{client: c}
	c.Trades = &TradeService{client: c}
	c.Realtime = NewRealtimeService(c)
	return c
}
