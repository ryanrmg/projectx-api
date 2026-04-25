package projectx

import (
	"context"
)

type BasicPositionResponse struct {
	Success      bool   `json:"success"`
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

type PositionService struct {
	client *ProjectXClient
}

type CloseContractRequest struct {
	AccountId  int    `json:"accountId"`
	ContractId string `json:"contractId"`
}

type PartialCloseContractRequest struct {
	AccountId  int    `json:"accountId"`
	ContractId string `json:"contractId"`
	Size       int    `json:"size"`
}

type SearchOpenContractRequest struct {
	AccountId int `json:"accountId"`
}

func (p *PositionService) CloseContract(ctx context.Context, req CloseContractRequest) (bool, error) {
	var resp BasicPositionResponse
	err := p.client.Post(
		ctx,
		"/Position/closeContract",
		req,
		&resp,
	)
	if err != nil {
		return false, err
	}

	return resp.Success, nil
}

func (p *PositionService) PartialCloseContract(ctx context.Context, req PartialCloseContractRequest) (bool, error) {
	var resp BasicPositionResponse
	err := p.client.Post(
		ctx,
		"/Position/partialCloseContract",
		req,
		&resp,
	)
	if err != nil {
		return false, err
	}

	return resp.Success, nil
}

func (p *PositionService) SearchOpenContract(ctx context.Context, req SearchOpenContractRequest) (bool, error) {
	var resp BasicPositionResponse
	err := p.client.Post(
		ctx,
		"/Position/searchOpen",
		req,
		&resp,
	)
	if err != nil {
		return false, err
	}

	return resp.Success, nil
}
