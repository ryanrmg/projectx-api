package projectx

import (
	"fmt"
)

type GatewayUserOrder struct {
	Id                int     `json:"id"`
	AccountId         int     `json:"accountId"`
	ContractId        string  `json:"contractId"`
	SymbolId          string  `json:"symbolId"`
	CreationTimestamp string  `json:"creationTimestamp"`
	UpdateTimestamp   string  `json:"updateTimestamp"`
	Status            int     `json:"status"`
	Type              int     `json:"type"`
	Side              int     `json:"side"`
	Size              int     `json:"size"`
	LimitPrice        float64 `json:"limitPrice"`
	StopPrice         float64 `json:"stopPrice"`
	FillVolume        int     `json:"fillVolume"`
	FilledPrice       float64 `json:"filledPrice"`
	CustomTag         string  `json:"customTag"`
}

func GatewayUserOrderCSVHeader() string {
	return "id,account_id,contract_id,symbol_id,creation_timestamp,update_timestamp,status,type,side,size,limit_price,stop_price,fill_volume,filled_price,custom_tag"
}

func (o GatewayUserOrder) String() string {
	return fmt.Sprintf(
		"GatewayUserOrder{Id:%d AccountId:%d ContractId:%s SymbolId:%s Status:%d Type:%d Side:%d Size:%d FillVolume:%d FilledPrice:%f}",
		o.Id,
		o.AccountId,
		o.ContractId,
		o.SymbolId,
		o.Status,
		o.Type,
		o.Side,
		o.Size,
		o.FillVolume,
		o.FilledPrice,
	)
}

func (o GatewayUserOrder) ToCSVRow() string {
	return fmt.Sprintf(
		"%d,%d,%s,%s,%s,%s,%d,%d,%d,%d,%.2f,%.2f,%d,%.2f,%s",
		o.Id,
		o.AccountId,
		o.ContractId,
		o.SymbolId,
		o.CreationTimestamp,
		o.UpdateTimestamp,
		o.Status,
		o.Type,
		o.Side,
		o.Size,
		o.LimitPrice,
		o.StopPrice,
		o.FillVolume,
		o.FilledPrice,
		o.CustomTag,
	)
}
