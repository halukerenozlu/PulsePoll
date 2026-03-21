package routes

import (
	"errors"
	"strings"
	"time"

	surveydomain "PulsePoll/internal/domain/survey"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type resultsReportHandler struct {
	db        *gorm.DB
	jwtSecret string
}

type surveyResultsModel struct {
	ID            string     `gorm:"column:id"`
	ResultsMode   string     `gorm:"column:results_mode"`
	Visibility    string     `gorm:"column:visibility"`
	VoteEndsAt    time.Time  `gorm:"column:vote_ends_at"`
	ResultsEndsAt time.Time  `gorm:"column:results_ends_at"`
	Options       []optionDB `gorm:"foreignKey:SurveyID;references:ID"`
}

func (surveyResultsModel) TableName() string { return "surveys" }

type surveyIDModel struct {
	ID string `gorm:"column:id"`
}

func (surveyIDModel) TableName() string { return "surveys" }

type reportModel struct {
	ID              string    `gorm:"column:id"`
	SurveyID        string    `gorm:"column:survey_id"`
	ReporterUserID  *string   `gorm:"column:reporter_user_id"`
	ReporterGuestID *string   `gorm:"column:reporter_guest_id"`
	Reason          string    `gorm:"column:reason"`
	Details         *string   `gorm:"column:details"`
	Status          string    `gorm:"column:status"`
	CreatedAt       time.Time `gorm:"column:created_at"`
}

func (reportModel) TableName() string { return "reports" }

type surveyResultOptionResponse struct {
	ID         string  `json:"id"`
	Text       string  `json:"text"`
	VoteCount  int64   `json:"vote_count"`
	Percentage float64 `json:"percentage"`
}

type surveyResultsResponse struct {
	SurveyID   string                       `json:"survey_id"`
	TotalVotes int64                        `json:"total_votes"`
	Options    []surveyResultOptionResponse `json:"options"`
}

type reportRequest struct {
	Reason  string `json:"reason"`
	Details string `json:"details"`
}

func RegisterResultsReportRoutes(app *fiber.App, db *gorm.DB, jwtSecret string) {
	h := &resultsReportHandler{
		db:        db,
		jwtSecret: jwtSecret,
	}

	api := app.Group("/api/v1")
	api.Get("/surveys/:id/results", h.getResults)
	api.Post("/surveys/:id/report", h.createReport)
}

func (h *resultsReportHandler) getResults(c *fiber.Ctx) error {
	surveyID, err := requireUUIDPathParam(c, "id")
	if err != nil {
		return nil
	}

	var survey surveyResultsModel
	err = h.db.
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Where("id = ?", surveyID).
		First(&survey).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return writeError(c, fiber.StatusNotFound, "NOT_FOUND", "survey not found")
		}
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to load survey")
	}

	flags := surveydomain.ComputeFlags(
		time.Now().UTC(),
		surveydomain.Visibility(survey.Visibility),
		surveydomain.ResultsMode(survey.ResultsMode),
		survey.VoteEndsAt,
		survey.ResultsEndsAt,
	)
	if !flags.ResultsVisible {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "results not visible")
	}

	var total int64
	for _, option := range survey.Options {
		total += option.VoteCount
	}

	respOptions := make([]surveyResultOptionResponse, 0, len(survey.Options))
	for _, option := range survey.Options {
		percentage := 0.0
		if total > 0 {
			percentage = (float64(option.VoteCount) / float64(total)) * 100
		}
		respOptions = append(respOptions, surveyResultOptionResponse{
			ID:         option.ID,
			Text:       option.Text,
			VoteCount:  option.VoteCount,
			Percentage: percentage,
		})
	}

	return c.JSON(surveyResultsResponse{
		SurveyID:   survey.ID,
		TotalVotes: total,
		Options:    respOptions,
	})
}

func (h *resultsReportHandler) createReport(c *fiber.Ctx) error {
	surveyID, err := requireUUIDPathParam(c, "id")
	if err != nil {
		return nil
	}

	var survey surveyIDModel
	if err := h.db.Where("id = ?", surveyID).First(&survey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return writeError(c, fiber.StatusNotFound, "NOT_FOUND", "survey not found")
		}
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to load survey")
	}

	var req reportRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid request body")
	}

	req.Reason = strings.TrimSpace(req.Reason)
	req.Details = strings.TrimSpace(req.Details)

	if req.Reason == "" {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "reason is required")
	}

	var reporterUserID *string
	authHeader := strings.TrimSpace(c.Get("Authorization"))
	if authHeader != "" {
		userID, err := h.userIDFromAuthorization(authHeader)
		if err != nil {
			return writeError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "invalid access token")
		}
		reporterUserID = &userID
	}

	var reporterGuestID *string
	if reporterUserID == nil {
		guestID := strings.TrimSpace(c.Cookies(guestIDCookieName))
		if guestID != "" {
			reporterGuestID = &guestID
		}
	}

	report := reportModel{
		SurveyID:        survey.ID,
		ReporterUserID:  reporterUserID,
		ReporterGuestID: reporterGuestID,
		Reason:          req.Reason,
		Status:          "open",
	}
	if req.Details != "" {
		report.Details = &req.Details
	}

	if err := h.db.Create(&report).Error; err != nil {
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to create report")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"ok": true})
}

func (h *resultsReportHandler) userIDFromAuthorization(authHeader string) (string, error) {
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
