package projectx

import (
	"fmt"
)

type GatewayUserPosition struct {
	Id                int     `json:"id"`
	AccountId         int     `json:"accountId"`
	ContractId        string  `json:"contractId"`
	CreationTimestamp string  `json:"creationTimestamp"`
	Type              int     `json:"type"`
	Size              int     `json:"size"`
	AveragePrice      float64 `json:"averagePrice"`
}

func GatewayUserPositionCSVHeader() string {
	return "id,account_id,contract_id,creation_timestamp,type,size,average_price"
}

func (p GatewayUserPosition) String() string {
	return fmt.Sprintf(
		"GatewayUserPosition{Id:%d AccountId:%d ContractId:%s Type:%d Size:%d AvgPrice:%f}",
		p.Id,
		p.AccountId,
		p.ContractId,
		p.Type,
		p.Size,
		p.AveragePrice,
	)
}

func (p GatewayUserPosition) ToCSVRow() string {
	return fmt.Sprintf(
		"%d,%d,%s,%s,%d,%d,%.2f",
		p.Id,
		p.AccountId,
		p.ContractId,
		p.CreationTimestamp,
		p.Type,
		p.Size,
		p.AveragePrice,
	)
}
