package repository

import (
	"soccer-service/internal/core/models"
	"sync"
)

type PlayerRepository struct {
	players map[string]models.Player
	mutex   sync.RWMutex
}

func NewPlayerRepository() *PlayerRepository {
	return &PlayerRepository{
		players: make(map[string]models.Player),
	}
}

func (r *PlayerRepository) Create(player models.Player) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.players[player.ID] = player
	return nil
}

func (r *PlayerRepository) Get(id string) (models.Player, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	player, exists := r.players[id]
	if !exists {
		return models.Player{}, ErrNotFound
	}
	return player, nil
}

func (r *PlayerRepository) Update(player models.Player) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.players[player.ID] = player
	return nil
}

func (r *PlayerRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.players, id)
	return nil
}

func (r *PlayerRepository) List() ([]models.Player, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	players := make([]models.Player, 0, len(r.players))
	for _, player := range r.players {
		players = append(players, player)
	}
	return players, nil
} 