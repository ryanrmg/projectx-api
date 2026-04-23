package projectx

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net"
)

const (
	Bid   = 0
	Ask   = 1
	Https = "https://"
	Wss   = "wss://"
)

type RealtimeService struct {
	userConn   net.Conn
	marketConn net.Conn
	client     *ProjectXClient

	Trades chan json.RawMessage
	Quotes chan json.RawMessage
	Orders chan json.RawMessage
}

func NewRealtimeService(c *ProjectXClient) *RealtimeService {
	return &RealtimeService{
		client: c,
		Trades: make(chan json.RawMessage, 100),
		Quotes: make(chan json.RawMessage, 100),
		Orders: make(chan json.RawMessage, 100),
	}
}

type SubscribeMsg struct {
	Type      int      `json:"type"`
	Target    string   `json:"target"`
	Arguments []string `json:"arguments"`
}

type GatewayTrade struct {
	SymbolId  string  `json:"symbolId"`
	Price     float64 `json:"price"`
	Timestamp string  `json:"timestamp"`
	Type      int     `json:"type"`
	Volume    int64   `json:"volume"`
}

func (r *RealtimeService) Connect(ctx context.Context) error {
	marketNegotiateUrl := Https + r.client.marketWsUrl + "/negotiate?negotiateVersion=1"
	userNegotiateUrl := Https + r.client.userWsUrl + "/negotiate?negotiateVersion=1"
	marketWsUrl := Wss + r.client.marketWsUrl
	userWsUrl := Wss + r.client.userWsUrl

	token, err := r.client.getToken(ctx)
	if err != nil {
		return err
	}

	marketConn, err := GetWsConn(marketNegotiateUrl, marketWsUrl, token)
	if err != nil {
		log.Println("Failed to get market connection")
		return err
	}
	userConn, err := GetWsConn(userNegotiateUrl, userWsUrl, token)
	if err != nil {
		log.Println("Failed to get user connection")
		return err
	}

	r.marketConn = marketConn
	r.userConn = userConn

	// start background readers
	go r.readMarketLoop()
	// go r.readUserLoop()

	return nil
}

func (r *RealtimeService) SubscribeContractTrades(contractId string) error {
	if contractId != "" {
		log.Println("Subscribing to ContractTrades")
		msg := SubscribeMsg{
			Type:      1,
			Target:    "SubscribeContractTrades",
			Arguments: []string{contractId},
		}
		b, err := json.Marshal(msg)
		if err != nil {
			log.Println("failed to subscribe to ContractTrades")
			return err
		}
		b = append(b, RecordSep)
		WriteTextFrame(r.marketConn, b)
	}
	return nil
}

func (r *RealtimeService) readMarketLoop() {
	log.Println("read loop market")
	for {
		msg, err := ReadFullText(r.marketConn)
		if err != nil {
			log.Println("market read error:", err)
			return
		}

		r.handleFrame(msg)
	}
}

func (r *RealtimeService) readUserLoop() {
	log.Println("read loop user")
	for {
		msg, err := ReadFullText(r.userConn)
		log.Println(msg)
		if err != nil {
			log.Println("user read error:", err)
			return
		}

		r.handleFrame(msg)
	}
}

func splitSignalRMessages(frame []byte) [][]byte {
	parts := bytes.Split(frame, []byte{RecordSep})
	var msgs [][]byte
	for _, p := range parts {
		if len(p) > 0 {
			msgs = append(msgs, p)
		}
	}
	return msgs
}

func (r *RealtimeService) handleFrame(frame []byte) {
	messages := splitSignalRMessages(frame)
	// log.Println(string(messages))
	for _, msg := range messages {
		// log.Println(string(msg))

		var envelope struct {
			Type   int             `json:"type"`
			Target string          `json:"target"`
			Args   json.RawMessage `json:"arguments"`
		}

		if err := json.Unmarshal(msg, &envelope); err != nil {
			log.Println("bad json:", string(msg))
			continue
		}

		// type 6 = keepalive
		if envelope.Type != 1 {
			continue
		}

		log.Println(envelope.Target)

		switch envelope.Target {

		case "GatewayTrade":
			select {
			case r.Trades <- msg:
			default:
			}

		case "GatewayQuote":
			select {
			case r.Quotes <- msg:
			default:
			}

		case "ReceiveUserOrders":
			select {
			case r.Orders <- msg:
			default:
			}
		}
	}
}

func (r *RealtimeService) TradesStream() <-chan json.RawMessage {
	return r.Trades
}

func (r *RealtimeService) QuotesStream() <-chan json.RawMessage {
	return r.Quotes
}

func (r *RealtimeService) OrdersStream() <-chan json.RawMessage {
	return r.Orders
}
