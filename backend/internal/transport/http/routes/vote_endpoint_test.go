package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"PulsePoll/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestVoteEndpointRegisteredUserSuccessIncrementsVoteCount(t *testing.T) {
	db := newVoteEndpointTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newVoteEndpointTestApp(db, redisClient)
	token, userID := registerVoteEndpointUser(t, app, "vote-registered@example.com")
	surveyID, optionID := createVoteEndpointSurvey(t, db, voteEndpointSurveyConfig{
		CreatorID:       userID,
		MaxVotesPerUser: 1,
	})

	resp := voteEndpointJSONRequest(t, app, surveyID, optionID, nil, token)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
	assertOKResponse(t, resp)
	if count := optionVoteCount(t, db, optionID); count != 1 {
		t.Fatalf("expected vote_count 1, got %d", count)
	}
}

func TestVoteEndpointGuestWithConsentSuccess(t *testing.T) {
	db := newVoteEndpointTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newVoteEndpointTestApp(db, redisClient)
	_, creatorID := registerVoteEndpointUser(t, app, "vote-guest-creator@example.com")
	surveyID, optionID := createVoteEndpointSurvey(t, db, voteEndpointSurveyConfig{
		CreatorID:       creatorID,
		MaxVotesPerUser: 1,
	})

	resp := voteEndpointJSONRequest(t, app, surveyID, optionID, []*http.Cookie{
		{Name: guestIDCookieName, Value: "guest-with-consent"},
	}, "")
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
	assertOKResponse(t, resp)
}

func TestVoteEndpointGuestWithoutConsentReturnsConsentRequired(t *testing.T) {
	db := newVoteEndpointTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newVoteEndpointTestApp(db, redisClient)
	_, creatorID := registerVoteEndpointUser(t, app, "vote-no-consent-creator@example.com")
	surveyID, optionID := createVoteEndpointSurvey(t, db, voteEndpointSurveyConfig{
		CreatorID:       creatorID,
		MaxVotesPerUser: 1,
	})

	resp := voteEndpointJSONRequest(t, app, surveyID, optionID, nil, "")
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status %d, got %d", fiber.StatusForbidden, resp.StatusCode)
	}
	assertErrorCode(t, resp, "CONSENT_REQUIRED")
}

func TestVoteEndpointExpiredSurveyReturnsPhaseNotVoting(t *testing.T) {
	db := newVoteEndpointTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newVoteEndpointTestApp(db, redisClient)
	token, creatorID := registerVoteEndpointUser(t, app, "vote-expired@example.com")
	surveyID, optionID := createVoteEndpointSurvey(t, db, voteEndpointSurveyConfig{
		CreatorID:       creatorID,
		MaxVotesPerUser: 1,
		VoteEndsAt:      time.Now().UTC().Add(-2 * time.Hour),
		ResultsEndsAt:   time.Now().UTC().Add(-time.Hour),
		RetentionEndsAt: time.Now().UTC().Add(time.Hour),
	})

	resp := voteEndpointJSONRequest(t, app, surveyID, optionID, nil, token)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status %d, got %d", fiber.StatusForbidden, resp.StatusCode)
	}
	assertErrorCode(t, resp, "PHASE_NOT_VOTING")
}

func TestVoteEndpointMaxVotesExceeded(t *testing.T) {
	db := newVoteEndpointTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newVoteEndpointTestApp(db, redisClient)
	token, creatorID := registerVoteEndpointUser(t, app, "vote-max@example.com")
	surveyID, optionID := createVoteEndpointSurvey(t, db, voteEndpointSurveyConfig{
		CreatorID:       creatorID,
		MaxVotesPerUser: 1,
	})

	first := voteEndpointJSONRequest(t, app, surveyID, optionID, nil, token)
	if first.StatusCode != fiber.StatusOK {
		first.Body.Close()
		t.Fatalf("first vote: expected status %d, got %d", fiber.StatusOK, first.StatusCode)
	}
	first.Body.Close()

	second := voteEndpointJSONRequest(t, app, surveyID, optionID, nil, token)
	defer second.Body.Close()

	if second.StatusCode != fiber.StatusForbidden && second.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status %d or %d, got %d", fiber.StatusForbidden, fiber.StatusBadRequest, second.StatusCode)
	}
}

func TestVoteEndpointPrivatePINWithoutPINOKReturnsPINRequired(t *testing.T) {
	db := newVoteEndpointTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newVoteEndpointTestApp(db, redisClient)
	token, creatorID := registerVoteEndpointUser(t, app, "vote-pin@example.com")
	surveyID, optionID := createVoteEndpointSurvey(t, db, voteEndpointSurveyConfig{
		CreatorID:       creatorID,
		Visibility:      "private_pin",
		AccessPIN:       "1234",
		MaxVotesPerUser: 1,
	})

	resp := voteEndpointJSONRequest(t, app, surveyID, optionID, nil, token)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status %d, got %d", fiber.StatusForbidden, resp.StatusCode)
	}
	assertErrorCode(t, resp, "PIN_REQUIRED")
}

func TestVoteEndpointSurveyNotFound(t *testing.T) {
	db := newVoteEndpointTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newVoteEndpointTestApp(db, redisClient)
	token, _ := registerVoteEndpointUser(t, app, "vote-missing@example.com")

	resp := voteEndpointJSONRequest(t, app, uuid.NewString(), uuid.NewString(), nil, token)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected status %d, got %d", fiber.StatusNotFound, resp.StatusCode)
	}
}

