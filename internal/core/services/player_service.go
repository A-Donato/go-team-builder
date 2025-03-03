package services

import (
	"soccer-service/internal/core/models"
	"soccer-service/internal/repository"
)

type PlayerService struct {
	repo *repository.PlayerRepository
}

func NewPlayerService(repo *repository.PlayerRepository) *PlayerService {
	return &PlayerService{repo: repo}
}

func (s *PlayerService) CreatePlayer(player models.Player) error {
	return s.repo.Create(player)
}

func (s *PlayerService) GetPlayer(id string) (models.Player, error) {
	return s.repo.Get(id)
}

func (s *PlayerService) UpdatePlayer(player models.Player) error {
	return s.repo.Update(player)
}

func (s *PlayerService) DeletePlayer(id string) error {
	return s.repo.Delete(id)
}

func (s *PlayerService) ListPlayers() ([]models.Player, error) {
	return s.repo.List()
} 