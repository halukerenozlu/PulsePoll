package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestVoteChangeEndpointRegisteredUserSuccessUpdatesVoteCounts(t *testing.T) {
	db := newVoteEndpointTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newVoteEndpointTestApp(db, redisClient)
	token, userID := registerVoteEndpointUser(t, app, "change-registered@example.com")
	surveyID, oldOptionID, newOptionID := createVoteChangeSurvey(t, db, voteChangeSurveyConfig{
		CreatorID:           userID,
		MaxVotesPerUser:     1,
		AllowVoteChangeOnce: true,
		OldOptionVoteCount:  1,
	})
	setVoteChangeReceipt(t, redisClient, surveyID, voterIdentity{UserID: userID}, voteReceipt{
		VotesUsed:    1,
		LastOptionID: oldOptionID,
	})

	resp := voteChangeEndpointJSONRequest(t, app, surveyID, newOptionID, nil, token)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
	assertOKResponse(t, resp)
	if count := optionVoteCount(t, db, oldOptionID); count != 0 {
		t.Fatalf("expected old option vote_count 0, got %d", count)
	}
	if count := optionVoteCount(t, db, newOptionID); count != 1 {
		t.Fatalf("expected new option vote_count 1, got %d", count)
	}
}

func TestVoteChangeEndpointGuestWithConsentSuccess(t *testing.T) {
	db := newVoteEndpointTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newVoteEndpointTestApp(db, redisClient)
	_, creatorID := registerVoteEndpointUser(t, app, "change-guest-creator@example.com")
	guestID := "guest-change-success"
	surveyID, oldOptionID, newOptionID := createVoteChangeSurvey(t, db, voteChangeSurveyConfig{
		CreatorID:           creatorID,
		MaxVotesPerUser:     1,
		AllowVoteChangeOnce: true,
		OldOptionVoteCount:  1,
	})
	setVoteChangeReceipt(t, redisClient, surveyID, voterIdentity{GuestID: guestID, IsGuest: true}, voteReceipt{
		VotesUsed:    1,
		LastOptionID: oldOptionID,
	})

	resp := voteChangeEndpointJSONRequest(t, app, surveyID, newOptionID, []*http.Cookie{
		{Name: guestIDCookieName, Value: guestID},
	}, "")
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
	assertOKResponse(t, resp)
}

func TestVoteChangeEndpointChangeAlreadyUsed(t *testing.T) {
	db := newVoteEndpointTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newVoteEndpointTestApp(db, redisClient)
	token, userID := registerVoteEndpointUser(t, app, "change-used@example.com")
	surveyID, oldOptionID, newOptionID := createVoteChangeSurvey(t, db, voteChangeSurveyConfig{
		CreatorID:           userID,
		MaxVotesPerUser:     1,
		AllowVoteChangeOnce: true,
		OldOptionVoteCount:  1,
	})
	setVoteChangeReceipt(t, redisClient, surveyID, voterIdentity{UserID: userID}, voteReceipt{
		VotesUsed:    1,
		LastOptionID: oldOptionID,
		ChangeUsed:   true,
	})

	resp := voteChangeEndpointJSONRequest(t, app, surveyID, newOptionID, nil, token)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status %d, got %d", fiber.StatusForbidden, resp.StatusCode)
	}
	assertErrorCode(t, resp, "VOTE_CHANGE_NOT_ALLOWED")
}

func TestVoteChangeEndpointAllowVoteChangeOnceFalse(t *testing.T) {
	db := newVoteEndpointTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newVoteEndpointTestApp(db, redisClient)
	token, userID := registerVoteEndpointUser(t, app, "change-disabled@example.com")
	surveyID, _, newOptionID := createVoteChangeSurvey(t, db, voteChangeSurveyConfig{
		CreatorID:           userID,
		MaxVotesPerUser:     1,
		AllowVoteChangeOnce: false,
	})

	resp := voteChangeEndpointJSONRequest(t, app, surveyID, newOptionID, nil, token)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status %d, got %d", fiber.StatusForbidden, resp.StatusCode)
	}
	assertErrorCode(t, resp, "VOTE_CHANGE_NOT_ALLOWED")
}

func TestVoteChangeEndpointMaxVotesPerUserGreaterThanOne(t *testing.T) {
	db := newVoteEndpointTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newVoteEndpointTestApp(db, redisClient)
	token, userID := registerVoteEndpointUser(t, app, "change-max-votes@example.com")
	surveyID, _, newOptionID := createVoteChangeSurvey(t, db, voteChangeSurveyConfig{
		CreatorID:           userID,
		MaxVotesPerUser:     2,
		AllowVoteChangeOnce: false,
	})

	resp := voteChangeEndpointJSONRequest(t, app, surveyID, newOptionID, nil, token)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status %d, got %d", fiber.StatusForbidden, resp.StatusCode)
	}
	assertErrorCode(t, resp, "VOTE_CHANGE_NOT_ALLOWED")
}

