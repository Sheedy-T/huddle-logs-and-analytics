package repositories

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/Sheedy-T/huddle-backend/internal/core/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) SaveSession(ctx context.Context, h domain.Huddle) error {
	query := `INSERT INTO huddles (id, channel_id, started_at, is_active) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, h.ID, h.ChannelID, h.StartedAt, h.IsActive)
	return err
}

func (r *PostgresRepository) SaveLog(ctx context.Context, l domain.HuddleLog) error {
	query := `INSERT INTO huddle_logs (huddle_id, user_id, event_type, event_timestamp, meta_data) VALUES ($1, $2, $3, $4, $5)`
	
	// Convert metadata map to JSON for Postgres
	metaJson, err := json.Marshal(l.MetaData)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, l.HuddleID, l.UserID, l.EventType, l.Timestamp, metaJson)
	return err
};

func (r *PostgresRepository) GetHuddleSummary(ctx context.Context, huddleID string) (*domain.HuddleSummary, error) {
	// SQL to count distinct users and calculate the time difference (in seconds) 
	// between the first and last logged event for the huddle.
	query := `
		SELECT 
			COUNT(DISTINCT user_id),
			EXTRACT(EPOCH FROM (MAX(event_timestamp) - MIN(event_timestamp)))
		FROM huddle_logs 
		WHERE huddle_id = $1
	`
	
	var summary domain.HuddleSummary
	summary.HuddleID = huddleID
	
	// Scan the result from the single row query
	// We use QueryRowContext because we only expect one result row.
	err := r.db.QueryRowContext(ctx, query, huddleID).Scan(
		&summary.TotalParticipants, 
		&summary.DurationSeconds,
	)
	
	// If no logs exist, DurationSeconds will be null/0 and TotalParticipants will be 0.
	// The standard library sql.QueryRowContext handles scanning NULLs into float64 as 0.
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	
	return &summary, nil
};