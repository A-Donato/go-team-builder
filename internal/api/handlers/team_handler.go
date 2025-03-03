package handlers

import (
	"net/http"
	"soccer-service/internal/core/models"
	"soccer-service/internal/core/services"

	"github.com/labstack/echo/v4"
)

type TeamHandler struct {
	teamService    *services.TeamService
	playerService  *services.PlayerService
	rankingService *services.RankingService
}

func NewTeamHandler(ts *services.TeamService, ps *services.PlayerService, rs *services.RankingService) *TeamHandler {
	return &TeamHandler{
		teamService:    ts,
		playerService:  ps,
		rankingService: rs,
	}
}

type GenerateTeamsRequest struct {
	PlayerIDs []string `json:"player_ids"`
	Strategy  string   `json:"strategy"` // "balanced", "competitive", "learning"
}

type TeamGenerationResponse struct {
	TeamA    models.TeamComposition `json:"team_a"`
	TeamB    models.TeamComposition `json:"team_b"`
	Analysis TeamAnalysis           `json:"analysis"`
}

type TeamAnalysis struct {
	PredictedScore struct {
		TeamA float64 `json:"team_a"`
		TeamB float64 `json:"team_b"`
	} `json:"predicted_score"`
	TeamAStrengths []string `json:"team_a_strengths"`
	TeamBStrengths []string `json:"team_b_strengths"`
	BalanceScore   float64  `json:"balance_score"` // 0-1, how well balanced the teams are
}

func (h *TeamHandler) GenerateTeams(c echo.Context) error {
	var req GenerateTeamsRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Get players
	var players []models.Player
	for _, id := range req.PlayerIDs {
		player, err := h.playerService.GetPlayer(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Player not found: "+id)
		}
		players = append(players, player)
	}

	// Generate balanced teams
	teamA, teamB, err := h.teamService.GenerateBalancedTeams(players)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Analyze team composition
	analysis := TeamAnalysis{
		BalanceScore:   calculateBalanceScore(teamA, teamB),
		TeamAStrengths: analyzeTeamStrengths(teamA),
		TeamBStrengths: analyzeTeamStrengths(teamB),
	}

	// Predict scores based on team composition
	analysis.PredictedScore.TeamA = predictTeamScore(teamA)
	analysis.PredictedScore.TeamB = predictTeamScore(teamB)

	response := TeamGenerationResponse{
		TeamA:    teamA,
		TeamB:    teamB,
		Analysis: analysis,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *TeamHandler) GetPlayerAffinities(c echo.Context) error {
	playerID := c.Param("id")

	// Get all players
	players, err := h.playerService.ListPlayers()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Calculate affinities
	affinities := make(map[string]float64)
	for _, player := range players {
		if player.ID != playerID {
			score, err := h.teamService.CalculatePlayerAffinity(playerID, player.ID)
			if err != nil {
				continue
			}
			affinities[player.ID] = score
		}
	}

	return c.JSON(http.StatusOK, affinities)
}

// Helper functions for team analysis
func calculateBalanceScore(teamA, teamB models.TeamComposition) float64 {
	// Compare average ratings and chemistry
	ratingDiff := abs(teamA.AvgRating - teamB.AvgRating)
	chemistryDiff := abs(teamA.Chemistry - teamB.Chemistry)

	// Return a score between 0-1, where 1 is perfectly balanced
	return 1.0 - ((ratingDiff/10.0 + chemistryDiff) / 2.0)
}

func analyzeTeamStrengths(team models.TeamComposition) []string {
	var strengths []string

	// Analyze based on formation and player roles
	// This is a simplified version - would be more sophisticated in production
	if len(team.Players) > 0 {
		strengths = append(strengths, "Full team composition")
		if team.Chemistry > 0.7 {
			strengths = append(strengths, "High team chemistry")
		}
		if team.AvgRating > 7.5 {
			strengths = append(strengths, "Strong individual skills")
		}
	}

	return strengths
}

func predictTeamScore(team models.TeamComposition) float64 {
	// Simple prediction based on average rating and chemistry
	// In production, this would use ML models
	return (team.AvgRating * 0.7) + (team.Chemistry * 0.3)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
