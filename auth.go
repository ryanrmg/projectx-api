package projectx

import (
	"context"
	"fmt"
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
