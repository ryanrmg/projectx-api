package projectx

import (
	"fmt"
)

type GatewayUserAccount struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	Balance   float64 `json:"balance"`
	CanTrade  bool    `json:"canTrade"`
	IsVisible bool    `json:"isVisible"`
	Simulated bool    `json:"Simulated"`
}

func GatewayUserAccountCSVHeader() string {
	return "id,name,balance,can_trade,is_visible,simulated"
}

func (a GatewayUserAccount) String() string {
	return fmt.Sprintf(
		"GatewayUserAccount{Id:%d Name:%s Balance:%f CanTrade:%t IsVisible:%t Simulated:%t}",
		a.Id,
		a.Name,
		a.Balance,
		a.CanTrade,
		a.IsVisible,
		a.Simulated,
	)
}

func (a GatewayUserAccount) ToCSVRow() string {
	return fmt.Sprintf(
		"%d,%s,%.2f,%t,%t,%t",
		a.Id,
		a.Name,
		a.Balance,
		a.CanTrade,
		a.IsVisible,
		a.Simulated,
	)
}
