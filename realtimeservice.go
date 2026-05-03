package projectx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"

	"sync"
)

const (
	Bid   = 0
	Ask   = 1
	Https = "https://"
	Wss   = "wss://"

	CHAN_SIZE = 100
)

type RealtimeService struct {
	userConn   net.Conn
	marketConn net.Conn
	client     *ProjectXClient

	mu sync.RWMutex
	// Trade chan GatewayTrade
	// Quote chan GatewayQuote
	// Depth chan GatewayDepth

	tradeSubs map[chan GatewayTrade]struct{}
	quoteSubs map[chan GatewayQuote]struct{}
	depthSubs map[chan GatewayDepth]struct{}

	// UserOrder    chan GatewayUserOrder
	// UserTrade    chan GatewayUserTrade
	// UserPosition chan GatewayUserPosition
	// UserAccount  chan GatewayUserAccount

	userOrderSubs    map[chan GatewayUserOrder]struct{}
	userTradeSubs    map[chan GatewayUserTrade]struct{}
	userPositionSubs map[chan GatewayUserPosition]struct{}
	userAccountSubs  map[chan GatewayUserAccount]struct{}
}

func NewRealtimeService(c *ProjectXClient) *RealtimeService {
	return &RealtimeService{
		client:    c,
		tradeSubs: make(map[chan GatewayTrade]struct{}),
		quoteSubs: make(map[chan GatewayQuote]struct{}),
		depthSubs: make(map[chan GatewayDepth]struct{}),

		userOrderSubs:    make(map[chan GatewayUserOrder]struct{}),
		userTradeSubs:    make(map[chan GatewayUserTrade]struct{}),
		userPositionSubs: make(map[chan GatewayUserPosition]struct{}),
		userAccountSubs:  make(map[chan GatewayUserAccount]struct{}),
		// Trade:  make(chan GatewayTrade, 100),
		// Quote:  make(chan GatewayQuote, 100),
		// Depth:  make(chan GatewayDepth, 100),

		// UserOrder:    make(chan GatewayUserOrder, 100),
		// UserTrade:    make(chan GatewayUserTrade, 100),
		// UserPosition: make(chan GatewayUserPosition, 100),
		// UserAccount:  make(chan GatewayUserAccount, 100),
	}
}

type SubscribeMsg struct {
	Type      int      `json:"type"`
	Target    string   `json:"target"`
	Arguments []string `json:"arguments"`
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
	go r.readUserLoop()

	return nil
}

func (r *RealtimeService) SubscribeContractTrades(contractId string) error {
	if contractId != "" {
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
		return nil
	}
	return errors.New("contract id is nil")
}

func (r *RealtimeService) SubscribeContractQuotes(contractId string) error {
	if contractId != "" {
		msg := SubscribeMsg{
			Type:      1,
			Target:    "SubscribeContractQuotes",
			Arguments: []string{contractId},
		}
		b, err := json.Marshal(msg)
		if err != nil {
			log.Println("failed to subscribe to ContractQuotes")
			return err
		}
		b = append(b, RecordSep)
		WriteTextFrame(r.marketConn, b)
		return nil
	}
	return errors.New("contract id is nil")
}

func (r *RealtimeService) SubscribeContractMarketDepth(contractId string) error {
	if contractId != "" {
		msg := SubscribeMsg{
			Type:      1,
			Target:    "SubscribeContractMarketDepth",
			Arguments: []string{contractId},
		}
		b, err := json.Marshal(msg)
		if err != nil {
			log.Println("failed to subscribe to ContractMarketDepth")
			return err
		}
		b = append(b, RecordSep)
		WriteTextFrame(r.marketConn, b)
		return nil
	}
	return errors.New("contract id is nil")
}

func (r *RealtimeService) SubscribeAccounts(contractId string) error {
	if contractId != "" {
		msg := SubscribeMsg{
			Type:      1,
			Target:    "SubscribeAccounts",
			Arguments: []string{contractId},
		}
		b, err := json.Marshal(msg)
		if err != nil {
			log.Println("failed to subscribe to Accounts")
			return err
		}
		b = append(b, RecordSep)
		WriteTextFrame(r.userConn, b)
		return nil
	}
	return errors.New("contract id is nil")
}

func (r *RealtimeService) SubscribeOrders(contractId string) error {
	if contractId != "" {
		msg := SubscribeMsg{
			Type:      1,
			Target:    "SubscribeOrders",
			Arguments: []string{contractId},
		}
		b, err := json.Marshal(msg)
		if err != nil {
			log.Println("failed to subscribe to Orders")
			return err
		}
		b = append(b, RecordSep)
		WriteTextFrame(r.userConn, b)
		return nil
	}
	return errors.New("contract id is nil")
}

func (r *RealtimeService) SubscribePositions(contractId string) error {
	if contractId != "" {
		msg := SubscribeMsg{
			Type:      1,
			Target:    "SubscribePositions",
			Arguments: []string{contractId},
		}
		b, err := json.Marshal(msg)
		if err != nil {
			log.Println("failed to subscribe to Positions")
			return err
		}
		b = append(b, RecordSep)
		WriteTextFrame(r.userConn, b)
		return nil
	}
	return errors.New("contract id is nil")
}

