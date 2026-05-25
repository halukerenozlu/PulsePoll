package routes

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestFeedEndpointSuccessReturnsPublicSurveysWithShape(t *testing.T) {
	db := newFeedEndpointTestDB(t)
	app := newFeedEndpointTestApp(db)
	creatorID := createFeedEndpointUser(t, db, "feed-success@example.com")
	createFeedEndpointSurvey(t, db, feedSurveyConfig{
		CreatorID:  creatorID,
		Title:      "Public feed survey",
		Visibility: "public",
		CreatedAt:  time.Now().UTC(),
	})

	resp := authJSONRequest(t, app, http.MethodGet, "/api/v1/feed", "", nil)
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		t.Fatalf("decode feed response: %v", err)
	}
	itemsRaw, ok := raw["items"]
	if !ok {
		t.Fatal("expected items field in feed response")
	}

	var items []map[string]json.RawMessage
	if err := json.Unmarshal(itemsRaw, &items); err != nil {
		t.Fatalf("decode feed items: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 feed item, got %d", len(items))
	}
	for _, field := range []string{"id", "title", "visibility", "results_mode", "created_at", "vote_ends_at", "results_ends_at", "phase", "can_vote", "results_visible", "requires_pin"} {
		if _, ok := items[0][field]; !ok {
			t.Fatalf("expected feed item field %q", field)
		}
	}
}

func TestFeedEndpointSearchReturnsOnlyMatchingSurveys(t *testing.T) {
	db := newFeedEndpointTestDB(t)
	app := newFeedEndpointTestApp(db)
	creatorID := createFeedEndpointUser(t, db, "feed-search@example.com")
	createFeedEndpointSurvey(t, db, feedSurveyConfig{
		CreatorID:   creatorID,
		Title:       "Weekend coffee poll",
		Description: "Find the best beans",
		Visibility:  "public",
		CreatedAt:   time.Now().UTC().Add(-time.Hour),
	})
	createFeedEndpointSurvey(t, db, feedSurveyConfig{
		CreatorID:  creatorID,
		Title:      "Lunch location poll",
		Visibility: "public",
		CreatedAt:  time.Now().UTC(),
	})

	resp := authJSONRequest(t, app, http.MethodGet, "/api/v1/feed?search=coffee", "", nil)
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	items := decodeFeedItems(t, resp)
	if len(items) != 1 {
		t.Fatalf("expected 1 matching item, got %d", len(items))
	}
	if items[0].Title != "Weekend coffee poll" {
		t.Fatalf("expected coffee survey, got %q", items[0].Title)
	}
}

func TestFeedEndpointSortNewOrdersByCreatedAtDescending(t *testing.T) {
	db := newFeedEndpointTestDB(t)
	app := newFeedEndpointTestApp(db)
	creatorID := createFeedEndpointUser(t, db, "feed-sort@example.com")
	now := time.Now().UTC()
	createFeedEndpointSurvey(t, db, feedSurveyConfig{
		CreatorID:  creatorID,
		Title:      "Old survey",
		Visibility: "public",
		CreatedAt:  now.Add(-2 * time.Hour),
	})
	createFeedEndpointSurvey(t, db, feedSurveyConfig{
		CreatorID:  creatorID,
		Title:      "New survey",
		Visibility: "public",
		CreatedAt:  now,
	})
	createFeedEndpointSurvey(t, db, feedSurveyConfig{
		CreatorID:  creatorID,
		Title:      "Middle survey",
		Visibility: "public",
		CreatedAt:  now.Add(-time.Hour),
	})

	resp := authJSONRequest(t, app, http.MethodGet, "/api/v1/feed?sort=new", "", nil)
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	items := decodeFeedItems(t, resp)
	if len(items) != 3 {
		t.Fatalf("expected 3 feed items, got %d", len(items))
	}
	got := []string{items[0].Title, items[1].Title, items[2].Title}
	want := []string{"New survey", "Middle survey", "Old survey"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected order %v, got %v", want, got)
		}
	}
}

