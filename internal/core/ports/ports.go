package ports

import (
	"context"
	"github.com/Sheedy-T/huddle-backend/internal/core/domain"
)

// HuddleRepository defines HOW we talk to the Database
type HuddleRepository interface {
	SaveSession(ctx context.Context, huddle domain.Huddle) error
	SaveLog(ctx context.Context, log domain.HuddleLog) error
  GetHuddleSummary(ctx context.Context, huddleID string) (*domain.HuddleSummary, error)
}

// HuddleService defines WHAT the business logic does
type HuddleService interface {
	StartHuddle(ctx context.Context, channelID string) (*domain.Huddle, error)
	LogActivity(ctx context.Context, huddleID string, userID string, event string, meta map[string]any) error
    GetHuddleSummary(ctx context.Context, huddleID string) (*domain.HuddleSummary, error)
};