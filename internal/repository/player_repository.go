package repository

import (
	"context"
	"errors"
	"soccer-service/internal/core/models"

	"github.com/jackc/pgx/v5"
)

var ErrNotFound = errors.New("record not found")

type PlayerRepository struct {
	db *pgx.Conn
}

func NewPlayerRepository(db *pgx.Conn) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (r *PlayerRepository) Create(player models.Player) error {
	query := `
		INSERT INTO players (id, name, position, goals_scored, assists, clean_sheets, matches_played, average_rating)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(context.Background(), query,
		player.ID,
		player.Name,
		player.Position,
		player.Stats.GoalsScored,
		player.Stats.Assists,
		player.Stats.CleanSheets,
		player.Stats.MatchesPlayed,
		player.Stats.AverageRating,
	)

	return err
}

func (r *PlayerRepository) Get(id string) (models.Player, error) {
	query := `
		SELECT id, name, position, goals_scored, assists, clean_sheets, matches_played, average_rating
		FROM players WHERE id = $1
	`

	var player models.Player
	var stats models.Stats

	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&player.ID,
		&player.Name,
		&player.Position,
		&stats.GoalsScored,
		&stats.Assists,
		&stats.CleanSheets,
		&stats.MatchesPlayed,
		&stats.AverageRating,
	)

	if err == pgx.ErrNoRows {
		return models.Player{}, ErrNotFound
	}

	if err != nil {
		return models.Player{}, err
	}

	player.Stats = stats
	return player, nil
}

func (r *PlayerRepository) Update(player models.Player) error {
	query := `
		UPDATE players 
		SET name = $2, position = $3, goals_scored = $4, assists = $5, 
			clean_sheets = $6, matches_played = $7, average_rating = $8
		WHERE id = $1
	`

	result, err := r.db.Exec(context.Background(), query,
		player.ID,
		player.Name,
		player.Position,
		player.Stats.GoalsScored,
		player.Stats.Assists,
		player.Stats.CleanSheets,
		player.Stats.MatchesPlayed,
		player.Stats.AverageRating,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *PlayerRepository) Delete(id string) error {
	query := `DELETE FROM players WHERE id = $1`

	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *PlayerRepository) List() ([]models.Player, error) {
	query := `
		SELECT id, name, position, goals_scored, assists, clean_sheets, matches_played, average_rating
		FROM players
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []models.Player
	for rows.Next() {
		var player models.Player
		var stats models.Stats

		err := rows.Scan(
			&player.ID,
			&player.Name,
			&player.Position,
			&stats.GoalsScored,
			&stats.Assists,
			&stats.CleanSheets,
			&stats.MatchesPlayed,
			&stats.AverageRating,
		)

		if err != nil {
			return nil, err
		}

		player.Stats = stats
		players = append(players, player)
	}

	return players, nil
}

func (r *PlayerRepository) UpdateAffinity(player1ID, player2ID string, score float64) error {
	query := `
		INSERT INTO player_affinities (player_id, target_player_id, compatibility_score)
		VALUES ($1, $2, $3)
		ON CONFLICT (player_id, target_player_id) 
		DO UPDATE SET compatibility_score = $3, updated_at = CURRENT_TIMESTAMP
	`

	_, err := r.db.Exec(context.Background(), query, player1ID, player2ID, score)
	return err
}
