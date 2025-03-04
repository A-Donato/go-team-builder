package models

import "time"

type EventType string

const (
	GoalScored   EventType = "goal"
	Assist       EventType = "assist"
	Pass         EventType = "pass"
	Tackle       EventType = "tackle"
	Save         EventType = "save"
	Interception EventType = "interception"
	YellowCard   EventType = "yellow_card"
	RedCard      EventType = "red_card"
	Substitution EventType = "substitution"
	Performance  EventType = "performance"
)

type Event struct {
	ID             string      `json:"id"`
	MatchID        string      `json:"match_id"`
	Type           EventType   `json:"type"`
	PlayerID       string      `json:"player_id"`
	TargetPlayerID string      `json:"target_player_id,omitempty"` // For passes, assists, tackles
	Description    string      `json:"description"`
	Timestamp      time.Time   `json:"timestamp"`
	Rating         float64     `json:"rating,omitempty"`
	Position       Position    `json:"position"` // Player's position at time of event
	Coordinates    Coordinates `json:"coordinates"`
	Success        bool        `json:"success"` // Whether the action was successful
	Impact         float64     `json:"impact"`  // Impact score of the event (0-1)
	CreatedAt      time.Time   `json:"created_at"` // Add this field
}

type Coordinates struct {
	X float64 `json:"x"` // Normalized field position (0-1)
	Y float64 `json:"y"` // Normalized field position (0-1)
}
