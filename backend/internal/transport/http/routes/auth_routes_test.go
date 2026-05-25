package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"PulsePoll/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestAuthRegisterEndpoint(t *testing.T) {
	db := newAuthTestDB(t)
	app, _ := newAuthTestApp(db)

	resp := authJSONRequest(t, app, http.MethodPost, "/api/v1/auth/register", `{
		"email":"new-user@example.com",
		"password":"StrongPass123!",
		"display_name":"New User"
	}`, nil)
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status %d, got %d", fiber.StatusCreated, resp.StatusCode)
	}

	var body authResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode register response: %v", err)
	}
	if body.AccessToken == "" {
		t.Fatal("expected access_token in register response")
	}

	duplicate := authJSONRequest(t, app, http.MethodPost, "/api/v1/auth/register", `{
		"email":"new-user@example.com",
		"password":"StrongPass123!",
		"display_name":"New User"
	}`, nil)
	defer duplicate.Body.Close()

	if duplicate.StatusCode != fiber.StatusConflict && duplicate.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected duplicate status %d or %d, got %d", fiber.StatusConflict, fiber.StatusBadRequest, duplicate.StatusCode)
	}
}

func TestAuthLoginEndpoint(t *testing.T) {
	db := newAuthTestDB(t)
	app, _ := newAuthTestApp(db)
	registerAuthUser(t, app, "login-user@example.com", "StrongPass123!")

	success := authJSONRequest(t, app, http.MethodPost, "/api/v1/auth/login", `{
		"email":"login-user@example.com",
		"password":"StrongPass123!"
	}`, nil)
	defer success.Body.Close()
	if success.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, success.StatusCode)
	}

	wrongPassword := authJSONRequest(t, app, http.MethodPost, "/api/v1/auth/login", `{
		"email":"login-user@example.com",
		"password":"WrongPass123!"
	}`, nil)
	defer wrongPassword.Body.Close()
	if wrongPassword.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected wrong-password status %d, got %d", fiber.StatusUnauthorized, wrongPassword.StatusCode)
	}

	unknownEmail := authJSONRequest(t, app, http.MethodPost, "/api/v1/auth/login", `{
		"email":"missing-user@example.com",
		"password":"StrongPass123!"
	}`, nil)
	defer unknownEmail.Body.Close()
	if unknownEmail.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected unknown-email status %d, got %d", fiber.StatusUnauthorized, unknownEmail.StatusCode)
	}
}

func TestAuthRefreshEndpoint(t *testing.T) {
	db := newAuthTestDB(t)
	app, cfg := newAuthTestApp(db)
	registerResp := registerAuthUser(t, app, "refresh-user@example.com", "StrongPass123!")
	refreshCookie := requireCookie(t, registerResp, cfg.RefreshCookieName)

	success := authJSONRequest(t, app, http.MethodPost, "/api/v1/auth/refresh", `{}`, []*http.Cookie{refreshCookie})
	defer success.Body.Close()
	if success.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, success.StatusCode)
	}

	var body authResponse
	if err := json.NewDecoder(success.Body).Decode(&body); err != nil {
		t.Fatalf("decode refresh response: %v", err)
	}
	if body.AccessToken == "" {
		t.Fatal("expected access_token in refresh response")
	}

	missingCookie := authJSONRequest(t, app, http.MethodPost, "/api/v1/auth/refresh", `{}`, nil)
	defer missingCookie.Body.Close()
	if missingCookie.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected missing-cookie status %d, got %d", fiber.StatusUnauthorized, missingCookie.StatusCode)
	}

	invalidCookie := &http.Cookie{Name: cfg.RefreshCookieName, Value: "not-a-real-refresh-token"}
	invalid := authJSONRequest(t, app, http.MethodPost, "/api/v1/auth/refresh", `{}`, []*http.Cookie{invalidCookie})
	defer invalid.Body.Close()
	if invalid.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected invalid-cookie status %d, got %d", fiber.StatusUnauthorized, invalid.StatusCode)
	}
}

func TestAuthLogoutEndpointRevokesSession(t *testing.T) {
	db := newAuthTestDB(t)
	app, cfg := newAuthTestApp(db)
	registerResp := registerAuthUser(t, app, "logout-user@example.com", "StrongPass123!")
	refreshCookie := requireCookie(t, registerResp, cfg.RefreshCookieName)

	logout := authJSONRequest(t, app, http.MethodPost, "/api/v1/auth/logout", `{}`, []*http.Cookie{refreshCookie})
	defer logout.Body.Close()
	if logout.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, logout.StatusCode)
	}

	var session authSession
	if err := db.Where("refresh_token_hash = ?", hashToken(refreshCookie.Value)).First(&session).Error; err != nil {
		t.Fatalf("load auth session: %v", err)
	}
	if session.RevokedAt == nil {
		t.Fatal("expected revoked_at to be set after logout")
	}
}

