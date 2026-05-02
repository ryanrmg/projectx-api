package projectx

import (
	"fmt"
)

type GatewayUserTrade struct {
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

func GatewayUserTradeCSVHeader() string {
	return "id,account_id,contract_id,creation_timestamp,price,profit_and_loss,fees,side,size,voided,order_id"
}

func (t GatewayUserTrade) String() string {
	return fmt.Sprintf(
		"GatewayUserTrade{Id:%d AccountId:%d ContractId:%s Price:%f PnL:%f Fees:%f Side:%d Size:%d Voided:%t}",
		t.Id,
		t.AccountId,
		t.ContractId,
		t.Price,
		t.ProfitAndLoss,
		t.Fees,
		t.Side,
		t.Size,
		t.Voided,
	)
}

func (t GatewayUserTrade) ToCSVRow() string {
	return fmt.Sprintf(
		"%d,%d,%s,%s,%.2f,%.2f,%.2f,%d,%d,%t,%d",
		t.Id,
		t.AccountId,
		t.ContractId,
		t.CreationTimestamp,
		t.Price,
		t.ProfitAndLoss,
		t.Fees,
		t.Side,
		t.Size,
		t.Voided,
		t.OrderId,
	)
}
