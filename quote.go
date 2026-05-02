package projectx

import (
	"fmt"
)

type GatewayQuote struct {
	Symbol        string  `json:"symbol"`
	SymbolName    string  `json:"symbolName,omitempty"`
	LastPrice     float64 `json:"lastPrice,omitempty"`
	BestBid       float64 `json:"bestBid,omitempty"`
	BestAsk       float64 `json:"bestAsk,omitempty"`
	Change        float64 `json:"change,omitempty"`
	ChangePercent float64 `json:"changePercent,omitempty"`
	Open          float64 `json:"open,omitempty"`
	High          float64 `json:"high,omitempty"`
	Low           float64 `json:"low,omitempty"`
	Volume        int64   `json:"volume,omitempty"`
	LastUpdated   string  `json:"lastUpdated"`
	Timestamp     string  `json:"timestamp"`
}

func GatewayQuoteCSVHeader() string {
	return "symbol,symbol_name,last_price,best_bid,best_ask,change,change_percent,open,high,low,volume,last_updated,timestamp"
}

func (q GatewayQuote) String() string {
	return fmt.Sprintf(
		"GatewayQuote{Symbol:%s Last:%f Bid:%f Ask:%f Volume:%d Timestamp:%s}",
		q.Symbol,
		q.LastPrice,
		q.BestBid,
		q.BestAsk,
		q.Volume,
		q.Timestamp,
	)
}

func (q GatewayQuote) ToCSVRow() string {
	return fmt.Sprintf(
		"%s,%s,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%d,%s,%s",
		q.Symbol,
		q.SymbolName,
		q.LastPrice,
		q.BestBid,
		q.BestAsk,
		q.Change,
		q.ChangePercent,
		q.Open,
		q.High,
		q.Low,
		q.Volume,
		q.LastUpdated,
		q.Timestamp,
	)
}