func TestMeEndpoint(t *testing.T) {
	db := newAuthTestDB(t)
	app, _ := newAuthTestApp(db)
	registerResp := registerAuthUser(t, app, "me-user@example.com", "StrongPass123!")

	var registerBody authResponse
	if err := json.NewDecoder(registerResp.Body).Decode(&registerBody); err != nil {
		t.Fatalf("decode register response: %v", err)
	}

	success := authJSONRequestWithBearer(t, app, http.MethodGet, "/api/v1/me", "", nil, registerBody.AccessToken)
	defer success.Body.Close()
	if success.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, success.StatusCode)
	}

	missingToken := authJSONRequest(t, app, http.MethodGet, "/api/v1/me", "", nil)
	defer missingToken.Body.Close()
	if missingToken.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected missing-token status %d, got %d", fiber.StatusUnauthorized, missingToken.StatusCode)
	}

	invalidToken := authJSONRequestWithBearer(t, app, http.MethodGet, "/api/v1/me", "", nil, "invalid")
	defer invalidToken.Body.Close()
	if invalidToken.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected invalid-token status %d, got %d", fiber.StatusUnauthorized, invalidToken.StatusCode)
	}
}

func newAuthTestApp(db *gorm.DB) (*fiber.App, config.AuthConfig) {
	cfg := config.AuthConfig{
		JWTSecret:           "auth-test-secret",
		AccessTokenTTLMin:   15,
		RefreshTokenTTLHour: 24,
		RefreshCookieName:   "refresh_token",
		RefreshCookieSecure: false,
	}
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	RegisterAuthRoutes(app, db, cfg)
	return app, cfg
}

func newAuthTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("PULSEPOLL_TEST_DATABASE_DSN")
	if dsn == "" {
		dsn = "host=localhost user=ephemeral password=ephemeral dbname=ephemeral port=5432 sslmode=disable"
	}

	base, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("auth endpoint integration tests require Postgres: %v", err)
	}
	sqlDB, err := base.DB()
	if err != nil {
		t.Fatalf("get sql db: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		t.Skipf("auth endpoint integration tests require Postgres: %v", err)
	}

	schema := "auth_test_" + strings.ReplaceAll(uuid.NewString(), "-", "_")
	if err := base.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto").Error; err != nil {
		t.Fatalf("create pgcrypto extension: %v", err)
	}
	if err := base.Exec(fmt.Sprintf("CREATE SCHEMA %s", schema)).Error; err != nil {
		t.Fatalf("create test schema: %v", err)
	}
	t.Cleanup(func() {
		_ = base.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schema)).Error
		_ = sqlDB.Close()
	})

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open schema db: %v", err)
	}
	schemaSQLDB, err := db.DB()
	if err != nil {
		t.Fatalf("get schema sql db: %v", err)
	}
	schemaSQLDB.SetMaxOpenConns(1)
	t.Cleanup(func() {
		_ = schemaSQLDB.Close()
	})
	if err := db.Exec(fmt.Sprintf("SET search_path TO %s", schema)).Error; err != nil {
		t.Fatalf("set search path: %v", err)
	}
	if err := db.AutoMigrate(&user{}, &authSession{}); err != nil {
		t.Fatalf("migrate auth tables: %v", err)
	}
	if err := db.Exec("CREATE UNIQUE INDEX users_email_key ON users (email)").Error; err != nil {
		t.Fatalf("create users email unique index: %v", err)
	}
	if err := db.Exec("CREATE UNIQUE INDEX auth_sessions_refresh_token_hash_key ON auth_sessions (refresh_token_hash)").Error; err != nil {
		t.Fatalf("create auth session refresh token unique index: %v", err)
	}

	return db
}

func registerAuthUser(t *testing.T, app *fiber.App, email string, password string) *http.Response {
	t.Helper()

	resp := authJSONRequest(t, app, http.MethodPost, "/api/v1/auth/register", fmt.Sprintf(`{
		"email":%q,
		"password":%q,
		"display_name":"Test User"
	}`, email, password), nil)
	if resp.StatusCode != fiber.StatusCreated {
		defer resp.Body.Close()
		t.Fatalf("register test user: expected status %d, got %d", fiber.StatusCreated, resp.StatusCode)
	}
	return resp
}

func authJSONRequest(
	t *testing.T,
	app *fiber.App,
	method string,
	path string,
	body string,
	cookies []*http.Cookie,
) *http.Response {
	t.Helper()

	return authJSONRequestWithBearer(t, app, method, path, body, cookies, "")
}

func authJSONRequestWithBearer(
	t *testing.T,
	app *fiber.App,
	method string,
	path string,
	body string,
	cookies []*http.Cookie,
	bearerToken string,
) *http.Response {
	t.Helper()

	var reader *bytes.Reader
	if body == "" {
		reader = bytes.NewReader(nil)
	} else {
		reader = bytes.NewReader([]byte(body))
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "auth-route-test")
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}
	return resp
}

func requireCookie(t *testing.T, resp *http.Response, name string) *http.Cookie {
	t.Helper()
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		if cookie.Name == name && cookie.Value != "" && cookie.Expires.After(time.Now().Add(-time.Minute)) {
			return cookie
		}
	}
	t.Fatalf("expected response cookie %q", name)
	return nil
}
