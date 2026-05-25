package routes

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"PulsePoll/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestReportEndpointRegisteredUserSuccessCreatesReport(t *testing.T) {
	db := newReportEndpointTestDB(t)
	app := newReportEndpointTestApp(db)
	token, userID := registerReportEndpointUser(t, app, "report-user@example.com")
	surveyID := createReportEndpointSurvey(t, db, userID)

	resp := authJSONRequestWithBearer(t, app, http.MethodPost, "/api/v1/surveys/"+surveyID+"/report", `{
		"reason":"abuse",
		"details":"contains harmful content"
	}`, nil, token)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status %d, got %d", fiber.StatusCreated, resp.StatusCode)
	}
	assertReportOKResponse(t, resp)

	var report reportModel
	if err := db.Where("survey_id = ?", surveyID).First(&report).Error; err != nil {
		t.Fatalf("load report: %v", err)
	}
	if report.ReporterUserID == nil || *report.ReporterUserID != userID {
		t.Fatalf("expected reporter_user_id %q, got %#v", userID, report.ReporterUserID)
	}
	if report.ReporterGuestID != nil {
		t.Fatalf("expected no reporter_guest_id, got %#v", report.ReporterGuestID)
	}
	if report.Reason != "abuse" {
		t.Fatalf("expected reason %q, got %q", "abuse", report.Reason)
	}
}

func TestReportEndpointGuestWithConsentSuccessCreatesReport(t *testing.T) {
	db := newReportEndpointTestDB(t)
	app := newReportEndpointTestApp(db)
	_, creatorID := registerReportEndpointUser(t, app, "report-guest-creator@example.com")
	surveyID := createReportEndpointSurvey(t, db, creatorID)
	guestID := "guest-report"

	resp := authJSONRequest(t, app, http.MethodPost, "/api/v1/surveys/"+surveyID+"/report", `{
		"reason":"spam"
	}`, []*http.Cookie{{Name: guestIDCookieName, Value: guestID}})
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status %d, got %d", fiber.StatusCreated, resp.StatusCode)
	}
	assertReportOKResponse(t, resp)

	var report reportModel
	if err := db.Where("survey_id = ?", surveyID).First(&report).Error; err != nil {
		t.Fatalf("load report: %v", err)
	}
	if report.ReporterUserID != nil {
		t.Fatalf("expected no reporter_user_id, got %#v", report.ReporterUserID)
	}
	if report.ReporterGuestID == nil || *report.ReporterGuestID != guestID {
		t.Fatalf("expected reporter_guest_id %q, got %#v", guestID, report.ReporterGuestID)
	}
}

func TestReportEndpointMissingReasonReturnsBadRequest(t *testing.T) {
	db := newReportEndpointTestDB(t)
	app := newReportEndpointTestApp(db)
	_, creatorID := registerReportEndpointUser(t, app, "report-missing-reason@example.com")
	surveyID := createReportEndpointSurvey(t, db, creatorID)

	resp := authJSONRequest(t, app, http.MethodPost, "/api/v1/surveys/"+surveyID+"/report", `{}`, nil)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

func TestReportEndpointSurveyNotFound(t *testing.T) {
	db := newReportEndpointTestDB(t)
	app := newReportEndpointTestApp(db)

	resp := authJSONRequest(t, app, http.MethodPost, "/api/v1/surveys/"+uuid.NewString()+"/report", `{
		"reason":"spam"
	}`, nil)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected status %d, got %d", fiber.StatusNotFound, resp.StatusCode)
	}
}

func newReportEndpointTestApp(db *gorm.DB) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	RegisterAuthRoutes(app, db, config.AuthConfig{
		JWTSecret:           "report-endpoint-test-secret",
		AccessTokenTTLMin:   15,
		RefreshTokenTTLHour: 24,
		RefreshCookieName:   "refresh_token",
		RefreshCookieSecure: false,
	})
	RegisterResultsReportRoutes(app, db, "report-endpoint-test-secret")
	return app
}

func newReportEndpointTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db := newAuthTestDB(t)
	if err := db.AutoMigrate(&surveyModel{}, &reportModel{}); err != nil {
		t.Fatalf("migrate report endpoint test tables: %v", err)
	}
	return db
}

func registerReportEndpointUser(t *testing.T, app *fiber.App, email string) (string, string) {
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

func createReportEndpointSurvey(t *testing.T, db *gorm.DB, creatorID string) string {
	t.Helper()

	now := time.Now().UTC()
	record := surveyModel{
		ID:                  uuid.NewString(),
		CreatorID:           creatorID,
		Title:               "Report endpoint survey",
		Visibility:          "public",
		ResultsMode:         "open_live",
		MaxVotesPerUser:     1,
		AllowVoteChangeOnce: false,
		VoteEndsAt:          now.Add(time.Hour),
		ResultsEndsAt:       now.Add(2 * time.Hour),
		RetentionEndsAt:     now.Add(2 * time.Hour),
		ModerationStatus:    "approved",
	}
	if err := db.Create(&record).Error; err != nil {
		t.Fatalf("create report endpoint survey: %v", err)
	}
	return record.ID
}

func assertReportOKResponse(t *testing.T, resp *http.Response) {
	t.Helper()

	var body map[string]bool
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode ok response: %v", err)
	}
	if !body["ok"] {
		t.Fatal("expected ok to be true")
	}
}