func TestVoteEndpointInvalidOptionID(t *testing.T) {
	db := newVoteEndpointTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newVoteEndpointTestApp(db, redisClient)
	token, creatorID := registerVoteEndpointUser(t, app, "vote-invalid-option@example.com")
	surveyID, _ := createVoteEndpointSurvey(t, db, voteEndpointSurveyConfig{
		CreatorID:       creatorID,
		MaxVotesPerUser: 1,
	})

	resp := voteEndpointJSONRequest(t, app, surveyID, uuid.NewString(), nil, token)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest && resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected status %d or %d, got %d", fiber.StatusBadRequest, fiber.StatusNotFound, resp.StatusCode)
	}
}

type voteEndpointSurveyConfig struct {
	CreatorID       string
	Visibility      string
	AccessPIN       string
	MaxVotesPerUser int
	VoteEndsAt      time.Time
	ResultsEndsAt   time.Time
	RetentionEndsAt time.Time
}

func newVoteEndpointTestApp(db *gorm.DB, redisClient *goredis.Client) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	RegisterAuthRoutes(app, db, config.AuthConfig{
		JWTSecret:           "vote-endpoint-test-secret",
		AccessTokenTTLMin:   15,
		RefreshTokenTTLHour: 24,
		RefreshCookieName:   "refresh_token",
		RefreshCookieSecure: false,
	})
	RegisterVoteRoutes(app, db, redisClient, "vote-endpoint-test-secret")
	return app
}

func newVoteEndpointTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db := newAuthTestDB(t)
	if err := db.AutoMigrate(&surveyModel{}, &optionDB{}); err != nil {
		t.Fatalf("migrate vote endpoint test tables: %v", err)
	}
	return db
}

func registerVoteEndpointUser(t *testing.T, app *fiber.App, email string) (string, string) {
	t.Helper()

	resp := registerAuthUser(t, app, email, "StrongPass123!")
	defer resp.Body.Close()

	var body authResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode register response: %v", err)
	}
	if body.AccessToken == "" {
		t.Fatal("expected access token from register")
	}
	return body.AccessToken, body.User.ID
}

func createVoteEndpointSurvey(t *testing.T, db *gorm.DB, cfg voteEndpointSurveyConfig) (string, string) {
	t.Helper()

	now := time.Now().UTC()
	visibility := cfg.Visibility
	if visibility == "" {
		visibility = "public"
	}
	maxVotes := cfg.MaxVotesPerUser
	if maxVotes == 0 {
		maxVotes = 1
	}
	voteEndsAt := cfg.VoteEndsAt
	if voteEndsAt.IsZero() {
		voteEndsAt = now.Add(time.Hour)
	}
	resultsEndsAt := cfg.ResultsEndsAt
	if resultsEndsAt.IsZero() {
		resultsEndsAt = now.Add(2 * time.Hour)
	}
	retentionEndsAt := cfg.RetentionEndsAt
	if retentionEndsAt.IsZero() {
		retentionEndsAt = now.Add(2 * time.Hour)
	}

	record := surveyModel{
		ID:                  uuid.NewString(),
		CreatorID:           cfg.CreatorID,
		Title:               "Vote endpoint survey",
		Visibility:          visibility,
		ResultsMode:         "open_live",
		MaxVotesPerUser:     maxVotes,
		AllowVoteChangeOnce: false,
		VoteEndsAt:          voteEndsAt,
		ResultsEndsAt:       resultsEndsAt,
		RetentionEndsAt:     retentionEndsAt,
		ModerationStatus:    "approved",
	}
	if cfg.AccessPIN != "" {
		hashBytes, err := bcrypt.GenerateFromPassword([]byte(cfg.AccessPIN), bcrypt.DefaultCost)
		if err != nil {
			t.Fatalf("hash access pin: %v", err)
		}
		hash := string(hashBytes)
		record.AccessPinHash = &hash
	}
	if err := db.Create(&record).Error; err != nil {
		t.Fatalf("create vote endpoint survey: %v", err)
	}

	option := optionDB{
		ID:       uuid.NewString(),
		SurveyID: record.ID,
		Text:     "Option A",
		Position: 1,
	}
	if err := db.Create(&option).Error; err != nil {
		t.Fatalf("create vote endpoint option: %v", err)
	}
	return record.ID, option.ID
}

func voteEndpointJSONRequest(
	t *testing.T,
	app *fiber.App,
	surveyID string,
	optionID string,
	cookies []*http.Cookie,
	bearerToken string,
) *http.Response {
	t.Helper()

	body := fmt.Sprintf(`{"option_id":%q}`, optionID)
	return authJSONRequestWithBearer(t, app, http.MethodPost, "/api/v1/surveys/"+surveyID+"/vote", body, cookies, bearerToken)
}

func assertOKResponse(t *testing.T, resp *http.Response) {
	t.Helper()

	var body map[string]bool
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode ok response: %v", err)
	}
	if !body["ok"] {
		t.Fatal("expected ok to be true")
	}
}

func assertErrorCode(t *testing.T, resp *http.Response, want string) {
	t.Helper()

	var body errorEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if body.Error.Code != want {
		t.Fatalf("expected error code %q, got %q", want, body.Error.Code)
	}
}

func optionVoteCount(t *testing.T, db *gorm.DB, optionID string) int64 {
	t.Helper()

	var option optionDB
	if err := db.Where("id = ?", optionID).First(&option).Error; err != nil {
		t.Fatalf("load option vote count: %v", err)
	}
	return option.VoteCount
}
