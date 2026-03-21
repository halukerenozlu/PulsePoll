package routes

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestErrorHandlerFormatsFiberErrorEnvelope(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	app.Get("/bad", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	})

	req := httptest.NewRequest("GET", "/bad", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}

	var body errorEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if body.Error.Code != "BAD_REQUEST" {
		t.Fatalf("expected error.code BAD_REQUEST, got %q", body.Error.Code)
	}
	if body.Error.Message != "invalid request" {
		t.Fatalf("expected error.message %q, got %q", "invalid request", body.Error.Message)
	}
}

func TestErrorHandlerFormatsUnexpectedErrorEnvelope(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	app.Get("/boom", func(c *fiber.Ctx) error {
		return errors.New("boom")
	})

	req := httptest.NewRequest("GET", "/boom", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", fiber.StatusInternalServerError, resp.StatusCode)
	}

	var body errorEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if body.Error.Code != "INTERNAL_SERVER_ERROR" {
		t.Fatalf("expected error.code INTERNAL_SERVER_ERROR, got %q", body.Error.Code)
	}
	if body.Error.Message != "internal server error" {
		t.Fatalf("expected generic internal error message, got %q", body.Error.Message)
	}
}