func (r *RealtimeService) SubscribeTrades(contractId string) error {
	if contractId != "" {
		msg := SubscribeMsg{
			Type:      1,
			Target:    "SubscribeTrades",
			Arguments: []string{contractId},
		}
		b, err := json.Marshal(msg)
		if err != nil {
			log.Println("failed to subscribe to Trades")
			return err
		}
		b = append(b, RecordSep)
		WriteTextFrame(r.userConn, b)
		return nil
	}
	return errors.New("contract id is nil")
}

func (r *RealtimeService) readMarketLoop() {
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
	for {
		msg, err := ReadFullText(r.userConn)
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
	for _, msg := range messages {

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

		switch envelope.Target {

		case "GatewayTrade":
			var msg []GatewayTrade
			if err := json.Unmarshal(envelope.Args, &msg); err == nil {
				for _, v := range msg {
					r.broadcastTrade(v)
				}
			}

		case "GatewayQuote":
			var msg []GatewayQuote
			if err := json.Unmarshal(envelope.Args, &msg); err == nil {
				for _, v := range msg {
					r.broadcastQuote(v)
				}
			}

		case "GatewayDepth":
			var msg []GatewayDepth
			if err := json.Unmarshal(envelope.Args, &msg); err == nil {
				for _, v := range msg {
					r.broadcastDepth(v)
				}
			}

		case "GatewayUserAccount":
			var msg []GatewayUserAccount
			if err := json.Unmarshal(envelope.Args, &msg); err == nil {
				for _, v := range msg {
					r.broadcastUserAccount(v)
				}
			}

		case "GatewayUserPosition":
			var msg []GatewayUserPosition
			if err := json.Unmarshal(envelope.Args, &msg); err == nil {
				for _, v := range msg {
					r.broadcastUserPosition(v)
				}
			}

		case "GatewayUserOrder":
			var msg []GatewayUserOrder
			if err := json.Unmarshal(envelope.Args, &msg); err == nil {
				for _, v := range msg {
					r.broadcastUserOrder(v)
				}
			}

		case "GatewayUserTrade":
			var msg []GatewayUserTrade
			if err := json.Unmarshal(envelope.Args, &msg); err == nil {
				for _, v := range msg {
					r.broadcastUserTrade(v)
				}
			}
		}
	}
}

func (r *RealtimeService) broadcastTrade(t GatewayTrade) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for sub := range r.tradeSubs {
		select {
		case sub <- t:
		default:
		}
	}
}

func (r *RealtimeService) broadcastQuote(q GatewayQuote) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for sub := range r.quoteSubs {
		select {
		case sub <- q:
		default:
		}
	}
}

func (r *RealtimeService) broadcastDepth(d GatewayDepth) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for sub := range r.depthSubs {
		select {
		case sub <- d:
		default:
		}
	}
}

func (r *RealtimeService) broadcastUserAccount(q GatewayUserAccount) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for sub := range r.userAccountSubs {
		select {
		case sub <- q:
		default:
		}
	}
}

func (r *RealtimeService) broadcastUserPosition(d GatewayUserPosition) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for sub := range r.userPositionSubs {
		select {
		case sub <- d:
		default:
		}
	}
}

func (r *RealtimeService) broadcastUserOrder(q GatewayUserOrder) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for sub := range r.userOrderSubs {
		select {
		case sub <- q:
		default:
		}
	}
}

func (r *RealtimeService) broadcastUserTrade(d GatewayUserTrade) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for sub := range r.userTradeSubs {
		select {
		case sub <- d:
		default:
		}
	}
}

func (r *RealtimeService) TradesStream() <-chan GatewayTrade {
	ch := make(chan GatewayTrade, CHAN_SIZE)

	r.mu.Lock()
	r.tradeSubs[ch] = struct{}{}
	r.mu.Unlock()

	return ch
}

func (r *RealtimeService) QuotesStream() <-chan GatewayQuote {
	ch := make(chan GatewayQuote, CHAN_SIZE)

	r.mu.Lock()
	r.quoteSubs[ch] = struct{}{}
	r.mu.Unlock()

	return ch
}

func (r *RealtimeService) DepthStream() <-chan GatewayDepth {
	ch := make(chan GatewayDepth, CHAN_SIZE)

	r.mu.Lock()
	r.depthSubs[ch] = struct{}{}
	r.mu.Unlock()

	return ch
}

func (r *RealtimeService) UserAccountStream() <-chan GatewayUserAccount {
	ch := make(chan GatewayUserAccount, CHAN_SIZE)

	r.mu.Lock()
	r.userAccountSubs[ch] = struct{}{}
	r.mu.Unlock()

	return ch
}

func (r *RealtimeService) UserPositionStream() <-chan GatewayUserPosition {
	ch := make(chan GatewayUserPosition, CHAN_SIZE)

	r.mu.Lock()
	r.userPositionSubs[ch] = struct{}{}
	r.mu.Unlock()

	return ch
}

func (r *RealtimeService) UserOrdersStream() <-chan GatewayUserOrder {
	ch := make(chan GatewayUserOrder, CHAN_SIZE)

	r.mu.Lock()
	r.userOrderSubs[ch] = struct{}{}
	r.mu.Unlock()

	return ch
}

func (r *RealtimeService) UserTradeStream() <-chan GatewayUserTrade {
	ch := make(chan GatewayUserTrade, CHAN_SIZE)

	r.mu.Lock()
	r.userTradeSubs[ch] = struct{}{}
	r.mu.Unlock()

	return ch
}
