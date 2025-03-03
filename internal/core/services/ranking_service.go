package services

import (
	"soccer-service/internal/core/models"
	"soccer-service/internal/repository"
	"sort"
)

type RankingService struct {
	playerRepo *repository.PlayerRepository
	matchRepo  *repository.MatchRepository
}

type PlayerRanking struct {
	PlayerID string  `json:"player_id"`
	Name     string  `json:"name"`
	Score    float64 `json:"score"`
	Position string  `json:"position"`
}

func NewRankingService(playerRepo *repository.PlayerRepository, matchRepo *repository.MatchRepository) *RankingService {
	return &RankingService{
		playerRepo: playerRepo,
		matchRepo:  matchRepo,
	}
}

func (s *RankingService) GetTopScorers(limit int) ([]PlayerRanking, error) {
	players, err := s.playerRepo.List()
	if err != nil {
		return nil, err
	}

	var rankings []PlayerRanking
	for _, p := range players {
		rankings = append(rankings, PlayerRanking{
			PlayerID: p.ID,
			Name:     p.Name,
			Score:    float64(p.Stats.GoalsScored),
			Position: string(p.Position),
		})
	}

	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Score > rankings[j].Score
	})

	if len(rankings) > limit {
		rankings = rankings[:limit]
	}
	return rankings, nil
}

func (s *RankingService) GetTopAssists(limit int) ([]PlayerRanking, error) {
	players, err := s.playerRepo.List()
	if err != nil {
		return nil, err
	}

	var rankings []PlayerRanking
	for _, p := range players {
		rankings = append(rankings, PlayerRanking{
			PlayerID: p.ID,
			Name:     p.Name,
			Score:    float64(p.Stats.Assists),
			Position: string(p.Position),
		})
	}

	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Score > rankings[j].Score
	})

	if len(rankings) > limit {
		rankings = rankings[:limit]
	}
	return rankings, nil
}

func (s *RankingService) GetTopRated(limit int) ([]PlayerRanking, error) {
	players, err := s.playerRepo.List()
	if err != nil {
		return nil, err
	}

	var rankings []PlayerRanking
	for _, p := range players {
		// Calculate composite score based on multiple factors
		score := s.calculateCompositeScore(p)
		rankings = append(rankings, PlayerRanking{
			PlayerID: p.ID,
			Name:     p.Name,
			Score:    score,
			Position: string(p.Position),
		})
	}

	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Score > rankings[j].Score
	})

	if len(rankings) > limit {
		rankings = rankings[:limit]
	}
	return rankings, nil
}

func (s *RankingService) calculateCompositeScore(p models.Player) float64 {
	// Weighted scoring based on multiple performance factors
	weights := map[string]float64{
		"rating":     0.3,
		"winRate":    0.2,
		"possession": 0.15,
		"passing":    0.15,
		"physical":   0.2,
	}

	score := weights["rating"] * p.Stats.AverageRating
	score += weights["winRate"] * p.Stats.WinRate
	score += weights["possession"] * p.Stats.BallPossession
	score += weights["passing"] * p.Stats.PassAccuracy
	score += weights["physical"] * (p.Stats.DistanceCovered / 10.0) // Normalize distance

	// Position-specific bonuses
	switch p.Position {
	case models.Forward:
		score += float64(p.Stats.GoalsScored) * 0.1
	case models.Midfielder:
		score += float64(p.Stats.Assists) * 0.1
	case models.Defender:
		score += float64(p.Stats.CleanSheets) * 0.1
	}

	return score
}

func (s *RankingService) GetPositionSpecialists(pos models.Position, limit int) ([]PlayerRanking, error) {
	players, err := s.playerRepo.List()
	if err != nil {
		return nil, err
	}

	var rankings []PlayerRanking
	for _, p := range players {
		if p.Position == pos {
			score := s.calculatePositionScore(p, pos)
			rankings = append(rankings, PlayerRanking{
				PlayerID: p.ID,
				Name:     p.Name,
				Score:    score,
				Position: string(p.Position),
			})
		}
	}

	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Score > rankings[j].Score
	})

	if len(rankings) > limit {
		rankings = rankings[:limit]
	}
	return rankings, nil
}

func (s *RankingService) calculatePositionScore(p models.Player, pos models.Position) float64 {
	baseScore := s.calculateCompositeScore(p)

	// Add position-specific metrics
	switch pos {
	case models.Forward:
		return baseScore * (1 + float64(p.Stats.GoalsScored)/100.0)
	case models.Midfielder:
		return baseScore * (1 + float64(p.Stats.Assists)/100.0)
	case models.Defender:
		return baseScore * (1 + float64(p.Stats.CleanSheets)/50.0)
	case models.Goalkeeper:
		return baseScore * (1 + float64(p.Stats.CleanSheets)/30.0)
	default:
		return baseScore
	}
}
