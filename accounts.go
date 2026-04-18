package projectx

import (
	"context"
)

type AccountSearchRequest struct {
	OnlyActiveAccounts bool `json:"onlyActiveAccounts"`
}

type AccountResponse struct {
	Accounts     []Account `json:"accounts"`
	Success      bool      `json:"success"`
	ErrorCode    int       `json:"errorCode"`
	ErrorMessage string    `json:"errorMessage"`
}

type Account struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	Balance   float64 `json:"balance"`
	CanTrade  bool    `json:"canTrade"`
	IsVisible bool    `json:"isVisible"`
}

type AccountService struct {
	client *ProjectXClient
}

func (s *AccountService) Search(ctx context.Context, req AccountSearchRequest) ([]Account, error) {
	var resp AccountResponse
	err := s.client.Post(
		ctx,
		"/Account/search",
		req,
		&resp,
	)
	if err != nil {
		return nil, err
	}

	return resp.Accounts, nil
}
