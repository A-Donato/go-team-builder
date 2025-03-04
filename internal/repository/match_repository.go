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
	tx, err := r.db.Begin(context.Background())
	if (err != nil) {
		return err
	}
	defer tx.Rollback(context.Background())

	// Insert match data
	matchQuery := `
		INSERT INTO matches (
			id, date, home_score, away_score, 
			home_mvp, away_mvp, result, duration
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = tx.Exec(context.Background(), matchQuery,
		match.ID,
		match.Date,
		match.HomeScore,
		match.AwayScore,
		match.HomeMVP,
		match.AwayMVP,
		match.Result,
		match.Duration,
	)
	if err != nil {
		return err
	}

	// Insert team compositions
	compQuery := `
		INSERT INTO team_compositions (
			match_id, is_home, formation, avg_rating, chemistry
		) VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(context.Background(), compQuery,
		match.ID, true, match.HomeTeam.Composition.Formation,
		match.HomeTeam.Composition.AvgRating, match.HomeTeam.Composition.Chemistry)
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(), compQuery,
		match.ID, false, match.AwayTeam.Composition.Formation,
		match.AwayTeam.Composition.AvgRating, match.AwayTeam.Composition.Chemistry)
	if err != nil {
		return err
	}

	// Insert match players
	playerQuery := `
		INSERT INTO match_players (
			match_id, player_id, is_home, position, role
		) VALUES ($1, $2, $3, $4, $5)
	`
	
	// Insert home team players
	for _, player := range match.HomeTeam.Composition.Players {
		_, err = tx.Exec(context.Background(), playerQuery,
			match.ID, player.PlayerID, true, player.Position, player.Role)
		if err != nil {
			return err
		}
	}

	// Insert away team players
	for _, player := range match.AwayTeam.Composition.Players {
		_, err = tx.Exec(context.Background(), playerQuery,
			match.ID, player.PlayerID, false, player.Position, player.Role)
		if err != nil {
			return err
		}
	}

	// Insert match analysis
	analysisQuery := `
		INSERT INTO match_analysis (
			match_id, possession_home, possession_away,
			shots_home, shots_away, passes_home, passes_away,
			fouls_home, fouls_away
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = tx.Exec(context.Background(), analysisQuery,
		match.ID,
		match.Analysis.PossessionHome,
		match.Analysis.PossessionAway,
		match.Analysis.ShotsHome,
		match.Analysis.ShotsAway,
		match.Analysis.PassesHome,
		match.Analysis.PassesAway,
		match.Analysis.FoulsHome,
		match.Analysis.FoulsAway,
	)
	if err != nil {
		return err
	}

	// Insert team synergy data
	for _, synergy := range match.Analysis.TeamSynergy {
		_, err = tx.Exec(context.Background(),
			`INSERT INTO team_synergy (match_id, player1_id, player2_id, synergy_score)
			 VALUES ($1, $2, $3, $4)`,
			match.ID, synergy.Player1ID, synergy.Player2ID, synergy.Score)
		if err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
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
