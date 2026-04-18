package projectx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *ProjectXClient) PostNoAuth(
	ctx context.Context,
	endpoint string,
	reqBody any,
	respBody any,
) error {
	return c.doPost(ctx, endpoint, "", reqBody, respBody)
}

func (c *ProjectXClient) Post(
	ctx context.Context,
	endpoint string,
	reqBody any,
	respBody any,
) error {

	// Get valid token (auto refresh)
	token, err := c.getToken(ctx)
	if err != nil {
		return err
	}

	// First attempt
	err = c.doPost(ctx, endpoint, token, reqBody, respBody)

	// If unauthorized, refresh token and retry once
	if apiErr, ok := err.(*APIError); ok && apiErr.StatusCode == http.StatusUnauthorized {
		// force refresh token
		if err := c.refreshToken(ctx); err != nil {
			return err
		}

		token, err = c.getToken(ctx)
		if err != nil {
			return err
		}

		return c.doPost(ctx, endpoint, token, reqBody, respBody)
	}

	return err
}

func (c *ProjectXClient) doPost(
	ctx context.Context,
	endpoint string,
	token string,
	reqBody any,
	respBody any,
) error {

	// Marshal request body
	var bodyReader io.Reader
	if reqBody != nil {
		b, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseUrl+endpoint,
		bodyReader,
	)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	if respBody != nil {
		if err := json.Unmarshal(bodyBytes, respBody); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}

	return nil
}
