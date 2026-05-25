package routes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	surveydomain "PulsePoll/internal/domain/survey"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	goredis "github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type voteHandler struct {
	db        *gorm.DB
	redis     *goredis.Client
	jwtSecret string
}

const (
	voteRateLimitTTLSeconds = 60
	voteRateLimitMaxHits    = 30
	pinFailTTLSeconds       = 15 * 60
	pinFailMaxAttempts      = 5
)

var voteRateLimitScript = goredis.NewScript(`
local current = redis.call("INCR", KEYS[1])
if current == 1 then
  redis.call("EXPIRE", KEYS[1], ARGV[1])
end
return current
`)

var pinFailScript = goredis.NewScript(`
local current = redis.call("INCR", KEYS[1])
redis.call("EXPIRE", KEYS[1], ARGV[1])
return current
`)

type voteRequest struct {
	OptionID string `json:"option_id"`
	PIN      string `json:"pin"`
}

type voteChangeRequest struct {
	NewOptionID string `json:"new_option_id"`
	PIN         string `json:"pin"`
}

type pinVerifyRequest struct {
	PIN string `json:"pin"`
}

type surveyVoteModel struct {
	ID                  string    `gorm:"column:id"`
	Visibility          string    `gorm:"column:visibility"`
	AccessPinHash       *string   `gorm:"column:access_pin_hash"`
	MaxVotesPerUser     int       `gorm:"column:max_votes_per_user"`
	AllowVoteChangeOnce bool      `gorm:"column:allow_vote_change_once"`
	VoteEndsAt          time.Time `gorm:"column:vote_ends_at"`
	ResultsEndsAt       time.Time `gorm:"column:results_ends_at"`
	RetentionEndsAt     time.Time `gorm:"column:retention_ends_at"`
}

func (surveyVoteModel) TableName() string { return "surveys" }

type voteReceipt struct {
	VotesUsed    int    `json:"votes_used"`
	LastOptionID string `json:"last_option_id"`
	ChangeUsed   bool   `json:"change_used"`
}

type voterIdentity struct {
	UserID  string
	GuestID string
	IsGuest bool
}

func RegisterVoteRoutes(app *fiber.App, db *gorm.DB, redisClient *goredis.Client, jwtSecret string) {
	h := &voteHandler{
		db:        db,
		redis:     redisClient,
		jwtSecret: jwtSecret,
	}

	api := app.Group("/api/v1")
	api.Post("/surveys/:id/pin/verify", h.verifyPIN)
	api.Post("/surveys/:id/vote", h.vote)
	api.Put("/surveys/:id/vote", h.changeVote)
}

func (h *voteHandler) verifyPIN(c *fiber.Ctx) error {
	surveyID, err := requireUUIDPathParam(c, "id")
	if err != nil {
		if isResponseSent(err) {
			return nil
		}
		return err
	}

	survey, err := h.getSurveyForVote(surveyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return writeError(c, fiber.StatusNotFound, "NOT_FOUND", "survey not found")
		}
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to load survey")
	}

	if survey.Visibility != "private_pin" {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "pin is not required for this survey")
	}
	if survey.AccessPinHash == nil || *survey.AccessPinHash == "" {
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "survey pin is not configured")
	}

	identity, err := h.resolveIdentity(c)
	if err != nil {
		if isResponseSent(err) {
			return nil
		}
		return err
	}

	ctx := context.Background()
	if identity.IsGuest {
		count, err := h.guestPINFailureCount(ctx, survey.ID, identity)
		if err != nil {
			return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to read pin failure count")
		}
		if count >= pinFailMaxAttempts {
			return writeError(c, fiber.StatusTooManyRequests, "TOO_MANY_REQUESTS", "too many pin attempts")
		}
	}

	var req pinVerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid request body")
	}
	req.PIN = strings.TrimSpace(req.PIN)
	if req.PIN == "" {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "pin is required")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*survey.AccessPinHash), []byte(req.PIN)); err != nil {
		if identity.IsGuest {
			count, err := h.recordGuestPINFailure(context.Background(), survey.ID, identity)
			if err != nil {
				return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to record pin failure")
			}
			if count > pinFailMaxAttempts {
				return writeError(c, fiber.StatusTooManyRequests, "TOO_MANY_REQUESTS", "too many pin attempts")
			}
		}
		return writeError(c, fiber.StatusForbidden, "PIN_REQUIRED", "invalid pin")
	}

	ttl := pinOKTTL(time.Now().UTC(), survey.VoteEndsAt)
	if ttl <= 0 {
		return writeError(c, fiber.StatusForbidden, "PHASE_NOT_VOTING", "voting not allowed in this phase")
	}

	if err := h.redis.Set(ctx, pinOKKey(survey.ID, identity), "1", ttl).Err(); err != nil {
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to store pin verification")
	}

	return c.JSON(fiber.Map{"ok": true})
}

