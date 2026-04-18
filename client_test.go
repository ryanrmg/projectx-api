package projectx

import (
	"net/http"
	"net/http/httptest"
	"time"
)

func newTestClient(handler http.HandlerFunc) (*ProjectXClient, func()) {
	server := httptest.NewServer(handler)

	client := NewProjectXClient(server.URL, "username", "test-key")
	client.httpClient = &http.Client{Timeout: 5 * time.Second}

	return client, server.Close
}
