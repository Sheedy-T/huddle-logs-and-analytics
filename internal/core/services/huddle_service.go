package services

import (
	"context"
	"encoding/json"
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

// StartHuddle creates a new Huddle session
func (s *HuddleService) StartHuddle(ctx context.Context, channelID string) (*domain.Huddle, error) {
	newID := uuid.New().String()

	huddle := domain.Huddle{
		ID:        newID,
		ChannelID: channelID,
		StartedAt: time.Now().UTC(),
		IsActive:  true,
	}

	if err := s.repo.SaveSession(ctx, huddle); err != nil {
		return nil, err
	}

	return &huddle, nil
}

// LogActivity saves a Huddle event
func (s *HuddleService) LogActivity(ctx context.Context, huddleID string, userID string, event string, meta map[string]any) error {
	var metaJson []byte
	var err error

	if meta != nil {
		metaJson, err = json.Marshal(meta)
		if err != nil {
			return err
		}
	}

	logEntry := domain.HuddleLog{
		HuddleID:  huddleID,
		UserID:    userID,
		EventType: event,
		Timestamp: time.Now().UTC(),
		MetaData:  metaJson,
	}

	return s.repo.SaveLog(ctx, logEntry)
}

// GetHuddleSummary fetches participants count and duration
func (s *HuddleService) GetHuddleSummary(ctx context.Context, huddleID string) (*domain.HuddleSummary, error) {
	return s.repo.GetHuddleSummary(ctx, huddleID)
}