func TestFeedEndpointExcludesNonPublicSurveys(t *testing.T) {
	db := newFeedEndpointTestDB(t)
	app := newFeedEndpointTestApp(db)
	creatorID := createFeedEndpointUser(t, db, "feed-visibility@example.com")
	createFeedEndpointSurvey(t, db, feedSurveyConfig{
		CreatorID:  creatorID,
		Title:      "Visible public survey",
		Visibility: "public",
		CreatedAt:  time.Now().UTC(),
	})
	createFeedEndpointSurvey(t, db, feedSurveyConfig{
		CreatorID:  creatorID,
		Title:      "Hidden unlisted survey",
		Visibility: "unlisted",
		CreatedAt:  time.Now().UTC().Add(time.Minute),
	})
	createFeedEndpointSurvey(t, db, feedSurveyConfig{
		CreatorID:  creatorID,
		Title:      "Hidden private pin survey",
		Visibility: "private_pin",
		AccessPIN:  "1234",
		CreatedAt:  time.Now().UTC().Add(2 * time.Minute),
	})

	resp := authJSONRequest(t, app, http.MethodGet, "/api/v1/feed", "", nil)
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	items := decodeFeedItems(t, resp)
	if len(items) != 1 {
		t.Fatalf("expected 1 public item, got %d", len(items))
	}
	if items[0].Title != "Visible public survey" {
		t.Fatalf("expected public survey, got %q", items[0].Title)
	}
}

func TestFeedEndpointEmptyResultReturnsEmptyList(t *testing.T) {
	db := newFeedEndpointTestDB(t)
	app := newFeedEndpointTestApp(db)
	creatorID := createFeedEndpointUser(t, db, "feed-empty@example.com")
	createFeedEndpointSurvey(t, db, feedSurveyConfig{
		CreatorID:  creatorID,
		Title:      "Only unlisted survey",
		Visibility: "unlisted",
		CreatedAt:  time.Now().UTC(),
	})

	resp := authJSONRequest(t, app, http.MethodGet, "/api/v1/feed", "", nil)
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	items := decodeFeedItems(t, resp)
	if len(items) != 0 {
		t.Fatalf("expected empty feed items, got %d", len(items))
	}
}

type feedSurveyConfig struct {
	CreatorID   string
	Title       string
	Description string
	Visibility  string
	AccessPIN   string
	CreatedAt   time.Time
}

type feedEndpointResponse struct {
	Items []feedItemResponse `json:"items"`
}

func newFeedEndpointTestApp(db *gorm.DB) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	RegisterSurveyRoutes(app, db, "feed-endpoint-test-secret")
	return app
}

func newFeedEndpointTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db := newAuthTestDB(t)
	if err := db.AutoMigrate(&surveyModel{}, &optionDB{}); err != nil {
		t.Fatalf("migrate feed endpoint test tables: %v", err)
	}
	return db
}

func createFeedEndpointUser(t *testing.T, db *gorm.DB, email string) string {
	t.Helper()

	record := user{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: "test-hash",
		DisplayName:  "Feed Test User",
		Status:       "active",
	}
	if err := db.Create(&record).Error; err != nil {
		t.Fatalf("create feed endpoint user: %v", err)
	}
	return record.ID
}

func createFeedEndpointSurvey(t *testing.T, db *gorm.DB, cfg feedSurveyConfig) string {
	t.Helper()

	now := time.Now().UTC()
	createdAt := cfg.CreatedAt
	if createdAt.IsZero() {
		createdAt = now
	}
	visibility := cfg.Visibility
	if visibility == "" {
		visibility = "public"
	}

	record := surveyModel{
		ID:                  uuid.NewString(),
		CreatorID:           cfg.CreatorID,
		Title:               cfg.Title,
		Visibility:          visibility,
		ResultsMode:         "open_live",
		MaxVotesPerUser:     1,
		AllowVoteChangeOnce: false,
		CreatedAt:           createdAt,
		UpdatedAt:           createdAt,
		VoteEndsAt:          now.Add(time.Hour),
		ResultsEndsAt:       now.Add(2 * time.Hour),
		RetentionEndsAt:     now.Add(2 * time.Hour),
		ModerationStatus:    "approved",
	}
	if cfg.Description != "" {
		record.Description = &cfg.Description
	}
	if cfg.AccessPIN != "" {
		hash := "test-pin-hash"
		record.AccessPinHash = &hash
	}
	if err := db.Create(&record).Error; err != nil {
		t.Fatalf("create feed endpoint survey: %v", err)
	}
	return record.ID
}

func decodeFeedItems(t *testing.T, resp *http.Response) []feedItemResponse {
	t.Helper()

	var body feedEndpointResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode feed response: %v", err)
	}
	return body.Items
}
