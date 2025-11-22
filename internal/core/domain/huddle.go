package domain

import "time"

// Huddle represents a call session
type Huddle struct {
	ID        string    `json:"id"`
	ChannelID string    `json:"channel_id"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at,omitempty"` // omitempty if still active
	IsActive  bool      `json:"is_active"`
}

// HuddleLog represents a specific event (Metadata)
type HuddleLog struct {
	ID        int64     `json:"id"`
	HuddleID  string    `json:"huddle_id"`
	UserID    string    `json:"user_id"`
	EventType string    `json:"event_type"` // e.g., "JOINED", "MUTED"
	Timestamp time.Time `json:"timestamp"`
	MetaData  map[string]any `json:"meta_data,omitempty"` // Robust: flexible for extra info
}

type HuddleSummary struct {
    HuddleID          string  `json:"huddle_id"`
    TotalParticipants int     `json:"total_participants"`
    DurationSeconds   float64 `json:"duration_seconds"`
}