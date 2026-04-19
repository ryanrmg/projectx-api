package projectx

import (
	"context"
	"time"
)

type OrderSearchRequest struct {
	AccountId      int    `json:"accountId"`
	StartTimestamp string `json:"startTimestamp"`
	EndTimestamp   string `json:"endTimestamp"`
}

type OpenOrderSearchRequest struct {
	AccountId int `json:"accountId"`
}

type OrderSearchResponse struct {
	Orders       []Order `json:"orders"`
	Success      bool    `json:"success"`
	ErrorCode    int     `json:"errorCode"`
	ErrorMessage string  `json:"errorMessage"`
}

type OpenOrderSearchResponse struct {
	Orders       []OpenOrder `json:"orders"`
	Success      bool        `json:"success"`
	ErrorCode    int         `json:"errorCode"`
	ErrorMessage string      `json:"errorMessage"`
}

type Order struct {
	ID                int64     `json:"id"`
	AccountID         int64     `json:"accountId"`
	ContractID        string    `json:"contractId"`
	SymbolID          string    `json:"symbolId"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
	UpdateTimestamp   time.Time `json:"updateTimestamp"`
	Status            int       `json:"status"`
	Type              int       `json:"type"`
	Side              int       `json:"side"`
	Size              int       `json:"size"`

	LimitPrice  *float64 `json:"limitPrice"`
	StopPrice   *float64 `json:"stopPrice"`
	FillVolume  int      `json:"fillVolume"`
	FilledPrice *float64 `json:"filledPrice"`

	CustomTag *string `json:"customTag"`
}

type OpenOrder struct {
	ID                int64     `json:"id"`
	AccountID         int64     `json:"accountId"`
	ContractID        string    `json:"contractId"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
	UpdateTimestamp   time.Time `json:"updateTimestamp"`

	Status int `json:"status"`
	Type   int `json:"type"`
	Side   int `json:"side"`
	Size   int `json:"size"`

	LimitPrice  *float64 `json:"limitPrice"`
	StopPrice   *float64 `json:"stopPrice"`
	FilledPrice *float64 `json:"filledPrice"`
}

type OrderService struct {
	client *ProjectXClient
}

func (s *OrderService) OrderSearch(ctx context.Context, req OrderSearchRequest) ([]Order, error) {
	var resp OrderSearchResponse
	err := s.client.Post(
		ctx,
		"/Order/search",
		req,
		&resp,
	)
	if err != nil {
		return nil, err
	}

	return resp.Orders, nil
}

func (s *OrderService) OpenOrderSearch(ctx context.Context, req OpenOrderSearchRequest) ([]OpenOrder, error) {
	var resp OpenOrderSearchResponse
	err := s.client.Post(
		ctx,
		"/Order/searchOpen",
		req,
		&resp,
	)
	if err != nil {
		return nil, err
	}

	return resp.Orders, nil
}
