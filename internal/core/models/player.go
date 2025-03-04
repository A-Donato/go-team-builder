package models

import "time"

type Position string
type Gender string
type AbilityLevel string

const (
	Forward    Position = "forward"
	Midfielder Position = "midfielder"
	Defender   Position = "defender"
	Goalkeeper Position = "goalkeeper"
)

const (
	Male   Gender = "male"
	Female Gender = "female"
)

const (
	Beginner     AbilityLevel = "beginner"     // Just starting, basic understanding
	Developing   AbilityLevel = "developing"   // Growing but needs improvement
	Intermediate AbilityLevel = "intermediate" // Solid amateur level
	Advanced     AbilityLevel = "advanced"     // Above average amateur
	Expert       AbilityLevel = "expert"       // Top amateur level
)

type Player struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Position    Position   `json:"position"`
	Gender      Gender     `json:"gender"`
	Stats       Stats      `json:"stats"`
	Affinities  Affinities `json:"affinities"`
	Versatility []Position `json:"versatility"` // Secondary positions they can play
	Abilities   []Ability  `json:"abilities"`   // List of specific abilities
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Stats struct {
	GoalsScored     int       `json:"goals_scored"`
	Assists         int       `json:"assists"`
	CleanSheets     int       `json:"clean_sheets"`
	MatchesPlayed   int       `json:"matches_played"`
	AverageRating   float64   `json:"average_rating"`
	WinRate         float64   `json:"win_rate"`
	PassAccuracy    float64   `json:"pass_accuracy"`
	BallPossession  float64   `json:"ball_possession"`
	Interceptions   int       `json:"interceptions"`
	Tackles         int       `json:"tackles"`
	DistanceCovered float64   `json:"distance_covered"` // in kilometers
	UpdatedAt       time.Time `json:"updated_at"`
}

// Affinities tracks how well a player performs with others
type Affinities struct {
	// PlayerID -> Compatibility Score (0-1)
	PlayerScores map[string]float64 `json:"player_scores"`
	// Position -> Performance Score (0-1)
	PositionScores map[Position]float64 `json:"position_scores"`
	// Preferred play style tags
	PlayStyles []string `json:"play_styles"`
}

// Ability represents a specific soccer skill with a word-based level
type Ability struct {
	Name        string       `json:"name"`
	Level       AbilityLevel `json:"level"`
	Description string       `json:"description,omitempty"`
}

// AbilityLevelWeight returns a normalized score (0-1) for each ability level
func AbilityLevelWeight(level AbilityLevel) float64 {
	switch level {
	case Expert:
		return 1.0
	case Advanced:
		return 0.8
	case Intermediate:
		return 0.6
	case Developing:
		return 0.4
	case Beginner:
		return 0.2
	default:
		return 0.0
	}
}

// AbilityInfo contains metadata about an ability
type AbilityInfo struct {
	Description string
	Weights     map[Position]float64 // Position-specific weights
}

