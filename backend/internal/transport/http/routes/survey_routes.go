package routes

import (
	"errors"
	"strings"
	"time"

	surveydomain "PulsePoll/internal/domain/survey"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type surveyHandler struct {
	db        *gorm.DB
	jwtSecret string
}

type surveyModel struct {
	ID                  string     `gorm:"column:id"`
	CreatorID           string     `gorm:"column:creator_id"`
	Title               string     `gorm:"column:title"`
	Description         *string    `gorm:"column:description"`
	Visibility          string     `gorm:"column:visibility"`
	AccessPinHash       *string    `gorm:"column:access_pin_hash"`
	ResultsMode         string     `gorm:"column:results_mode"`
	MaxVotesPerUser     int        `gorm:"column:max_votes_per_user"`
	AllowVoteChangeOnce bool       `gorm:"column:allow_vote_change_once"`
	CreatedAt           time.Time  `gorm:"column:created_at"`
	UpdatedAt           time.Time  `gorm:"column:updated_at"`
	VoteEndsAt          time.Time  `gorm:"column:vote_ends_at"`
	ResultsEndsAt       time.Time  `gorm:"column:results_ends_at"`
	RetentionEndsAt     time.Time  `gorm:"column:retention_ends_at"`
	ModerationStatus    string     `gorm:"column:moderation_status"`
	ModerationReason    *string    `gorm:"column:moderation_reason"`
	Options             []optionDB `gorm:"foreignKey:SurveyID;references:ID"`
}

func (surveyModel) TableName() string { return "surveys" }

type optionDB struct {
	ID        string    `gorm:"column:id"`
	SurveyID  string    `gorm:"column:survey_id"`
	Text      string    `gorm:"column:text"`
	Position  int       `gorm:"column:position"`
	VoteCount int64     `gorm:"column:vote_count"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (optionDB) TableName() string { return "survey_options" }

type createSurveyRequest struct {
	Title               string     `json:"title"`
	Description         string     `json:"description"`
	Options             []string   `json:"options"`
	Visibility          string     `json:"visibility"`
	AccessPIN           string     `json:"access_pin"`
	ResultsMode         string     `json:"results_mode"`
	MaxVotesPerUser     int        `json:"max_votes_per_user"`
	AllowVoteChangeOnce bool       `json:"allow_vote_change_once"`
	VoteEndsAt          *time.Time `json:"vote_ends_at"`
	ResultsEndsAt       *time.Time `json:"results_ends_at"`
	RetentionEndsAt     *time.Time `json:"retention_ends_at"`
}

type surveyOptionResponse struct {
	ID       string `json:"id"`
	Text     string `json:"text"`
	Position int    `json:"position"`
}

type surveyDetailResponse struct {
	ID                  string                 `json:"id"`
	CreatorID           string                 `json:"creator_id"`
	Title               string                 `json:"title"`
	Description         string                 `json:"description,omitempty"`
	Visibility          string                 `json:"visibility"`
	ResultsMode         string                 `json:"results_mode"`
	MaxVotesPerUser     int                    `json:"max_votes_per_user"`
	AllowVoteChangeOnce bool                   `json:"allow_vote_change_once"`
	CreatedAt           time.Time              `json:"created_at"`
	VoteEndsAt          time.Time              `json:"vote_ends_at"`
	ResultsEndsAt       time.Time              `json:"results_ends_at"`
	RetentionEndsAt     time.Time              `json:"retention_ends_at"`
	Phase               string                 `json:"phase"`
	CanVote             bool                   `json:"can_vote"`
	ResultsVisible      bool                   `json:"results_visible"`
	RequiresPIN         bool                   `json:"requires_pin"`
	Options             []surveyOptionResponse `json:"options"`
}

type feedItemResponse struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description,omitempty"`
	Visibility     string    `json:"visibility"`
	ResultsMode    string    `json:"results_mode"`
	CreatedAt      time.Time `json:"created_at"`
	VoteEndsAt     time.Time `json:"vote_ends_at"`
	ResultsEndsAt  time.Time `json:"results_ends_at"`
	Phase          string    `json:"phase"`
	CanVote        bool      `json:"can_vote"`
	ResultsVisible bool      `json:"results_visible"`
	RequiresPIN    bool      `json:"requires_pin"`
}

var blockedTerms = []string{
	"spam",
	"scam",
	"fraud",
}

func RegisterSurveyRoutes(app *fiber.App, db *gorm.DB, jwtSecret string) {
	h := &surveyHandler{db: db, jwtSecret: jwtSecret}

	api := app.Group("/api/v1")
	api.Post("/surveys", h.createSurvey)
	api.Get("/surveys/:id", h.getSurvey)
	api.Get("/feed", h.getFeed)
}

func (h *surveyHandler) createSurvey(c *fiber.Ctx) error {
	userID, err := h.userIDFromAuthorization(c.Get("Authorization"))
	if err != nil {
		return writeError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "auth required")
	}

	var req createSurveyRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid request body")
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.Visibility = strings.TrimSpace(req.Visibility)
	req.AccessPIN = strings.TrimSpace(req.AccessPIN)
	req.ResultsMode = strings.TrimSpace(req.ResultsMode)

	if req.Title == "" {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "title is required")
	}
	if len(req.Options) == 0 {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "options are required")
	}
	if req.Visibility == "" || req.ResultsMode == "" {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "visibility and results_mode are required")
	}
	switch req.Visibility {
	case "public", "unlisted", "private_pin":
	default:
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid visibility")
	}
	switch req.ResultsMode {
	case "open_live", "closed_hidden_until_end":
	default:
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid results_mode")
	}
	if req.Visibility == "private_pin" && req.AccessPIN == "" {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "access_pin is required for private_pin")
	}

	if req.MaxVotesPerUser == 0 {
		req.MaxVotesPerUser = 1
	}
	if req.MaxVotesPerUser < 1 {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "max_votes_per_user must be >= 1")
	}
	if req.AllowVoteChangeOnce && req.MaxVotesPerUser != 1 {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "allow_vote_change_once requires max_votes_per_user == 1")
	}

	if hasBlockedKeyword(req.Title) || hasBlockedKeyword(req.Description) {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "survey contains blocked terms")
	}

	options := make([]string, 0, len(req.Options))
	for _, raw := range req.Options {
		opt := strings.TrimSpace(raw)
		if opt == "" {
			return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "option text cannot be empty")
		}
		if hasBlockedKeyword(opt) {
			return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "survey contains blocked terms")
		}
		options = append(options, opt)
	}

	now := time.Now().UTC()
	voteEndsAt := now.Add(24 * time.Hour)
	resultsEndsAt := now.Add(48 * time.Hour)
	retentionEndsAt := now.Add(48 * time.Hour)

	if req.VoteEndsAt != nil {
		voteEndsAt = req.VoteEndsAt.UTC()
	}
	if req.ResultsEndsAt != nil {
		resultsEndsAt = req.ResultsEndsAt.UTC()
	}
	if req.RetentionEndsAt != nil {
		retentionEndsAt = req.RetentionEndsAt.UTC()
	}
	if resultsEndsAt.Before(voteEndsAt) {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "results_ends_at must be >= vote_ends_at")
	}
	if retentionEndsAt.Before(resultsEndsAt) {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "retention_ends_at must be >= results_ends_at")
	}

	record := surveyModel{
		CreatorID:           userID,
		Title:               req.Title,
		Visibility:          req.Visibility,
		ResultsMode:         req.ResultsMode,
		MaxVotesPerUser:     req.MaxVotesPerUser,
		AllowVoteChangeOnce: req.AllowVoteChangeOnce,
		VoteEndsAt:          voteEndsAt,
		ResultsEndsAt:       resultsEndsAt,
		RetentionEndsAt:     retentionEndsAt,
		ModerationStatus:    "approved",
	}
	if req.Description != "" {
		record.Description = &req.Description
	}
	if req.Visibility == "private_pin" {
		accessPINHash, err := bcrypt.GenerateFromPassword([]byte(req.AccessPIN), bcrypt.DefaultCost)
		if err != nil {
			return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to hash access pin")
		}
		hash := string(accessPINHash)
		record.AccessPinHash = &hash
	}

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}

		optionRows := make([]optionDB, 0, len(options))
		for idx, opt := range options {
			optionRows = append(optionRows, optionDB{
				SurveyID: record.ID,
				Text:     opt,
				Position: idx + 1,
			})
		}

		return tx.Create(&optionRows).Error
	}); err != nil {
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to create survey")
	}

	return h.getSurveyByID(c, record.ID, fiber.StatusCreated)
}

