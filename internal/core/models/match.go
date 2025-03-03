package models

import "time"

type Match struct {
	ID          string    `json:"id"`
	Date        time.Time `json:"date"`
	HomeTeam    Team      `json:"home_team"`
	AwayTeam    Team      `json:"away_team"`
	HomeScore   int       `json:"home_score"`
	AwayScore   int       `json:"away_score"`
	Events      []Event   `json:"events"`
	HomeMVP     string    `json:"home_mvp"`
	AwayMVP     string    `json:"away_mvp"`
}

type Team struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Players []Player `json:"players"`
} 