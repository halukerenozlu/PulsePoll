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

type consentErrorEnvelope struct {
	Error consentErrorBody `json:"error"`
}

type consentErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

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
			return c.Status(fiber.StatusForbidden).JSON(consentErrorEnvelope{
				Error: consentErrorBody{
					Code:    "CONSENT_REQUIRED",
					Message: "guest voting requires consent cookie",
				},
			})
		}

		return c.Next()
	}
}

func acceptConsent(c *fiber.Ctx) error {
	guestID, err := newGuestID()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(consentErrorEnvelope{
			Error: consentErrorBody{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "failed to generate guest id",
			},
		})
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
