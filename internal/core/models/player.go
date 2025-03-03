package models

type Position string

const (
	Goalkeeper  Position = "goalkeeper"
	Defender    Position = "defender"
	Midfielder  Position = "midfielder"
	Forward     Position = "forward"
)

type Player struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Position Position `json:"position"`
	Stats    Stats    `json:"stats"`
}

type Stats struct {
	GoalsScored    int     `json:"goals_scored"`
	Assists        int     `json:"assists"`
	CleanSheets    int     `json:"clean_sheets"`
	MatchesPlayed  int     `json:"matches_played"`
	AverageRating  float64 `json:"average_rating"`
} 