func (h *voteHandler) vote(c *fiber.Ctx) error {
	surveyID, err := requireUUIDPathParam(c, "id")
	if err != nil {
		if isResponseSent(err) {
			return nil
		}
		return err
	}
	if err := h.enforceVoteRateLimit(c); err != nil {
		return err
	}

	survey, err := h.getSurveyForVote(surveyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return writeError(c, fiber.StatusNotFound, "NOT_FOUND", "survey not found")
		}
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to load survey")
	}

	phase := surveydomain.CalculatePhase(time.Now().UTC(), survey.VoteEndsAt, survey.ResultsEndsAt)
	if phase != surveydomain.PhaseVoting {
		return writeError(c, fiber.StatusForbidden, "PHASE_NOT_VOTING", "voting not allowed in this phase")
	}

	identity, err := h.resolveIdentity(c)
	if err != nil {
		if isResponseSent(err) {
			return nil
		}
		return err
	}

	var req voteRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid request body")
	}
	req.OptionID = strings.TrimSpace(req.OptionID)
	req.PIN = strings.TrimSpace(req.PIN)
	if req.OptionID == "" {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "option_id is required")
	}
	if err := validateUUIDField(c, "option_id", req.OptionID); err != nil {
		if isResponseSent(err) {
			return nil
		}
		return err
	}

	if err := h.ensurePINVerified(c, survey, identity, req.PIN); err != nil {
		if isResponseSent(err) {
			return nil
		}
		return err
	}

	ctx := context.Background()
	receipt, err := h.loadVoteReceipt(ctx, survey, identity)
	if err != nil {
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to load vote receipt")
	}
	if receipt.VotesUsed >= survey.MaxVotesPerUser {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "max votes reached")
	}

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&optionDB{}).
			Where("id = ? AND survey_id = ?", req.OptionID, survey.ID).
			UpdateColumn("vote_count", gorm.Expr("vote_count + 1"))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errOptionNotFound
		}
		return nil
	}); err != nil {
		if errors.Is(err, errOptionNotFound) {
			return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid option_id")
		}
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to record vote")
	}

	receipt.VotesUsed++
	receipt.LastOptionID = req.OptionID
	if err := h.storeVoteReceipt(ctx, survey, identity, receipt); err != nil {
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to store vote receipt")
	}

	return c.JSON(fiber.Map{"ok": true})
}

