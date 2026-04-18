package projectx

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestPost_Success(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Fatalf("missing auth header")
		}

		resp := map[string]string{"status": "ok"}
		json.NewEncoder(w).Encode(resp)
	}

	client, closeServer := newTestClient(handler)
	defer closeServer()

	var result map[string]string

	err := client.doPost(
		context.Background(),
		"/test",
		"test-token",
		map[string]string{"hello": "world"},
		&result,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result["status"] != "ok" {
		t.Fatalf("unexpected response: %+v", result)
	}
}

func TestPost_APIError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`bad request`))
	}

	client, closeServer := newTestClient(handler)
	defer closeServer()

	err := client.PostNoAuth(
		context.Background(),
		"/test",
		nil,
		nil,
	)

	if err == nil {
		t.Fatal("expected error")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}

	if apiErr.StatusCode != 400 {
		t.Fatalf("wrong status: %d", apiErr.StatusCode)
	}
}

func TestPost_InvalidJSON(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`not-json`))
	}

	client, closeServer := newTestClient(handler)
	defer closeServer()

	var result map[string]string

	err := client.PostNoAuth(
		context.Background(),
		"/test",
		nil,
		&result,
	)

	if err == nil {
		t.Fatal("expected unmarshal error")
	}
}

func TestPost_ContextTimeout(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}

	client, closeServer := newTestClient(handler)
	defer closeServer()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err := client.PostNoAuth(ctx, "/test", nil, nil)

	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestPost_SendsBody(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)

		if body["foo"] != "bar" {
			t.Fatalf("body not received correctly")
		}
	}

	client, closeServer := newTestClient(handler)
	defer closeServer()

	err := client.PostNoAuth(
		context.Background(),
		"/test",
		map[string]string{"foo": "bar"},
		nil,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
