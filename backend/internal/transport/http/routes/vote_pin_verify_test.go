package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestPINVerifyEndpointSuccessSetsGuestPINOK(t *testing.T) {
	db := newPINVerifyTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newPINVerifyTestApp(db, redisClient)

	surveyID := createPINVerifySurvey(t, db, "1234")
	guestID := "guest-success"

	resp := pinVerifyJSONRequest(t, app, surveyID, guestID, `{"pin":"1234"}`)
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	got, err := redisClient.Get(context.Background(), pinOKKey(surveyID, voterIdentity{
		GuestID: guestID,
		IsGuest: true,
	})).Result()
	if err != nil {
		t.Fatalf("expected pinok key to be set: %v", err)
	}
	if got != "1" {
		t.Fatalf("expected pinok value %q, got %q", "1", got)
	}
}

func TestPINVerifyEndpointWrongPINIncrementsGuestPINFail(t *testing.T) {
	db := newPINVerifyTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newPINVerifyTestApp(db, redisClient)

	surveyID := createPINVerifySurvey(t, db, "1234")
	guestID := "guest-wrong"

	resp := pinVerifyJSONRequest(t, app, surveyID, guestID, `{"pin":"9999"}`)
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status %d, got %d", fiber.StatusForbidden, resp.StatusCode)
	}

	var body errorEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	if body.Error.Code != "PIN_REQUIRED" {
		t.Fatalf("expected PIN_REQUIRED code, got %q", body.Error.Code)
	}

	count := pinFailCount(t, redisClient, surveyID, guestID)
	if count != 1 {
		t.Fatalf("expected pinfail counter 1, got %d", count)
	}
}

func TestPINVerifyEndpointBruteForceReturnsTooManyRequests(t *testing.T) {
	db := newPINVerifyTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newPINVerifyTestApp(db, redisClient)

	surveyID := createPINVerifySurvey(t, db, "1234")
	guestID := "guest-bruteforce"

	for i := 0; i < pinFailMaxAttempts; i++ {
		resp := pinVerifyJSONRequest(t, app, surveyID, guestID, `{"pin":"9999"}`)
		if resp.StatusCode != fiber.StatusForbidden {
			resp.Body.Close()
			t.Fatalf("attempt %d: expected status %d, got %d", i+1, fiber.StatusForbidden, resp.StatusCode)
		}
		resp.Body.Close()
	}

	resp := pinVerifyJSONRequest(t, app, surveyID, guestID, `{"pin":"9999"}`)
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

func TestPINVerifyEndpointSurveyNotFound(t *testing.T) {
	db := newPINVerifyTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newPINVerifyTestApp(db, redisClient)

	resp := pinVerifyJSONRequest(t, app, uuid.NewString(), "guest-missing", `{"pin":"1234"}`)
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected status %d, got %d", fiber.StatusNotFound, resp.StatusCode)
	}
}

func newPINVerifyTestApp(db *gorm.DB, redisClient *goredis.Client) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	RegisterVoteRoutes(app, db, redisClient, "pin-verify-test-secret")
	return app
}

func newPINVerifyTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db := newAuthTestDB(t)
	if err := db.AutoMigrate(&surveyModel{}, &optionDB{}); err != nil {
		t.Fatalf("migrate pin verify test tables: %v", err)
	}
	return db
}

func newPINVerifyTestRedis(t *testing.T) *goredis.Client {
	t.Helper()

	addr := os.Getenv("PULSEPOLL_TEST_REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	dbNum := 15
	if raw := os.Getenv("PULSEPOLL_TEST_REDIS_DB"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			t.Fatalf("invalid PULSEPOLL_TEST_REDIS_DB: %v", err)
		}
		dbNum = parsed
	}

	client := goredis.NewClient(&goredis.Options{Addr: addr, DB: dbNum})
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Skipf("pin verify endpoint tests require Redis: %v", err)
	}
	if err := client.FlushDB(context.Background()).Err(); err != nil {
		t.Fatalf("flush test redis db: %v", err)
	}
	t.Cleanup(func() {
		_ = client.FlushDB(context.Background()).Err()
		_ = client.Close()
	})
	return client
}

func createPINVerifySurvey(t *testing.T, db *gorm.DB, pin string) string {
	t.Helper()

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash pin: %v", err)
	}
	hash := string(hashBytes)
	now := time.Now().UTC()
	record := surveyModel{
		ID:                  uuid.NewString(),
		CreatorID:           uuid.NewString(),
		Title:               "PIN verify survey",
		Visibility:          "private_pin",
		AccessPinHash:       &hash,
		ResultsMode:         "open_live",
		MaxVotesPerUser:     1,
		AllowVoteChangeOnce: false,
		VoteEndsAt:          now.Add(time.Hour),
		ResultsEndsAt:       now.Add(2 * time.Hour),
		RetentionEndsAt:     now.Add(2 * time.Hour),
		ModerationStatus:    "approved",
	}
	if err := db.Create(&record).Error; err != nil {
		t.Fatalf("create pin verify survey: %v", err)
	}
	return record.ID
}

func pinVerifyJSONRequest(t *testing.T, app *fiber.App, surveyID string, guestID string, body string) *http.Response {
	t.Helper()

	req := httptest.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/api/v1/surveys/%s/pin/verify", surveyID),
		bytes.NewReader([]byte(body)),
	)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: guestIDCookieName, Value: guestID})

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}
	return resp
}

func pinFailCount(t *testing.T, redisClient *goredis.Client, surveyID string, guestID string) int64 {
	t.Helper()

	raw, err := redisClient.Get(context.Background(), pinFailKey(surveyID, voterIdentity{
		GuestID: guestID,
		IsGuest: true,
	})).Result()
	if err != nil {
		t.Fatalf("get pinfail counter: %v", err)
	}
	count, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		t.Fatalf("parse pinfail counter: %v", err)
	}
	return count
}