func (h *surveyHandler) getSurvey(c *fiber.Ctx) error {
	surveyID, err := requireUUIDPathParam(c, "id")
	if err != nil {
		return nil
	}
	return h.getSurveyByID(c, surveyID, fiber.StatusOK)
}

func (h *surveyHandler) getSurveyByID(c *fiber.Ctx, surveyID string, status int) error {
	var record surveyModel
	err := h.db.
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Where("id = ?", surveyID).
		First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return writeError(c, fiber.StatusNotFound, "NOT_FOUND", "survey not found")
		}
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to load survey")
	}

	flags := surveydomain.ComputeFlags(
		time.Now().UTC(),
		surveydomain.Visibility(record.Visibility),
		surveydomain.ResultsMode(record.ResultsMode),
		record.VoteEndsAt,
		record.ResultsEndsAt,
	)

	respOptions := make([]surveyOptionResponse, 0, len(record.Options))
	for _, opt := range record.Options {
		respOptions = append(respOptions, surveyOptionResponse{
			ID:       opt.ID,
			Text:     opt.Text,
			Position: opt.Position,
		})
	}

	resp := surveyDetailResponse{
		ID:                  record.ID,
		CreatorID:           record.CreatorID,
		Title:               record.Title,
		Visibility:          record.Visibility,
		ResultsMode:         record.ResultsMode,
		MaxVotesPerUser:     record.MaxVotesPerUser,
		AllowVoteChangeOnce: record.AllowVoteChangeOnce,
		CreatedAt:           record.CreatedAt,
		VoteEndsAt:          record.VoteEndsAt,
		ResultsEndsAt:       record.ResultsEndsAt,
		RetentionEndsAt:     record.RetentionEndsAt,
		Phase:               string(flags.Phase),
		CanVote:             flags.CanVote,
		ResultsVisible:      flags.ResultsVisible,
		RequiresPIN:         flags.RequiresPIN,
		Options:             respOptions,
	}
	if record.Description != nil {
		resp.Description = *record.Description
	}

	return c.Status(status).JSON(resp)
}