func (h *voteHandler) changeVote(c *fiber.Ctx) error {
	surveyID, err := requireUUIDPathParam(c, "id")
	if err != nil {
		if isResponseSent(err) {
			return nil
		}
		return err
	}
	if err := h.enforceVoteRateLimit(c); err != nil {
		return err
	}

	survey, err := h.getSurveyForVote(surveyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return writeError(c, fiber.StatusNotFound, "NOT_FOUND", "survey not found")
		}
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to load survey")
	}

	phase := surveydomain.CalculatePhase(time.Now().UTC(), survey.VoteEndsAt, survey.ResultsEndsAt)
	if phase != surveydomain.PhaseVoting {
		return writeError(c, fiber.StatusForbidden, "PHASE_NOT_VOTING", "voting not allowed in this phase")
	}
	if survey.MaxVotesPerUser != 1 || !survey.AllowVoteChangeOnce {
		return writeError(c, fiber.StatusForbidden, "VOTE_CHANGE_NOT_ALLOWED", "vote change not allowed")
	}

	identity, err := h.resolveIdentity(c)
	if err != nil {
		if isResponseSent(err) {
			return nil
		}
		return err
	}

	var req voteChangeRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid request body")
	}
	req.NewOptionID = strings.TrimSpace(req.NewOptionID)
	req.PIN = strings.TrimSpace(req.PIN)
	if req.NewOptionID == "" {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "new_option_id is required")
	}
	if err := validateUUIDField(c, "new_option_id", req.NewOptionID); err != nil {
		if isResponseSent(err) {
			return nil
		}
		return err
	}

	if err := h.ensurePINVerified(c, survey, identity, req.PIN); err != nil {
		if isResponseSent(err) {
			return nil
		}
		return err
	}

	ctx := context.Background()
	receipt, err := h.loadVoteReceipt(ctx, survey, identity)
	if err != nil {
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to load vote receipt")
	}
	if receipt.VotesUsed == 0 || receipt.LastOptionID == "" || receipt.ChangeUsed {
		return writeError(c, fiber.StatusForbidden, "VOTE_CHANGE_NOT_ALLOWED", "vote change not allowed")
	}
	if receipt.LastOptionID == req.NewOptionID {
		return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "new_option_id must be different")
	}

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		prevRes := tx.Model(&optionDB{}).
			Where("id = ? AND survey_id = ?", receipt.LastOptionID, survey.ID).
			UpdateColumn("vote_count", gorm.Expr("vote_count - 1"))
		if prevRes.Error != nil {
			return prevRes.Error
		}
		if prevRes.RowsAffected == 0 {
			return errPreviousOptionNotFound
		}

		newRes := tx.Model(&optionDB{}).
			Where("id = ? AND survey_id = ?", req.NewOptionID, survey.ID).
			UpdateColumn("vote_count", gorm.Expr("vote_count + 1"))
		if newRes.Error != nil {
			return newRes.Error
		}
		if newRes.RowsAffected == 0 {
			return errNewOptionNotFound
		}

		return nil
	}); err != nil {
		if errors.Is(err, errNewOptionNotFound) {
			return writeError(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid new_option_id")
		}
		if errors.Is(err, errPreviousOptionNotFound) {
			return writeError(c, fiber.StatusForbidden, "VOTE_CHANGE_NOT_ALLOWED", "vote change not allowed")
		}
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to change vote")
	}

	receipt.LastOptionID = req.NewOptionID
	receipt.ChangeUsed = true
	if err := h.storeVoteReceipt(ctx, survey, identity, receipt); err != nil {
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to store vote receipt")
	}

	return c.JSON(fiber.Map{"ok": true})
}

func (h *voteHandler) getSurveyForVote(surveyID string) (surveyVoteModel, error) {
	var survey surveyVoteModel
	err := h.db.Where("id = ?", surveyID).First(&survey).Error
	return survey, err
}

func (h *voteHandler) resolveIdentity(c *fiber.Ctx) (voterIdentity, error) {
	userID, err := h.userIDFromAuthorization(c.Get("Authorization"))
	if err == nil {
		return voterIdentity{UserID: userID}, nil
	}

	guestID := strings.TrimSpace(c.Cookies(guestIDCookieName))
	if guestID == "" {
		if err := writeError(c, fiber.StatusForbidden, "CONSENT_REQUIRED", "guest voting requires consent cookie"); err != nil {
			return voterIdentity{}, err
		}
		return voterIdentity{}, errResponseSent
	}
	return voterIdentity{GuestID: guestID, IsGuest: true}, nil
}

func (h *voteHandler) ensurePINVerified(
	c *fiber.Ctx,
	survey surveyVoteModel,
	identity voterIdentity,
	pin string,
) error {
	if survey.Visibility != "private_pin" {
		return nil
	}

	if pin != "" {
		if survey.AccessPinHash == nil || *survey.AccessPinHash == "" {
			if err := writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "survey pin is not configured"); err != nil {
				return err
			}
			return errResponseSent
		}
		if err := bcrypt.CompareHashAndPassword([]byte(*survey.AccessPinHash), []byte(pin)); err != nil {
			if err := writeError(c, fiber.StatusForbidden, "PIN_REQUIRED", "pin verification required"); err != nil {
				return err
			}
			return errResponseSent
		}
		ttl := pinOKTTL(time.Now().UTC(), survey.VoteEndsAt)
		if ttl <= 0 {
			if err := writeError(c, fiber.StatusForbidden, "PHASE_NOT_VOTING", "voting not allowed in this phase"); err != nil {
				return err
			}
			return errResponseSent
		}
		if err := h.redis.Set(context.Background(), pinOKKey(survey.ID, identity), "1", ttl).Err(); err != nil {
			if err := writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to store pin verification"); err != nil {
				return err
			}
			return errResponseSent
		}
		return nil
	}

	ok, err := h.redis.Get(context.Background(), pinOKKey(survey.ID, identity)).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			if err := writeError(c, fiber.StatusForbidden, "PIN_REQUIRED", "pin verification required"); err != nil {
				return err
			}
			return errResponseSent
		}
		if err := writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to read pin verification"); err != nil {
			return err
		}
		return errResponseSent
	}
	if ok != "1" {
		if err := writeError(c, fiber.StatusForbidden, "PIN_REQUIRED", "pin verification required"); err != nil {
			return err
		}
		return errResponseSent
	}
	return nil
}

