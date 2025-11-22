package services

import (
	"context"
	"time"

	"github.com/Sheedy-T/huddle-backend/internal/core/domain"
	"github.com/Sheedy-T/huddle-backend/internal/core/ports"
	"github.com/google/uuid"
)

type HuddleService struct {
	repo ports.HuddleRepository
}

func NewHuddleService(repo ports.HuddleRepository) *HuddleService {
	return &HuddleService{repo: repo}
}

func (s *HuddleService) StartHuddle(ctx context.Context, channelID string) (*domain.Huddle, error) {
	// Business Logic: Generate a new UUID here
	newID := uuid.New().String()

	huddle := domain.Huddle{
		ID:        newID,
		ChannelID: channelID,
		StartedAt: time.Now().In(time.UTC),
		IsActive:  true,
	}

	if err := s.repo.SaveSession(ctx, huddle); err != nil {
		return nil, err
	}

	return &huddle, nil
}

func (s *HuddleService) LogActivity(ctx context.Context, huddleID string, userID string, event string, meta map[string]any) error {
	logEntry := domain.HuddleLog{
		HuddleID:  huddleID,
		UserID:    userID,
		EventType: event,
		Timestamp: time.Now().In(time.UTC),
		MetaData:  meta,
	}

	return s.repo.SaveLog(ctx, logEntry)
}

func (s *HuddleService) GetHuddleSummary(ctx context.Context, huddleID string) (*domain.HuddleSummary, error) {
	return s.repo.GetHuddleSummary(ctx, huddleID)
};