func TestVoteChangeEndpointExpiredSurveyReturnsPhaseNotVoting(t *testing.T) {
	db := newVoteEndpointTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newVoteEndpointTestApp(db, redisClient)
	token, userID := registerVoteEndpointUser(t, app, "change-expired@example.com")
	surveyID, _, newOptionID := createVoteChangeSurvey(t, db, voteChangeSurveyConfig{
		CreatorID:           userID,
		MaxVotesPerUser:     1,
		AllowVoteChangeOnce: true,
		VoteEndsAt:          time.Now().UTC().Add(-2 * time.Hour),
		ResultsEndsAt:       time.Now().UTC().Add(-time.Hour),
		RetentionEndsAt:     time.Now().UTC().Add(time.Hour),
	})

	resp := voteChangeEndpointJSONRequest(t, app, surveyID, newOptionID, nil, token)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status %d, got %d", fiber.StatusForbidden, resp.StatusCode)
	}
	assertErrorCode(t, resp, "PHASE_NOT_VOTING")
}

func TestVoteChangeEndpointGuestWithoutConsentReturnsConsentRequired(t *testing.T) {
	db := newVoteEndpointTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newVoteEndpointTestApp(db, redisClient)
	_, creatorID := registerVoteEndpointUser(t, app, "change-no-consent-creator@example.com")
	surveyID, _, newOptionID := createVoteChangeSurvey(t, db, voteChangeSurveyConfig{
		CreatorID:           creatorID,
		MaxVotesPerUser:     1,
		AllowVoteChangeOnce: true,
	})

	resp := voteChangeEndpointJSONRequest(t, app, surveyID, newOptionID, nil, "")
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status %d, got %d", fiber.StatusForbidden, resp.StatusCode)
	}
	assertErrorCode(t, resp, "CONSENT_REQUIRED")
}

func TestVoteChangeEndpointPrivatePINWithoutPINOKReturnsPINRequired(t *testing.T) {
	db := newVoteEndpointTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newVoteEndpointTestApp(db, redisClient)
	token, userID := registerVoteEndpointUser(t, app, "change-pin@example.com")
	surveyID, oldOptionID, newOptionID := createVoteChangeSurvey(t, db, voteChangeSurveyConfig{
		CreatorID:           userID,
		Visibility:          "private_pin",
		AccessPIN:           "1234",
		MaxVotesPerUser:     1,
		AllowVoteChangeOnce: true,
		OldOptionVoteCount:  1,
	})
	setVoteChangeReceipt(t, redisClient, surveyID, voterIdentity{UserID: userID}, voteReceipt{
		VotesUsed:    1,
		LastOptionID: oldOptionID,
	})

	resp := voteChangeEndpointJSONRequest(t, app, surveyID, newOptionID, nil, token)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status %d, got %d", fiber.StatusForbidden, resp.StatusCode)
	}
	assertErrorCode(t, resp, "PIN_REQUIRED")
}

func TestVoteChangeEndpointSurveyNotFound(t *testing.T) {
	db := newVoteEndpointTestDB(t)
	redisClient := newPINVerifyTestRedis(t)
	app := newVoteEndpointTestApp(db, redisClient)
	token, _ := registerVoteEndpointUser(t, app, "change-missing@example.com")

	resp := voteChangeEndpointJSONRequest(t, app, uuid.NewString(), uuid.NewString(), nil, token)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected status %d, got %d", fiber.StatusNotFound, resp.StatusCode)
	}
}

type voteChangeSurveyConfig struct {
	CreatorID           string
	Visibility          string
	AccessPIN           string
	MaxVotesPerUser     int
	AllowVoteChangeOnce bool
	VoteEndsAt          time.Time
	ResultsEndsAt       time.Time
	RetentionEndsAt     time.Time
	OldOptionVoteCount  int64
}

func createVoteChangeSurvey(t *testing.T, db *gorm.DB, cfg voteChangeSurveyConfig) (string, string, string) {
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
		Title:               "Vote change endpoint survey",
		Visibility:          visibility,
		ResultsMode:         "open_live",
		MaxVotesPerUser:     maxVotes,
		AllowVoteChangeOnce: cfg.AllowVoteChangeOnce,
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
		t.Fatalf("create vote change survey: %v", err)
	}

	oldOption := optionDB{
		ID:        uuid.NewString(),
		SurveyID:  record.ID,
		Text:      "Old option",
		Position:  1,
		VoteCount: cfg.OldOptionVoteCount,
	}
	newOption := optionDB{
		ID:       uuid.NewString(),
		SurveyID: record.ID,
		Text:     "New option",
		Position: 2,
	}
	if err := db.Create(&oldOption).Error; err != nil {
		t.Fatalf("create old option: %v", err)
	}
	if err := db.Create(&newOption).Error; err != nil {
		t.Fatalf("create new option: %v", err)
	}
	return record.ID, oldOption.ID, newOption.ID
}

func voteChangeEndpointJSONRequest(
	t *testing.T,
	app *fiber.App,
	surveyID string,
	newOptionID string,
	cookies []*http.Cookie,
	bearerToken string,
) *http.Response {
	t.Helper()

	body := fmt.Sprintf(`{"new_option_id":%q}`, newOptionID)
	return authJSONRequestWithBearer(t, app, http.MethodPut, "/api/v1/surveys/"+surveyID+"/vote", body, cookies, bearerToken)
}

func setVoteChangeReceipt(t *testing.T, redisClient *goredis.Client, surveyID string, identity voterIdentity, receipt voteReceipt) {
	t.Helper()

	raw, err := json.Marshal(receipt)
	if err != nil {
		t.Fatalf("marshal vote receipt: %v", err)
	}
	if err := redisClient.Set(context.Background(), voteReceiptKey(surveyID, identity), raw, time.Hour).Err(); err != nil {
		t.Fatalf("set vote receipt: %v", err)
	}
}
