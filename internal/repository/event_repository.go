package repository

import (
	"context"
	"soccer-service/internal/core/models"

	"github.com/jackc/pgx/v5"
)

type EventRepository struct {
	db *pgx.Conn
}

func NewEventRepository(db *pgx.Conn) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(event models.Event) error {
	query := `
		INSERT INTO events (id, match_id, player_id, event_type, description, timestamp, rating)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(context.Background(), query,
		event.ID,
		event.MatchID,
		event.PlayerID,
		event.Type,
		event.Description,
		event.Timestamp,
		event.Rating,
	)

	return err
}

func (r *EventRepository) GetByMatch(matchID string) ([]models.Event, error) {
	query := `SELECT id, match_id, player_id, event_type, description, timestamp, rating FROM events WHERE match_id = $1`

	rows, err := r.db.Query(context.Background(), query, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID,
			&event.MatchID,
			&event.PlayerID,
			&event.Type,
			&event.Description,
			&event.Timestamp,
			&event.Rating,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}
