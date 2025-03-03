package repository

import (
	"context"
	"soccer-service/internal/core/models"

	"github.com/jackc/pgx/v5"
)

type MatchRepository struct {
	db *pgx.Conn
}

func NewMatchRepository(db *pgx.Conn) *MatchRepository {
	return &MatchRepository{db: db}
}

func (r *MatchRepository) Create(match models.Match) error {
	query := `
		INSERT INTO matches (id, date, home_team, away_team, home_score, away_score, home_mvp, away_mvp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(context.Background(), query,
		match.ID,
		match.Date,
		match.HomeTeam.ID,
		match.AwayTeam.ID,
		match.HomeScore,
		match.AwayScore,
		match.HomeMVP,
		match.AwayMVP,
	)

	return err
}

func (r *MatchRepository) Get(id string) (models.Match, error) {
	query := `SELECT id, date, home_team, away_team, home_score, away_score, home_mvp, away_mvp FROM matches WHERE id = $1`

	var match models.Match
	match.HomeTeam = models.Team{}
	match.AwayTeam = models.Team{}

	var homeTeamID, awayTeamID string
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&match.ID,
		&match.Date,
		&homeTeamID,
		&awayTeamID,
		&match.HomeScore,
		&match.AwayScore,
		&match.HomeMVP,
		&match.AwayMVP,
	)

	if err == pgx.ErrNoRows {
		return models.Match{}, ErrNotFound
	}

	if err != nil {
		return models.Match{}, err
	}

	// Set team IDs
	match.HomeTeam.ID = homeTeamID
	match.AwayTeam.ID = awayTeamID

	return match, nil
}

func (r *MatchRepository) Update(match models.Match) error {
	query := `
		UPDATE matches
		SET date = $1,
			home_team = $2,
			away_team = $3,
			home_score = $4,
			away_score = $5,
			home_mvp = $6,
			away_mvp = $7
		WHERE id = $8
	`

	result, err := r.db.Exec(context.Background(), query,
		match.Date,
		match.HomeTeam.ID,
		match.AwayTeam.ID,
		match.HomeScore,
		match.AwayScore,
		match.HomeMVP,
		match.AwayMVP,
		match.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *MatchRepository) Delete(id string) error {
	query := `DELETE FROM matches WHERE id = $1`

	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *MatchRepository) GetMatchesWithPlayers(player1ID, player2ID string) ([]models.Match, error) {
	query := `
		SELECT DISTINCT m.* 
		FROM matches m
		JOIN match_players mp1 ON m.id = mp1.match_id
		JOIN match_players mp2 ON m.id = mp2.match_id
		WHERE (mp1.player_id = $1 AND mp2.player_id = $2)
		ORDER BY m.date DESC
	`

	rows, err := r.db.Query(context.Background(), query, player1ID, player2ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []models.Match
	for rows.Next() {
		var match models.Match
		err := rows.Scan(
			&match.ID,
			&match.Date,
			&match.HomeScore,
			&match.AwayScore,
			&match.HomeMVP,
			&match.AwayMVP,
			&match.Result,
			&match.Duration,
		)
		if err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	return matches, nil
}