func (h *voteHandler) loadVoteReceipt(
	ctx context.Context,
	survey surveyVoteModel,
	identity voterIdentity,
) (voteReceipt, error) {
	key := voteReceiptKey(survey.ID, identity)
	raw, err := h.redis.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return voteReceipt{}, nil
		}
		return voteReceipt{}, err
	}

	var receipt voteReceipt
	if err := json.Unmarshal([]byte(raw), &receipt); err != nil {
		return voteReceipt{}, err
	}
	return receipt, nil
}

func (h *voteHandler) storeVoteReceipt(
	ctx context.Context,
	survey surveyVoteModel,
	identity voterIdentity,
	receipt voteReceipt,
) error {
	ttl := retentionTTL(time.Now().UTC(), survey.RetentionEndsAt)
	if ttl <= 0 {
		ttl = time.Second
	}

	raw, err := json.Marshal(receipt)
	if err != nil {
		return err
	}
	return h.redis.Set(ctx, voteReceiptKey(survey.ID, identity), raw, ttl).Err()
}

func voteReceiptKey(surveyID string, identity voterIdentity) string {
	if identity.IsGuest {
		return fmt.Sprintf("vote:survey:%s:guest:%s", surveyID, identity.GuestID)
	}
	return fmt.Sprintf("vote:survey:%s:user:%s", surveyID, identity.UserID)
}

func pinOKKey(surveyID string, identity voterIdentity) string {
	if identity.IsGuest {
		return fmt.Sprintf("pinok:survey:%s:guest:%s", surveyID, identity.GuestID)
	}
	return fmt.Sprintf("pinok:survey:%s:user:%s", surveyID, identity.UserID)
}

func pinFailKey(surveyID string, identity voterIdentity) string {
	return fmt.Sprintf("pinfail:survey:%s:guest:%s", surveyID, identity.GuestID)
}

func voteRateLimitKey(ip string) string {
	return fmt.Sprintf("rl:ip:%s:vote", ip)
}

func (h *voteHandler) recordGuestPINFailure(ctx context.Context, surveyID string, identity voterIdentity) (int64, error) {
	return pinFailScript.Run(ctx, h.redis, []string{pinFailKey(surveyID, identity)}, int64(pinFailTTLSeconds)).Int64()
}

func (h *voteHandler) guestPINFailureCount(ctx context.Context, surveyID string, identity voterIdentity) (int64, error) {
	raw, err := h.redis.Get(ctx, pinFailKey(surveyID, identity)).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	count, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (h *voteHandler) enforceVoteRateLimit(c *fiber.Ctx) error {
	return enforceVoteRateLimitWithFuncs(
		c,
		voteRateLimitKey(c.IP()),
		voteRateLimitMaxHits,
		func(ctx context.Context, key string) (int64, error) {
			ttlSeconds := int64(voteRateLimitTTLSeconds)
			return voteRateLimitScript.Run(ctx, h.redis, []string{key}, ttlSeconds).Int64()
		},
	)
}

func enforceVoteRateLimitWithFuncs(
	c *fiber.Ctx,
	key string,
	maxHits int64,
	nextCount func(context.Context, string) (int64, error),
) error {
	ctx := context.Background()
	count, err := nextCount(ctx, key)
	if err != nil {
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to enforce rate limit")
	}
	if count > maxHits {
		return writeError(c, fiber.StatusTooManyRequests, "TOO_MANY_REQUESTS", "too many requests")
	}
	return nil
}

func pinOKTTL(now, voteEndsAt time.Time) time.Duration {
	remaining := voteEndsAt.Sub(now)
	if remaining <= 0 {
		return 0
	}
	max := 30 * time.Minute
	if remaining < max {
		return remaining
	}
	return max
}

func retentionTTL(now, retentionEndsAt time.Time) time.Duration {
	return retentionEndsAt.Sub(now)
}

func (h *voteHandler) userIDFromAuthorization(authHeader string) (string, error) {
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
