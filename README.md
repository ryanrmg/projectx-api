# projectx-api

Unofficial Go client for the ProjectX Gateway API.

## Install

go get github.com/ryangess/projectx-api

## Example

```go
client := projectx.NewProjectXClient("https://gateway-api-demo.s2f.projectx.com", "username", "api-key")

accounts, err := client.Accounts.Search(ctx)```