// Common abilities catalog with descriptions and position relevance
var CommonAbilities = map[string]AbilityInfo{
	// Technical Skills
	"Ball Control": {
		Description: "Ability to receive and maintain control of the ball",
		Weights: map[Position]float64{
			Forward:    1.0,  // Critical for forwards to control under pressure
			Midfielder: 0.95, // Essential for midfield play
			Defender:   0.7,  // Important but not critical
			Goalkeeper: 0.3,  // Less important for goalkeepers
		},
	},
	"Short Passing": {
		Description: "Accuracy and effectiveness of passes under 20 meters",
		Weights: map[Position]float64{
			Forward:    0.7,  // Important for link-up play
			Midfielder: 1.0,  // Critical for midfield control
			Defender:   0.85, // Important for building from back
			Goalkeeper: 0.5,  // Needed for distribution
		},
	},
	"Long Passing": {
		Description: "Accuracy and effectiveness of passes over 20 meters",
		Weights: map[Position]float64{
			Forward:    0.4,  // Less critical for forwards
			Midfielder: 0.85, // Important for switching play
			Defender:   0.9,  // Critical for long clearances
			Goalkeeper: 0.8,  // Important for goal kicks
		},
	},
	"Shooting": {
		Description: "Power and accuracy when shooting at goal",
		Weights: map[Position]float64{
			Forward:    1.0, // Essential for forwards
			Midfielder: 0.7, // Important for attacking mids
			Defender:   0.2, // Rarely needed
			Goalkeeper: 0.1, // Almost never needed
		},
	},
	"First Touch": {
		Description: "Initial ball control when receiving passes",
		Weights: map[Position]float64{
			Forward:    1.0,  // Critical under pressure
			Midfielder: 0.95, // Essential in tight spaces
			Defender:   0.8,  // Important for control
			Goalkeeper: 0.6,  // Needed for catches
		},
	},
	"Weak Foot": {
		Description: "Ability to use non-dominant foot",
		Weights: map[Position]float64{
			Forward:    0.9,  // Very important for unpredictability
			Midfielder: 0.85, // Important for all-round play
			Defender:   0.7,  // Useful for clearances
			Goalkeeper: 0.4,  // Less critical
		},
	},
	"Dribbling": {
		Description: "Ball control while moving",
		Weights: map[Position]float64{
			Forward:    0.95, // Essential for attackers
			Midfielder: 0.85, // Important for progression
			Defender:   0.4,  // Less important
			Goalkeeper: 0.1,  // Rarely needed
		},
	},

	// Physical Attributes
	"Speed": {
		Description: "Pace and acceleration",
		Weights: map[Position]float64{
			Forward:    0.95, // Critical for attacking
			Midfielder: 0.8,  // Important for transitions
			Defender:   0.85, // Essential for recovery
			Goalkeeper: 0.3,  // Less important
		},
	},
	"Stamina": {
		Description: "Endurance and fatigue resistance",
		Weights: map[Position]float64{
			Forward:    0.75, // Important for pressing
			Midfielder: 1.0,  // Critical for box-to-box
			Defender:   0.8,  // Important for 90 minutes
			Goalkeeper: 0.3,  // Less critical
		},
	},
	"Strength": {
		Description: "Physical power in duels",
		Weights: map[Position]float64{
			Forward:    0.8,  // Important for hold-up play
			Midfielder: 0.7,  // Useful in duels
			Defender:   0.95, // Critical for defending
			Goalkeeper: 0.75, // Important for presence
		},
	},
	"Agility": {
		Description: "Quick changes in direction",
		Weights: map[Position]float64{
			Forward:    0.9,  // Essential for attackers
			Midfielder: 0.85, // Important for evasion
			Defender:   0.7,  // Useful for positioning
			Goalkeeper: 0.85, // Critical for saves
		},
	},
	"Jump": {
		Description: "Vertical leap ability",
		Weights: map[Position]float64{
			Forward:    0.85, // Important for headers
			Midfielder: 0.6,  // Useful for duels
			Defender:   0.9,  // Critical for aerial duels
			Goalkeeper: 1.0,  // Essential for saves
		},
	},

	// Mental/Tactical
	"Positioning": {
		Description: "Finding optimal positions during play",
		Weights: map[Position]float64{
			Forward:    0.9,  // Critical for scoring
			Midfielder: 0.95, // Essential for play
			Defender:   1.0,  // Critical for defense
			Goalkeeper: 0.95, // Essential for saves
		},
	},
	"Game Reading": {
		Description: "Understanding and anticipating play development",
		Weights: map[Position]float64{
			Forward:    0.8,  // Important for movement
			Midfielder: 0.95, // Critical for control
			Defender:   1.0,  // Essential for defense
			Goalkeeper: 0.9,  // Critical for positioning
		},
	},
	"Decision Making": {
		Description: "Making good choices under pressure",
		Weights: map[Position]float64{
			Forward:    0.9,  // Critical for chances
			Midfielder: 0.95, // Essential for play
			Defender:   0.9,  // Critical for defense
			Goalkeeper: 1.0,  // Essential for saves
		},
	},
	"Teamwork": {
		Description: "Ability to play effectively with teammates",
		Weights: map[Position]float64{
			Forward:    0.85, // Important for attack
			Midfielder: 1.0,  // Critical for linking
			Defender:   0.9,  // Essential for line
			Goalkeeper: 0.8,  // Important for organization
		},
	},
	"Communication": {
		Description: "Effectiveness in on-field communication",
		Weights: map[Position]float64{
			Forward:    0.7,  // Useful for pressing
			Midfielder: 0.85, // Important for organization
			Defender:   0.95, // Critical for defensive line
			Goalkeeper: 1.0,  // Essential for organizing defense
		},
	},

	// Defensive Skills
	"Tackling": {
		Description: "Clean and effective defensive challenges",
		Weights: map[Position]float64{
			Forward:    0.3,  // Occasional pressing
			Midfielder: 0.75, // Important for recovery
			Defender:   1.0,  // Essential for defense
			Goalkeeper: 0.2,  // Rarely needed
		},
	},
	"Marking": {
		Description: "Following and containing opponents",
		Weights: map[Position]float64{
			Forward:    0.3, // For pressing only
			Midfielder: 0.8, // Important for midfield
			Defender:   1.0, // Essential for defense
			Goalkeeper: 0.1, // Rarely needed
		},
	},
	"Interception": {
		Description: "Reading and cutting off passes",
		Weights: map[Position]float64{
			Forward:    0.4, // For pressing
			Midfielder: 0.9, // Critical for control
			Defender:   1.0, // Essential for defense
			Goalkeeper: 0.7, // Important for crosses
		},
	},

	// Specialized Skills
	"Heading": {
		Description: "Aerial ball control and power",
		Weights: map[Position]float64{
			Forward:    0.9,  // Critical for goals
			Midfielder: 0.7,  // Important for duels
			Defender:   0.95, // Essential for defense
			Goalkeeper: 0.8,  // Important for punches
		},
	},
	"Free Kicks": {
		Description: "Set piece execution",
		Weights: map[Position]float64{
			Forward:    0.8,  // Important for direct
			Midfielder: 0.85, // Important for all types
			Defender:   0.6,  // Long free kicks
			Goalkeeper: 0.1,  // Rarely needed
		},
	},
	"Penalties": {
		Description: "Penalty kick effectiveness",
		Weights: map[Position]float64{
			Forward:    0.9, // Primary takers
			Midfielder: 0.8, // Secondary takers
			Defender:   0.4, // Rare takers
			Goalkeeper: 0.2, // Emergency only
		},
	},
	"Goalkeeping": {
		Description: "Specific goalkeeper skills",
		Weights: map[Position]float64{
			Forward:    0.0, // Not needed
			Midfielder: 0.0, // Not needed
			Defender:   0.1, // Emergency only
			Goalkeeper: 1.0, // Essential
		},
	},
}
