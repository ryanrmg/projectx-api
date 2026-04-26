package projectx

import (
	"context"
)

type TradeService struct {
	client *ProjectXClient
}

type TradeSearchRequest struct {
	AccountId      int    `json:"accountId"`
	StartTimestamp string `json:"startTimestamp"`
	EndTimestamp   string `json:"endTimestamp"`
}

type Trade struct {
	Id                int     `json:"id"`
	AccountId         int     `json:"accountId"`
	ContractId        string  `json:"contractId"`
	CreationTimestamp string  `json:"creationTimestamp"`
	Price             float64 `json:"price"`
	ProfitAndLoss     float64 `json:"profitAndLoss"`
	Fees              float64 `json:"fees"`
	Side              int     `json:"side"`
	Size              int     `json:"size"`
	Voided            bool    `json:"voided"`
	OrderId           int     `json:"orderId"`
}

type TradeSearchResponse struct {
	Trades       []Trade `json:"trades"`
	Success      bool    `json:"success"`
	ErrorCode    int     `json:"errorCode"`
	ErrorMessage string  `json:"errorMessage"`
}

func (s *TradeService) Search(ctx context.Context, req TradeSearchRequest) ([]Trade, error) {
	var resp TradeSearchResponse
	err := s.client.Post(
		ctx,
		"/Trade/search",
		req,
		&resp,
	)
	if err != nil {
		return nil, err
	}

	return resp.Trades, nil
}
