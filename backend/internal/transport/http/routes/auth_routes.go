package routes

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"PulsePoll/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authHandler struct {
	db  *gorm.DB
	cfg config.AuthConfig
}

type user struct {
	ID           string    `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Email        string    `gorm:"column:email"`
	PasswordHash string    `gorm:"column:password_hash"`
	DisplayName  string    `gorm:"column:display_name"`
	Status       string    `gorm:"column:status"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (user) TableName() string { return "users" }

type authSession struct {
	ID               string     `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	UserID           string     `gorm:"column:user_id"`
	RefreshTokenHash string     `gorm:"column:refresh_token_hash"`
	UserAgent        string     `gorm:"column:user_agent"`
	IP               string     `gorm:"column:ip"`
	ExpiresAt        time.Time  `gorm:"column:expires_at"`
	RevokedAt        *time.Time `gorm:"column:revoked_at"`
	CreatedAt        time.Time  `gorm:"column:created_at"`
}

func (authSession) TableName() string { return "auth_sessions" }

type authRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type authResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

type userResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func RegisterAuthRoutes(app *fiber.App, db *gorm.DB, cfg config.AuthConfig) {
	h := &authHandler{db: db, cfg: cfg}

	api := app.Group("/api/v1")
	auth := api.Group("/auth")
	auth.Post("/register", h.register)
	auth.Post("/login", h.login)
	auth.Post("/refresh", h.refresh)
	auth.Post("/logout", h.logout)

	api.Get("/me", h.me)
}

func (h *authHandler) register(c *fiber.Ctx) error {
	var req authRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid request body")
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.DisplayName = strings.TrimSpace(req.DisplayName)

	if req.Email == "" || req.Password == "" || req.DisplayName == "" {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "email, password and display_name are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to hash password")
	}

	newUser := user{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		DisplayName:  req.DisplayName,
		Status:       "active",
	}

	if err := h.db.Create(&newUser).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "email already exists")
		}
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to create user")
	}

	res, refreshToken, refreshExpiry, err := h.buildAuthResponse(newUser, c.Get("User-Agent"), c.IP())
	if err != nil {
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to create session")
	}

	setRefreshCookie(c, h.cfg, refreshToken, refreshExpiry)
	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *authHandler) login(c *fiber.Ctx) error {
	var req authRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid request body")
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "email and password are required")
	}

	var existing user
	if err := h.db.Where("email = ?", req.Email).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return writeError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "invalid credentials")
		}
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to load user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existing.PasswordHash), []byte(req.Password)); err != nil {
		return writeError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "invalid credentials")
	}
	if existing.Status != "active" {
		return writeError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "invalid credentials")
	}

	res, refreshToken, refreshExpiry, err := h.buildAuthResponse(existing, c.Get("User-Agent"), c.IP())
	if err != nil {
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to create session")
	}

	setRefreshCookie(c, h.cfg, refreshToken, refreshExpiry)
	return c.JSON(res)
}

func (h *authHandler) refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies(h.cfg.RefreshCookieName)
	if refreshToken == "" {
		return writeError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "missing refresh token")
	}

	now := time.Now().UTC()
	oldHash := hashToken(refreshToken)

	var session authSession
	if err := h.db.
		Where("refresh_token_hash = ? AND revoked_at IS NULL AND expires_at > ?", oldHash, now).
		First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return writeError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "invalid refresh token")
		}
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to load session")
	}

	var existing user
	if err := h.db.Where("id = ?", session.UserID).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return writeError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "user not found")
		}
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to load user")
	}

	accessToken, err := h.createAccessToken(existing.ID)
	if err != nil {
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to create access token")
	}

	newRefreshToken, err := generateToken(32)
	if err != nil {
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to rotate refresh token")
	}
	newRefreshExpiry := now.Add(time.Duration(h.cfg.RefreshTokenTTLHour) * time.Hour)
	newSession := authSession{
		UserID:           existing.ID,
		RefreshTokenHash: hashToken(newRefreshToken),
		UserAgent:        c.Get("User-Agent"),
		IP:               c.IP(),
		ExpiresAt:        newRefreshExpiry,
	}

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		revokedAt := now
		if err := tx.Model(&authSession{}).
			Where("id = ? AND revoked_at IS NULL", session.ID).
			Update("revoked_at", revokedAt).Error; err != nil {
			return err
		}
		if err := tx.Create(&newSession).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to rotate session")
	}

	setRefreshCookie(c, h.cfg, newRefreshToken, newRefreshExpiry)
	return c.JSON(authResponse{
		AccessToken: accessToken,
		User: userResponse{
			ID:          existing.ID,
			Email:       existing.Email,
			DisplayName: existing.DisplayName,
		},
	})
}

func (h *authHandler) logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies(h.cfg.RefreshCookieName)
	if refreshToken != "" {
		hash := hashToken(refreshToken)
		now := time.Now().UTC()
		_ = h.db.Model(&authSession{}).
			Where("refresh_token_hash = ? AND revoked_at IS NULL", hash).
			Update("revoked_at", now).Error
	}

	clearRefreshCookie(c, h.cfg)
	return c.JSON(fiber.Map{"ok": true})
}

func (h *authHandler) me(c *fiber.Ctx) error {
	userID, err := h.userIDFromAuthorization(c.Get("Authorization"))
	if err != nil {
		return writeError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "invalid access token")
	}

	var existing user
	if err := h.db.Where("id = ?", userID).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return writeError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "user not found")
		}
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to load user")
	}

	return c.JSON(userResponse{
		ID:          existing.ID,
		Email:       existing.Email,
		DisplayName: existing.DisplayName,
	})
}

func (h *authHandler) buildAuthResponse(
	u user,
	userAgent string,
	ip string,
) (authResponse, string, time.Time, error) {
	accessToken, err := h.createAccessToken(u.ID)
	if err != nil {
		return authResponse{}, "", time.Time{}, err
	}

	refreshToken, err := generateToken(32)
	if err != nil {
		return authResponse{}, "", time.Time{}, err
	}

	refreshExpiry := time.Now().UTC().Add(time.Duration(h.cfg.RefreshTokenTTLHour) * time.Hour)
	session := authSession{
		UserID:           u.ID,
		RefreshTokenHash: hashToken(refreshToken),
		UserAgent:        userAgent,
		IP:               ip,
		ExpiresAt:        refreshExpiry,
	}
	if err := h.db.Create(&session).Error; err != nil {
		return authResponse{}, "", time.Time{}, err
	}

	return authResponse{
		AccessToken: accessToken,
		User: userResponse{
			ID:          u.ID,
			Email:       u.Email,
			DisplayName: u.DisplayName,
		},
	}, refreshToken, refreshExpiry, nil
}

func (h *authHandler) createAccessToken(userID string) (string, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(time.Duration(h.cfg.AccessTokenTTLMin) * time.Minute)
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWTSecret))
}

func (h *authHandler) userIDFromAuthorization(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid authorization header")
	}

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(h.cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid || claims.Subject == "" {
		return "", errors.New("invalid token")
	}
	return claims.Subject, nil
}

func generateToken(byteLen int) (string, error) {
	buf := make([]byte, byteLen)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func setRefreshCookie(c *fiber.Ctx, cfg config.AuthConfig, value string, expires time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:     cfg.RefreshCookieName,
		Value:    value,
		HTTPOnly: true,
		Secure:   cfg.RefreshCookieSecure,
		SameSite: "Lax",
		Path:     "/",
		Expires:  expires,
	})
}

func clearRefreshCookie(c *fiber.Ctx, cfg config.AuthConfig) {
	c.Cookie(&fiber.Cookie{
		Name:     cfg.RefreshCookieName,
		Value:    "",
		HTTPOnly: true,
		Secure:   cfg.RefreshCookieSecure,
		SameSite: "Lax",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}

func writeError(c *fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(errorEnvelope{
		Error: errorBody{
			Code:    code,
			Message: message,
		},
	})
}
