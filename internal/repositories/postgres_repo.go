package repositories

import (
	"context"
	"database/sql"

	"github.com/Sheedy-T/huddle-backend/internal/core/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// SaveSession inserts a new huddle session
func (r *PostgresRepository) SaveSession(ctx context.Context, h domain.Huddle) error {
	query := `
		INSERT INTO huddles (id, channel_id, started_at, is_active) 
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, h.ID, h.ChannelID, h.StartedAt, h.IsActive)
	return err
}

// SaveLog inserts a new huddle log with JSON metadata
func (r *PostgresRepository) SaveLog(ctx context.Context, l domain.HuddleLog) error {
	query := `
		INSERT INTO huddle_logs 
			(huddle_id, user_id, event_type, event_timestamp, meta_data) 
		VALUES ($1, $2, $3, $4, $5::jsonb)
	`

	_, err := r.db.ExecContext(ctx, query, l.HuddleID, l.UserID, l.EventType, l.Timestamp, l.MetaData)
	return err
}

// GetHuddleSummary returns a summary of participants and duration
func (r *PostgresRepository) GetHuddleSummary(ctx context.Context, huddleID string) (*domain.HuddleSummary, error) {
	query := `
		SELECT 
			COUNT(DISTINCT user_id),
			EXTRACT(EPOCH FROM (MAX(event_timestamp) - MIN(event_timestamp)))
		FROM huddle_logs 
		WHERE huddle_id = $1
	`

	var summary domain.HuddleSummary
	summary.HuddleID = huddleID

	err := r.db.QueryRowContext(ctx, query, huddleID).Scan(
		&summary.TotalParticipants,
		&summary.DurationSeconds,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &summary, nil
};
