package projectx

import (
	"context"
)

type BarHistoryRequest struct {
	ContractId        string `json:"contractId"`
	Live              bool   `json:"live"`
	StartTime         string `json:"startTime"`
	EndTime           string `json:"endTime"`
	Unit              int    `json:"unit"`       // second, minute, hour, day, week, month
	UnitNumber        int    `json:"unitNumber"` // number of units to aggregate
	Limit             int    `json:"limit"`
	IncludePartialBar bool   `json:"includePartialBar"`
}

type HistoryResponse struct {
	Bars         []Bar  `json:"bars"`
	Success      bool   `json:"success"`
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

type Bar struct {
	Timestamp string  `json:"t"`
	Open      float64 `json:"o"`
	High      float64 `json:"h"`
	Low       float64 `json:"l"`
	Close     float64 `json:"c"`
	Volume    float64 `json:"v"`
}

type AvailableContractRequest struct {
	Live bool `json:"live"`
}

type AvailableContractResponse struct {
	Contracts    []Contract `json:"contracts"`
	Success      bool       `json:"success"`
	ErrorCode    int        `json:"errorCode"`
	ErrorMessage string     `json:"errorMessage"`
}

type Contract struct {
	Id             string  `json:"id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	TickSize       float64 `json:"tickSize"`
	TickValue      float64 `json:"tickValue"`
	ActiveContract bool    `json:"activeContract"`
	SymbolId       string  `json:"symbolId"`
}

type MarketService struct {
	client *ProjectXClient
}

func (s *MarketService) History(ctx context.Context, req BarHistoryRequest) ([]Bar, error) {
	var resp HistoryResponse
	err := s.client.Post(
		ctx,
		"/History/retrieveBars",
		req,
		&resp,
	)
	if err != nil {
		return nil, err
	}

	return resp.Bars, nil
}

func (s *MarketService) AvailableContracts(ctx context.Context, req AvailableContractRequest) ([]Contract, error) {
	var resp AvailableContractResponse
	err := s.client.Post(
		ctx,
		"/Contract/available",
		req,
		&resp,
	)
	if err != nil {
		return nil, err
	}

	return resp.Contracts, nil
}
