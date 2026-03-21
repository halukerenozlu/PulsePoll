package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHealthRouteHealthy(t *testing.T) {
	app := fiber.New()
	registerHealthRoute(
		app,
		func(_ context.Context) error { return nil },
		func(_ context.Context) error { return nil },
	)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if body["ok"] != true {
		t.Fatalf("expected ok=true, got %v", body["ok"])
	}
	if body["db"] != "up" {
		t.Fatalf("expected db=up, got %v", body["db"])
	}
	if body["redis"] != "up" {
		t.Fatalf("expected redis=up, got %v", body["redis"])
	}
}

func TestHealthRouteServiceUnavailableWhenDependencyFails(t *testing.T) {
	tests := []struct {
		name     string
		pgErr    error
		redisErr error
		expected map[string]any
	}{
		{
			name:     "postgres down",
			pgErr:    errors.New("postgres down"),
			expected: map[string]any{"ok": false, "db": "down", "redis": "up"},
		},
		{
			name:     "redis down",
			redisErr: errors.New("redis down"),
			expected: map[string]any{"ok": false, "db": "up", "redis": "down"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			registerHealthRoute(
				app,
				func(_ context.Context) error { return tt.pgErr },
				func(_ context.Context) error { return tt.redisErr },
			)

			req := httptest.NewRequest("GET", "/health", nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("app.Test returned error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != fiber.StatusServiceUnavailable {
				t.Fatalf("expected status %d, got %d", fiber.StatusServiceUnavailable, resp.StatusCode)
			}

			var body map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("decode response body: %v", err)
			}

			if body["ok"] != tt.expected["ok"] {
				t.Fatalf("expected ok=%v, got %v", tt.expected["ok"], body["ok"])
			}
			if body["db"] != tt.expected["db"] {
				t.Fatalf("expected db=%v, got %v", tt.expected["db"], body["db"])
			}
			if body["redis"] != tt.expected["redis"] {
				t.Fatalf("expected redis=%v, got %v", tt.expected["redis"], body["redis"])
			}
		})
	}
}
