package routes

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAcceptConsentEndpointSetsGuestCookie(t *testing.T) {
	app := newConsentTestApp()

	resp := authJSONRequest(t, app, http.MethodPost, "/api/v1/consent/accept", `{}`, nil)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	var body map[string]bool
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	ok, exists := body["ok"]
	if !exists {
		t.Fatal("expected ok field in response")
	}
	if !ok {
		t.Fatal("expected ok to be true")
	}

	cookie := requireConsentCookie(t, resp)
	if !cookie.HttpOnly {
		t.Fatal("expected guest_id cookie to be HttpOnly")
	}
}

func TestAcceptConsentEndpointIsIdempotent(t *testing.T) {
	app := newConsentTestApp()

	first := authJSONRequest(t, app, http.MethodPost, "/api/v1/consent/accept", `{}`, nil)
	firstCookie := requireConsentCookie(t, first)
	first.Body.Close()

	second := authJSONRequest(t, app, http.MethodPost, "/api/v1/consent/accept", `{}`, []*http.Cookie{firstCookie})
	defer second.Body.Close()

	if second.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, second.StatusCode)
	}

	var body map[string]bool
	if err := json.NewDecoder(second.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	ok, exists := body["ok"]
	if !exists {
		t.Fatal("expected ok field in response")
	}
	if !ok {
		t.Fatal("expected ok to be true")
	}
}

func newConsentTestApp() *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	RegisterConsentRoutes(app)
	return app
}

func requireConsentCookie(t *testing.T, resp *http.Response) *http.Cookie {
	t.Helper()

	for _, cookie := range resp.Cookies() {
		if cookie.Name == guestIDCookieName && cookie.Value != "" {
			return cookie
		}
	}
	t.Fatalf("expected response cookie %q", guestIDCookieName)
	return nil
}
