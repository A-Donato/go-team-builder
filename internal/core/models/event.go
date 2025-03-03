package models

import "time"

type EventType string

const (
	GoalScored    EventType = "goal"
	YellowCard    EventType = "yellow_card"
	RedCard       EventType = "red_card"
	Substitution  EventType = "substitution"
	Performance   EventType = "performance"
)

type Event struct {
	ID          string    `json:"id"`
	MatchID     string    `json:"match_id"`
	Type        EventType `json:"type"`
	PlayerID    string    `json:"player_id"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	Rating      float64   `json:"rating,omitempty"`
} 