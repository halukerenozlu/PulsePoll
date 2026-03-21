package routes

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// errResponseSent signals that an HTTP error response has already been written
// to the client. Callers must return nil (not this error) so Fiber sends the
// buffered response instead of invoking its error handler.
var errResponseSent = errors.New("response already sent")

// Sentinel errors for vote transaction operations.
var (
	errOptionNotFound         = errors.New("option not found")
	errPreviousOptionNotFound = errors.New("previous option not found")
	errNewOptionNotFound      = errors.New("new option not found")
)

func isResponseSent(err error) bool {
	return errors.Is(err, errResponseSent)
}

// ErrorHandler returns a Fiber error handler that produces the standard
// error envelope for any unhandled error returned from a handler.
func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			code := statusToCode(fiberErr.Code)
			return writeError(c, fiberErr.Code, code, fiberErr.Message)
		}

		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
	}
}

func writeError(c *fiber.Ctx, status int, code, message string) error {
	if status >= fiber.StatusInternalServerError {
		log.Printf("http 5xx response method=%s path=%s status=%d code=%s message=%s", c.Method(), c.Path(), status, code, message)
	}

	return c.Status(status).JSON(errorEnvelope{
		Error: errorBody{
			Code:    code,
			Message: message,
		},
	})
}

// statusToCode converts an HTTP status code to an uppercase underscore code
// string (e.g. 400 -> "BAD_REQUEST", 404 -> "NOT_FOUND").
func statusToCode(status int) string {
	text := http.StatusText(status)
	if text == "" {
		return "INTERNAL_SERVER_ERROR"
	}
	return strings.ReplaceAll(strings.ToUpper(text), " ", "_")
}
