package repository

import (
	"context"
	"errors"
	"fmt"
	"soccer-service/internal/core/models"
	"strings"

	"github.com/google/uuid"
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
	tx, err := r.db.Begin(context.Background())
	if (err != nil) {
		return err
	}
	defer tx.Rollback(context.Background())

	// Insert base player data
	query := `
		INSERT INTO players (id, name, position, gender)
		VALUES ($1, $2, $3, $4)
	`
	_, err = tx.Exec(context.Background(), query,
		player.ID,
		player.Name,
		player.Position,
		player.Gender,
	)
	if err != nil {
		return err
	}

	// Insert player stats
	statsQuery := `
		INSERT INTO player_stats (
			player_id, goals_scored, assists, clean_sheets, 
			matches_played, average_rating, win_rate, pass_accuracy,
			ball_possession, interceptions, tackles, distance_covered
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err = tx.Exec(context.Background(), statsQuery,
		player.ID,
		player.Stats.GoalsScored,
		player.Stats.Assists,
		player.Stats.CleanSheets,
		player.Stats.MatchesPlayed,
		player.Stats.AverageRating,
		player.Stats.WinRate,
		player.Stats.PassAccuracy,
		player.Stats.BallPossession,
		player.Stats.Interceptions,
		player.Stats.Tackles,
		player.Stats.DistanceCovered,
	)
	if err != nil {
		return err
	}

	// Insert versatility positions
	for _, pos := range player.Versatility {
		_, err = tx.Exec(context.Background(),
			"INSERT INTO player_versatility (player_id, position) VALUES ($1, $2)",
			player.ID, pos)
		if err != nil {
			return err
		}
	}

	// Insert abilities
	for _, ability := range player.Abilities {
		_, err = tx.Exec(context.Background(),
			`INSERT INTO player_abilities (player_id, ability_name, ability_level, description)
			 VALUES ($1, $2, $3, $4)`,
			player.ID, ability.Name, ability.Level, ability.Description)
		if err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}

func (r *PlayerRepository) Get(id string) (models.Player, error) {
	var player models.Player
	
	// Get base player data
	baseQuery := `
		SELECT p.id, p.name, p.position, p.gender, 
			   ps.goals_scored, ps.assists, ps.clean_sheets, ps.matches_played,
			   ps.average_rating, ps.win_rate, ps.pass_accuracy, ps.ball_possession,
			   ps.interceptions, ps.tackles, ps.distance_covered
		FROM players p
		LEFT JOIN player_stats ps ON p.id = ps.player_id
		WHERE p.id = $1
	`
	
	var stats models.Stats
	err := r.db.QueryRow(context.Background(), baseQuery, id).Scan(
		&player.ID,
		&player.Name,
		&player.Position,
		&player.Gender,
		&stats.GoalsScored,
		&stats.Assists,
		&stats.CleanSheets,
		&stats.MatchesPlayed,
		&stats.AverageRating,
		&stats.WinRate,
		&stats.PassAccuracy,
		&stats.BallPossession,
		&stats.Interceptions,
		&stats.Tackles,
		&stats.DistanceCovered,
	)

	if err == pgx.ErrNoRows {
		return models.Player{}, ErrNotFound
	}
	if err != nil {
		return models.Player{}, err
	}

	player.Stats = stats

	// Get versatility positions
	versRows, err := r.db.Query(context.Background(),
		"SELECT position FROM player_versatility WHERE player_id = $1", id)
	if err != nil {
		return models.Player{}, err
	}
	defer versRows.Close()

	for versRows.Next() {
		var pos models.Position
		if err := versRows.Scan(&pos); err != nil {
			return models.Player{}, err
		}
		player.Versatility = append(player.Versatility, pos)
	}

	// Get abilities
	abRows, err := r.db.Query(context.Background(),
		`SELECT ability_name, ability_level, description 
		 FROM player_abilities WHERE player_id = $1`, id)
	if err != nil {
		return models.Player{}, err
	}
	defer abRows.Close()

	for abRows.Next() {
		var ability models.Ability
		if err := abRows.Scan(&ability.Name, &ability.Level, &ability.Description); err != nil {
			return models.Player{}, err
		}
		player.Abilities = append(player.Abilities, ability)
	}

	return player, nil
}

func (r *PlayerRepository) Update(player models.Player) error {
    // Begin transaction
    tx, err := r.db.Begin(context.Background())
    if err != nil {
        return err
    }
    defer tx.Rollback(context.Background())

    // Update player base data
    playerQuery := `
        UPDATE players 
        SET name = $2, position = $3, gender = $4
        WHERE id = $1
    `
    result, err := tx.Exec(context.Background(), playerQuery,
        player.ID,
        player.Name,
        player.Position,
        player.Gender,
    )
    if err != nil {
        return err
    }
    if result.RowsAffected() == 0 {
        return ErrNotFound
    }

    // Update player stats
    statsQuery := `
        UPDATE player_stats 
        SET goals_scored = $2, assists = $3, clean_sheets = $4,
            matches_played = $5, average_rating = $6, win_rate = $7,
            pass_accuracy = $8, ball_possession = $9, interceptions = $10,
            tackles = $11, distance_covered = $12
        WHERE player_id = $1
    `
    _, err = tx.Exec(context.Background(), statsQuery,
        player.ID,
        player.Stats.GoalsScored,
        player.Stats.Assists,
        player.Stats.CleanSheets,
        player.Stats.MatchesPlayed,
        player.Stats.AverageRating,
        player.Stats.WinRate,
        player.Stats.PassAccuracy,
        player.Stats.BallPossession,
        player.Stats.Interceptions,
        player.Stats.Tackles,
        player.Stats.DistanceCovered,
    )
    if err != nil {
        return err
    }

    // Update versatility positions
    _, err = tx.Exec(context.Background(), "DELETE FROM player_versatility WHERE player_id = $1", player.ID)
    if err != nil {
        return err
    }
    for _, pos := range player.Versatility {
        _, err = tx.Exec(context.Background(),
            "INSERT INTO player_versatility (player_id, position) VALUES ($1, $2)",
            player.ID, pos)
        if err != nil {
            return err
        }
    }

    // Update abilities
    _, err = tx.Exec(context.Background(), "DELETE FROM player_abilities WHERE player_id = $1", player.ID)
    if err != nil {
        return err
    }
    for _, ability := range player.Abilities {
        _, err = tx.Exec(context.Background(),
            `INSERT INTO player_abilities (player_id, ability_name, ability_level, description)
             VALUES ($1, $2, $3, $4)`,
            player.ID, ability.Name, ability.Level, ability.Description)
        if err != nil {
            return err
        }
    }

    return tx.Commit(context.Background())
}

func (r *PlayerRepository) Delete(id string) error {
    tx, err := r.db.Begin(context.Background())
    if err != nil {
        return err
    }
    defer tx.Rollback(context.Background())

    // Delete from child tables first (in correct order)
    deleteQueries := []string{
        "DELETE FROM player_abilities WHERE player_id = $1",
        "DELETE FROM player_versatility WHERE player_id = $1",
        "DELETE FROM player_affinities WHERE player_id = $1 OR target_player_id = $1",
        "DELETE FROM position_performance WHERE player_id = $1",
        "DELETE FROM play_styles WHERE player_id = $1",
        "DELETE FROM player_stats WHERE player_id = $1",
        "DELETE FROM players WHERE id = $1",
    }

    for _, query := range deleteQueries {
        result, err := tx.Exec(context.Background(), query, id)
        if err != nil {
            return err
        }
        
        // Only check rows affected for the final players table delete
        if query == deleteQueries[len(deleteQueries)-1] {
            if result.RowsAffected() == 0 {
                return ErrNotFound
            }
        }
    }

    return tx.Commit(context.Background())
}

func (r *PlayerRepository) List() ([]models.Player, error) {
	query := `
		SELECT p.id, p.name, p.position, 
			   ps.goals_scored, ps.assists, ps.clean_sheets, 
			   ps.matches_played, ps.average_rating
		FROM players p
		LEFT JOIN player_stats ps ON p.id = ps.player_id
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return []models.Player{}, err  // Return empty array instead of nil
	}
	defer rows.Close()

	players := []models.Player{} // Initialize empty slice
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
			return []models.Player{}, err // Return empty array instead of nil
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