func (h *surveyHandler) getFeed(c *fiber.Ctx) error {
	sort := strings.TrimSpace(strings.ToLower(c.Query("sort")))
	if sort == "" {
		sort = "new"
	}
	if sort != "new" {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "only sort=new is supported in MVP")
	}

	visibility := strings.TrimSpace(strings.ToLower(c.Query("visibility")))
	if visibility == "" {
		visibility = "public"
	}
	if visibility != "public" {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "only visibility=public is supported in MVP")
	}

	search := strings.TrimSpace(c.Query("search"))

	query := h.db.Where("visibility = ?", "public").Order("created_at DESC")
	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(COALESCE(description, '')) LIKE ?", like, like)
	}

	var rows []surveyModel
	if err := query.Limit(50).Find(&rows).Error; err != nil {
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to load feed")
	}

	items := make([]feedItemResponse, 0, len(rows))
	now := time.Now().UTC()
	for _, row := range rows {
		flags := surveydomain.ComputeFlags(
			now,
			surveydomain.Visibility(row.Visibility),
			surveydomain.ResultsMode(row.ResultsMode),
			row.VoteEndsAt,
			row.ResultsEndsAt,
		)

		item := feedItemResponse{
			ID:             row.ID,
			Title:          row.Title,
			Visibility:     row.Visibility,
			ResultsMode:    row.ResultsMode,
			CreatedAt:      row.CreatedAt,
			VoteEndsAt:     row.VoteEndsAt,
			ResultsEndsAt:  row.ResultsEndsAt,
			Phase:          string(flags.Phase),
			CanVote:        flags.CanVote,
			ResultsVisible: flags.ResultsVisible,
			RequiresPIN:    flags.RequiresPIN,
		}
		if row.Description != nil {
			item.Description = *row.Description
		}
		items = append(items, item)
	}

	return c.JSON(fiber.Map{"items": items})
}

func (h *surveyHandler) userIDFromAuthorization(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid authorization header")
	}

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(h.jwtSecret), nil
	})
	if err != nil || !token.Valid || claims.Subject == "" {
		return "", errors.New("invalid token")
	}
	return claims.Subject, nil
}

func hasBlockedKeyword(value string) bool {
	lower := strings.ToLower(value)
	for _, blocked := range blockedTerms {
		if strings.Contains(lower, blocked) {
			return true
		}
	}
	return false
}
