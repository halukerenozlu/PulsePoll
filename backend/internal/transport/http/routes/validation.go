package routes

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// requireUUIDPathParam validates a path parameter is a non-empty valid UUID.
// On failure it writes the error response and returns errResponseSent.
// Callers must return nil (not the error) to let Fiber send the buffered response.
func requireUUIDPathParam(c *fiber.Ctx, name string) (string, error) {
	value := strings.TrimSpace(c.Params(name))
	if value == "" {
		writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", fmt.Sprintf("%s path parameter is required", name))
		return "", errResponseSent
	}
	if _, err := uuid.Parse(value); err != nil {
		writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", fmt.Sprintf("%s path parameter must be a valid UUID", name))
		return "", errResponseSent
	}
	return value, nil
}

// validateUUIDField validates a request body field is a valid UUID.
// On failure it writes the error response and returns errResponseSent.
// Callers must return nil (not the error) to let Fiber send the buffered response.
func validateUUIDField(c *fiber.Ctx, fieldName, value string) error {
	if _, err := uuid.Parse(value); err != nil {
		writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", fmt.Sprintf("%s must be a valid UUID", fieldName))
		return errResponseSent
	}
	return nil
}
