package routes

import (
	"PulsePoll/internal/config"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestCreateSurveyEndpointSuccessAndPersistence(t *testing.T) {
	db := newSurveyEndpointTestDB(t)
	app := newSurveyEndpointTestApp(db)
	token := registerSurveyEndpointUser(t, app, "survey-create@example.com")

	resp := authJSONRequestWithBearer(t, app, http.MethodPost, "/api/v1/surveys", `{
		"title":"Best backend check?",
		"options":["Go tests","Manual curl"],
		"visibility":"public",
		"results_mode":"open_live",
		"max_votes_per_user":1
	}`, nil, token)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status %d, got %d", fiber.StatusCreated, resp.StatusCode)
	}

	var created surveyDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected created survey id")
	}
	if len(created.Options) != 2 {
		t.Fatalf("expected 2 options, got %d", len(created.Options))
	}
	if created.Options[0].ID == "" || created.Options[1].ID == "" {
		t.Fatal("expected option ids in create response")
	}

	fetchedResp := authJSONRequest(t, app, http.MethodGet, "/api/v1/surveys/"+created.ID, "", nil)
	defer fetchedResp.Body.Close()
	if fetchedResp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected fetch status %d, got %d", fiber.StatusOK, fetchedResp.StatusCode)
	}

	var fetched surveyDetailResponse
	if err := json.NewDecoder(fetchedResp.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode fetched survey: %v", err)
	}
	if fetched.ID != created.ID {
		t.Fatalf("expected fetched id %q, got %q", created.ID, fetched.ID)
	}
	if len(fetched.Options) != 2 {
		t.Fatalf("expected fetched survey to include 2 options, got %d", len(fetched.Options))
	}
}

func TestCreateSurveyEndpointAuthFailure(t *testing.T) {
	db := newSurveyEndpointTestDB(t)
	app := newSurveyEndpointTestApp(db)

	resp := authJSONRequest(t, app, http.MethodPost, "/api/v1/surveys", `{
		"title":"No token survey",
		"options":["A","B"],
		"visibility":"public",
		"results_mode":"open_live"
	}`, nil)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}
}

