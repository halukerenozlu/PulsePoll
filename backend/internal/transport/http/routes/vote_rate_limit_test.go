package routes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestVoteRateLimitKeyShape(t *testing.T) {
	got := voteRateLimitKey("127.0.0.1")
	want := "rl:ip:127.0.0.1:vote"
	if got != want {
		t.Fatalf("voteRateLimitKey() = %q, want %q", got, want)
	}
}

func TestEnforceVoteRateLimitWithFuncsAllowedPath(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})

	app.Post("/rate", func(c *fiber.Ctx) error {
		err := enforceVoteRateLimitWithFuncs(
			c,
			"rl:ip:127.0.0.1:vote",
			30,
			func(_ context.Context, _ string) (int64, error) {
				return 1, nil
			},
		)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"ok": true})
	})

	req := httptest.NewRequest("POST", "/rate", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}

func TestEnforceVoteRateLimitWithFuncsBlockedPath(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})

	app.Post("/rate", func(c *fiber.Ctx) error {
		return enforceVoteRateLimitWithFuncs(
			c,
			"rl:ip:127.0.0.1:vote",
			30,
			func(_ context.Context, _ string) (int64, error) {
				return 31, nil
			},
		)
	})

	req := httptest.NewRequest("POST", "/rate", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusTooManyRequests {
		t.Fatalf("expected status %d, got %d", fiber.StatusTooManyRequests, resp.StatusCode)
	}

	var body errorEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	if body.Error.Code != "TOO_MANY_REQUESTS" {
		t.Fatalf("expected TOO_MANY_REQUESTS code, got %q", body.Error.Code)
	}
}

func TestEnforceVoteRateLimitWithFuncsRedisError(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})

	app.Post("/rate", func(c *fiber.Ctx) error {
		return enforceVoteRateLimitWithFuncs(
			c,
			"rl:ip:127.0.0.1:vote",
			30,
			func(_ context.Context, _ string) (int64, error) {
				return 0, errors.New("redis unavailable")
			},
		)
	})

	req := httptest.NewRequest("POST", "/rate", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", fiber.StatusInternalServerError, resp.StatusCode)
	}
}
