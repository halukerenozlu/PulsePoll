package routes

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestSurveyPathValidationReturnsBadRequest(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	RegisterSurveyRoutes(app, nil, "test-secret")

	req := httptest.NewRequest("GET", "/api/v1/surveys/not-a-uuid", nil)
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
		t.Fatalf("expected BAD_REQUEST code, got %q", body.Error.Code)
	}
	if !strings.Contains(body.Error.Message, "id path parameter") {
		t.Fatalf("expected path parameter message, got %q", body.Error.Message)
	}
}

func TestFeedQueryValidationReturnsBadRequest(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	RegisterSurveyRoutes(app, nil, "test-secret")

	tests := []struct {
		name     string
		url      string
		contains string
	}{
		{name: "invalid sort", url: "/api/v1/feed?sort=trending", contains: "sort=new"},
		{name: "invalid visibility", url: "/api/v1/feed?visibility=private_pin", contains: "visibility=public"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
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
				t.Fatalf("expected BAD_REQUEST code, got %q", body.Error.Code)
			}
			if !strings.Contains(body.Error.Message, tt.contains) {
				t.Fatalf("expected message containing %q, got %q", tt.contains, body.Error.Message)
			}
		})
	}
}

func TestUUIDFieldValidationReturnsBadRequest(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	app.Post("/validate", func(c *fiber.Ctx) error {
		var req struct {
			OptionID string `json:"option_id"`
		}
		if err := c.BodyParser(&req); err != nil {
			return err
		}
		if err := validateUUIDField(c, "option_id", req.OptionID); err != nil {
			if err == errResponseSent {
				return nil
			}
			return err
		}
		return c.JSON(fiber.Map{"ok": true})
	})

	req := httptest.NewRequest("POST", "/validate", bytes.NewBufferString(`{"option_id":"not-a-uuid"}`))
	req.Header.Set("Content-Type", "application/json")
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
		t.Fatalf("expected BAD_REQUEST code, got %q", body.Error.Code)
	}
	if !strings.Contains(body.Error.Message, "option_id") {
		t.Fatalf("expected option_id mention in message, got %q", body.Error.Message)
	}
}
