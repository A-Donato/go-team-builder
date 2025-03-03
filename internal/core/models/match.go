package models

import "time"

type MatchResult string

const (
	HomeWin MatchResult = "home_win"
	AwayWin MatchResult = "away_win"
	Draw    MatchResult = "draw"
)

type Match struct {
	ID        string        `json:"id"`
	Date      time.Time     `json:"date"`
	HomeTeam  Team          `json:"home_team"`
	AwayTeam  Team          `json:"away_team"`
	HomeScore int           `json:"home_score"`
	AwayScore int           `json:"away_score"`
	Events    []Event       `json:"events"`
	HomeMVP   string        `json:"home_mvp"`
	AwayMVP   string        `json:"away_mvp"`
	Result    MatchResult   `json:"result"`
	Duration  int           `json:"duration"` // in minutes
	Analysis  MatchAnalysis `json:"analysis"`
}

type Team struct {
	ID          string          `json:"id"`
	Composition TeamComposition `json:"composition"`
}

type TeamComposition struct {
	Players   []PlayerRole `json:"players"`
	Formation string       `json:"formation"` // e.g., "4-4-2"
	AvgRating float64      `json:"avg_rating"`
	Chemistry float64      `json:"chemistry"` // 0-1 score based on player affinities
}

type PlayerRole struct {
	PlayerID string   `json:"player_id"`
	Position Position `json:"position"`
	Role     string   `json:"role"` // e.g., "striker", "defensive_mid", "sweeper"
}

type MatchAnalysis struct {
	PossessionHome float64            `json:"possession_home"`
	PossessionAway float64            `json:"possession_away"`
	ShotsHome      int                `json:"shots_home"`
	ShotsAway      int                `json:"shots_away"`
	PassesHome     int                `json:"passes_home"`
	PassesAway     int                `json:"passes_away"`
	FoulsHome      int                `json:"fouls_home"`
	FoulsAway      int                `json:"fouls_away"`
	TeamSynergy    map[string]float64 `json:"team_synergy"` // player pair -> synergy score
}
