package projectx

import (
	"fmt"
)

type GatewayTrade struct {
	SymbolId  string  `json:"symbolId"`
	Price     float64 `json:"price"`
	Timestamp string  `json:"timestamp"`
	Type      int     `json:"type"`
	Volume    int64   `json:"volume"`
}

func GatewayTradeCSVHeader() string {
	return "symbol_id,price,timestamp,type,volume"
}

func (t GatewayTrade) String() string {
	return fmt.Sprintf(
		"GatewayTrade{SymbolId:%s Price:%f Timestamp:%s Type:%d Volume:%d}",
		t.SymbolId,
		t.Price,
		t.Timestamp,
		t.Type,
		t.Volume,
	)
}

func (t GatewayTrade) ToCSVRow() string {
	return fmt.Sprintf(
		"%s,%.2f,%s,%d,%d",
		t.SymbolId,
		t.Price,
		t.Timestamp,
		t.Type,
		t.Volume,
	)
}