func (r *PlayerRepository) GetPlayerAffinities(playerID string) (models.Affinities, error) {
	affinities := models.Affinities{
		PlayerScores: make(map[string]float64),
		PositionScores: make(map[models.Position]float64),
		PlayStyles: make([]string, 0),
	}

	// Get player-to-player affinities
	affRows, err := r.db.Query(context.Background(),
		`SELECT target_player_id, compatibility_score 
		 FROM player_affinities WHERE player_id = $1`, playerID)
	if err != nil {
		return affinities, err
	}
	defer affRows.Close()

	for affRows.Next() {
		var targetID string
		var score float64
		if err := affRows.Scan(&targetID, &score); err != nil {
			return affinities, err
		}
		affinities.PlayerScores[targetID] = score
	}

	// Get position performance scores
	posRows, err := r.db.Query(context.Background(),
		`SELECT position, performance_score 
		 FROM position_performance WHERE player_id = $1`, playerID)
	if err != nil {
		return affinities, err
	}
	defer posRows.Close()

	for posRows.Next() {
		var pos models.Position
		var score float64
		if err := posRows.Scan(&pos, &score); err != nil {
			return affinities, err
		}
		affinities.PositionScores[pos] = score
	}

	// Get play styles
	styleRows, err := r.db.Query(context.Background(),
		`SELECT style FROM play_styles WHERE player_id = $1`, playerID)
	if err != nil {
		return affinities, err
	}
	defer styleRows.Close()

	for styleRows.Next() {
		var style string
		if err := styleRows.Scan(&style); err != nil {
			return affinities, err
		}
		affinities.PlayStyles = append(affinities.PlayStyles, style)
	}

	return affinities, nil
}

// Helper function to strip prefix and validate UUID
func stripAndValidatePlayerID(id string) (string, error) {
    if (!strings.HasPrefix(id, "ply-")) {
        return "", fmt.Errorf("invalid player ID format")
    }
    
    // Strip the prefix and validate the remaining UUID
    uuidStr := strings.TrimPrefix(id, "ply-")
    _, err := uuid.Parse(uuidStr)
    if (err != nil) {
        return "", fmt.Errorf("invalid UUID format after prefix")
    }
    
    return uuidStr, nil
}
