package services

import (
	"soccer-service/internal/core/models"
	"soccer-service/internal/repository"
)

type MatchService struct {
	matchRepo  *repository.MatchRepository
	playerRepo *repository.PlayerRepository
}

func NewMatchService(matchRepo *repository.MatchRepository, playerRepo *repository.PlayerRepository) *MatchService {
	return &MatchService{
		matchRepo:  matchRepo,
		playerRepo: playerRepo,
	}
}

func (s *MatchService) CreateMatch(match models.Match) error {
	return s.matchRepo.Create(match)
}

func (s *MatchService) GetMatch(id string) (models.Match, error) {
	return s.matchRepo.Get(id)
}

func (s *MatchService) UpdateMatch(match models.Match) error {
	return s.matchRepo.Update(match)
}

func (s *MatchService) DeleteMatch(id string) error {
	return s.matchRepo.Delete(id)
}
