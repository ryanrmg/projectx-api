# projectx-api

Unofficial Go client for the ProjectX Gateway API.

## Install

go get github.com/ryanrmg/projectx-api

## Example

## Basic 
```go
client := NewProjectXClient(
	"https://api.topstepx.com/api",
	"https://rtc.topstepx.com/hubs/",
	"username",
	"api-key",
)

accounts, err := client.Accounts.Search(
	ctx,
	AccountSearchRequest{OnlyActiveAccounts: true},
)

bars, err := client.Markets.History(
	ctx,
	BarHistoryRequest{
		ContractId:        "CON.F.US.MNQ.M26",
		Live:              false,
		EndTime:           time.Now().UTC().Format(time.RFC3339),
		StartTime:         time.Now().UTC().Add(-time.Duration(24*5) * time.Hour).Format(time.RFC3339),
		Unit:              2,
		UnitNumber:        5,
		Limit:             100,
		IncludePartialBar: false,
	},
)

contracts, err := client.Markets.AvailableContracts(
	ctx,
	AvailableContractRequest{Live: true},
)

```

## Real Time

```go
	// connect
	if err := client.Realtime.Connect(ctx); err != nil {
		t.Fatalf("connect error: %v", err)
	}

	// subscribe
	contract := "CON.F.US.MNQ.M26"
	if err := client.Realtime.SubscribeContractTrades(contract); err != nil {
		t.Fatalf("subscribe error: %v", err)
	}

	t.Log("subscribed — waiting for trade message...")

	trades := client.Realtime.TradesStream()

	for {
		select {

		case msg := <-trades:
			// do something with your message

	}

```
