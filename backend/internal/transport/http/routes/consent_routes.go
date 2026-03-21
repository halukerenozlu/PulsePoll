package routes

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	guestIDCookieName = "guest_id"
	guestIDTTL        = 48 * time.Hour
)

func RegisterConsentRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	consent := api.Group("/consent")
	consent.Post("/accept", acceptConsent)
}

func RequireGuestConsentForGuest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if isAuthenticatedRequest(c) {
			return c.Next()
		}

		if c.Cookies(guestIDCookieName) == "" {
			return writeError(c, fiber.StatusForbidden, "CONSENT_REQUIRED", "guest voting requires consent cookie")
		}

		return c.Next()
	}
}

func acceptConsent(c *fiber.Ctx) error {
	guestID, err := newGuestID()
	if err != nil {
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to generate guest id")
	}

	c.Cookie(&fiber.Cookie{
		Name:     guestIDCookieName,
		Value:    guestID,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
		Expires:  time.Now().UTC().Add(guestIDTTL),
	})

	return c.JSON(fiber.Map{"ok": true})
}

func newGuestID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func isAuthenticatedRequest(c *fiber.Ctx) bool {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return false
	}

	parts := strings.SplitN(authHeader, " ", 2)
	return len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") && strings.TrimSpace(parts[1]) != ""
}