func TestCreateSurveyEndpointValidationFailures(t *testing.T) {
	db := newSurveyEndpointTestDB(t)
	app := newSurveyEndpointTestApp(db)
	token := registerSurveyEndpointUser(t, app, "survey-validation@example.com")

	tests := []struct {
		name string
		body string
	}{
		{
			name: "missing title",
			body: `{
				"options":["A","B"],
				"visibility":"public",
				"results_mode":"open_live"
			}`,
		},
		{
			name: "empty options",
			body: `{
				"title":"Empty options survey",
				"options":[],
				"visibility":"public",
				"results_mode":"open_live"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := authJSONRequestWithBearer(t, app, http.MethodPost, "/api/v1/surveys", tt.body, nil, token)
			defer resp.Body.Close()

			if resp.StatusCode != fiber.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
			}
		})
	}
}

func TestCreateSurveyEndpointModerationBlockedTitle(t *testing.T) {
	db := newSurveyEndpointTestDB(t)
	app := newSurveyEndpointTestApp(db)
	token := registerSurveyEndpointUser(t, app, "survey-moderation@example.com")

	resp := authJSONRequestWithBearer(t, app, http.MethodPost, "/api/v1/surveys", `{
		"title":"Obvious spam survey",
		"options":["A","B"],
		"visibility":"public",
		"results_mode":"open_live"
	}`, nil, token)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest && resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status %d or %d, got %d", fiber.StatusBadRequest, fiber.StatusForbidden, resp.StatusCode)
	}
}

func TestGetSurveyEndpointSuccessIncludesComputedFields(t *testing.T) {
	db := newSurveyEndpointTestDB(t)
	app := newSurveyEndpointTestApp(db)
	token := registerSurveyEndpointUser(t, app, "survey-detail@example.com")
	surveyID := createSurveyViaEndpoint(t, app, token, "Computed fields survey", "private_pin", "open_live")

	resp := authJSONRequest(t, app, http.MethodGet, "/api/v1/surveys/"+surveyID, "", nil)
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		t.Fatalf("decode survey response: %v", err)
	}
	for _, field := range []string{"phase", "can_vote", "results_visible", "requires_pin"} {
		if _, ok := raw[field]; !ok {
			t.Fatalf("expected computed field %q in response", field)
		}
	}

	var detail surveyDetailResponse
	mustUnmarshalSurveyField(t, raw, "phase", &detail.Phase)
	mustUnmarshalSurveyField(t, raw, "can_vote", &detail.CanVote)
	mustUnmarshalSurveyField(t, raw, "results_visible", &detail.ResultsVisible)
	mustUnmarshalSurveyField(t, raw, "requires_pin", &detail.RequiresPIN)

	if detail.Phase == "" {
		t.Fatal("expected phase value")
	}
	if !detail.CanVote {
		t.Fatal("expected can_vote to be true during voting")
	}
	if !detail.ResultsVisible {
		t.Fatal("expected open_live results to be visible during voting")
	}
	if !detail.RequiresPIN {
		t.Fatal("expected private_pin survey to require pin")
	}
}

func TestGetSurveyEndpointNotFound(t *testing.T) {
	db := newSurveyEndpointTestDB(t)
	app := newSurveyEndpointTestApp(db)

	resp := authJSONRequest(t, app, http.MethodGet, "/api/v1/surveys/"+uuid.NewString(), "", nil)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected status %d, got %d", fiber.StatusNotFound, resp.StatusCode)
	}
}

func TestResultsEndpointVisibleOpenLiveDuringVoting(t *testing.T) {
	db := newSurveyEndpointTestDB(t)
	app := newSurveyEndpointTestApp(db)
	surveyID := createSurveyRecordWithOptions(t, db, "open_live", []int64{3, 1})

	resp := authJSONRequest(t, app, http.MethodGet, "/api/v1/surveys/"+surveyID+"/results", "", nil)
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	var results surveyResultsResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		t.Fatalf("decode results response: %v", err)
	}
	if results.SurveyID != surveyID {
		t.Fatalf("expected survey_id %q, got %q", surveyID, results.SurveyID)
	}
	if results.TotalVotes != 4 {
		t.Fatalf("expected total_votes 4, got %d", results.TotalVotes)
	}
	if len(results.Options) != 2 {
		t.Fatalf("expected 2 result options, got %d", len(results.Options))
	}
	if results.Options[0].VoteCount != 3 || results.Options[0].Percentage != 75 {
		t.Fatalf("expected first option count 3 and percentage 75, got count %d percentage %v", results.Options[0].VoteCount, results.Options[0].Percentage)
	}
	if results.Options[1].VoteCount != 1 || results.Options[1].Percentage != 25 {
		t.Fatalf("expected second option count 1 and percentage 25, got count %d percentage %v", results.Options[1].VoteCount, results.Options[1].Percentage)
	}
}

func TestResultsEndpointHiddenClosedUntilEndDuringVoting(t *testing.T) {
	db := newSurveyEndpointTestDB(t)
	app := newSurveyEndpointTestApp(db)
	surveyID := createSurveyRecordWithOptions(t, db, "closed_hidden_until_end", []int64{2, 0})

	resp := authJSONRequest(t, app, http.MethodGet, "/api/v1/surveys/"+surveyID+"/results", "", nil)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusForbidden && resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d or %d, got %d", fiber.StatusForbidden, fiber.StatusOK, resp.StatusCode)
	}
	if resp.StatusCode == fiber.StatusOK {
		var results surveyResultsResponse
		if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
			t.Fatalf("decode hidden results response: %v", err)
		}
		if results.TotalVotes != 0 || len(results.Options) != 0 {
			t.Fatalf("expected empty hidden results, got total %d options %d", results.TotalVotes, len(results.Options))
		}
	}
}

func TestResultsEndpointNotFound(t *testing.T) {
	db := newSurveyEndpointTestDB(t)
	app := newSurveyEndpointTestApp(db)

	resp := authJSONRequest(t, app, http.MethodGet, "/api/v1/surveys/"+uuid.NewString()+"/results", "", nil)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected status %d, got %d", fiber.StatusNotFound, resp.StatusCode)
	}
}

func newSurveyEndpointTestApp(db *gorm.DB) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	RegisterAuthRoutes(app, db, config.AuthConfig{
		JWTSecret:           "survey-endpoint-test-secret",
		AccessTokenTTLMin:   15,
		RefreshTokenTTLHour: 24,
		RefreshCookieName:   "refresh_token",
		RefreshCookieSecure: false,
	})
	RegisterSurveyRoutes(app, db, "survey-endpoint-test-secret")
	RegisterResultsReportRoutes(app, db, "survey-endpoint-test-secret")
	return app
}

func newSurveyEndpointTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db := newAuthTestDB(t)
	if err := db.AutoMigrate(&surveyModel{}, &optionDB{}, &reportModel{}); err != nil {
		t.Fatalf("migrate survey endpoint test tables: %v", err)
	}
	return db
}

func registerSurveyEndpointUser(t *testing.T, app *fiber.App, email string) string {
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
	return body.AccessToken
}

func createSurveyViaEndpoint(t *testing.T, app *fiber.App, token string, title string, visibility string, resultsMode string) string {
	t.Helper()

	body := fmt.Sprintf(`{
		"title":%q,
		"options":["A","B"],
		"visibility":%q,
		"access_pin":"1234",
		"results_mode":%q
	}`, title, visibility, resultsMode)
	resp := authJSONRequestWithBearer(t, app, http.MethodPost, "/api/v1/surveys", body, nil, token)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("create survey: expected status %d, got %d", fiber.StatusCreated, resp.StatusCode)
	}

	var created surveyDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode created survey: %v", err)
	}
	return created.ID
}

func createSurveyRecordWithOptions(t *testing.T, db *gorm.DB, resultsMode string, voteCounts []int64) string {
	t.Helper()

	now := time.Now().UTC()
	record := surveyModel{
		ID:                  uuid.NewString(),
		CreatorID:           uuid.NewString(),
		Title:               "Results survey",
		Visibility:          "public",
		ResultsMode:         resultsMode,
		MaxVotesPerUser:     1,
		AllowVoteChangeOnce: false,
		VoteEndsAt:          now.Add(time.Hour),
		ResultsEndsAt:       now.Add(2 * time.Hour),
		RetentionEndsAt:     now.Add(2 * time.Hour),
		ModerationStatus:    "approved",
	}
	if err := db.Create(&record).Error; err != nil {
		t.Fatalf("create results survey: %v", err)
	}

	for i, voteCount := range voteCounts {
		option := optionDB{
			ID:        uuid.NewString(),
			SurveyID:  record.ID,
			Text:      fmt.Sprintf("Option %d", i+1),
			Position:  i + 1,
			VoteCount: voteCount,
		}
		if err := db.Create(&option).Error; err != nil {
			t.Fatalf("create result option: %v", err)
		}
	}
	return record.ID
}

func mustUnmarshalSurveyField(t *testing.T, raw map[string]json.RawMessage, field string, target any) {
	t.Helper()

	if err := json.Unmarshal(raw[field], target); err != nil {
		t.Fatalf("decode field %q: %v", field, err)
	}
}
