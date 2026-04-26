# projectx-api

Unofficial Go client for the ProjectX Gateway API.

## Install

go get github.com/ryangess/projectx-api

## Example

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
