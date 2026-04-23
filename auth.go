package projectx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type loginRequest struct {
	Username string `json:"username"`
	APIKey   string `json:"apiKey"`
}

type loginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expiresIn"` // seconds
}

func (c *ProjectXClient) login(ctx context.Context) error {
	req := loginRequest{
		Username: c.username,
		APIKey:   c.apiKey,
	}

	var resp loginResponse
	err := c.PostNoAuth(ctx, "/Auth/loginKey", req, &resp)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	c.mu.Lock()
	c.token = resp.Token
	c.expiresAt = time.Now().Add(time.Duration(resp.ExpiresIn-30) * time.Second) // refresh early
	c.mu.Unlock()

	return nil
}

func (c *ProjectXClient) refreshToken(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// clear token so next getToken() logs in
	c.token = ""
	c.expiresAt = time.Time{}

	err := c.login(ctx)
	return err
}

// setTokenForTest injects a token and bypasses login.
// Only compiled in tests.
func (c *ProjectXClient) setTokenForTest(token string) {
	c.token = token
	c.expiresAt = time.Now().Add(24 * time.Hour) // make it valid
}

func (c *ProjectXClient) negotiate(ctx context.Context, hubHttp string, token string) (string, error) {
	url := fmt.Sprintf("%s/negotiate?negotiateVersion=1", hubHttp)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var res struct {
		ConnectionToken string `json:"connectionToken"`
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}

	return res.ConnectionToken, nil
}
