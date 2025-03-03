package repository

import (
	"soccer-service/internal/core/models"
	"sync"
)

type MatchRepository struct {
	matches map[string]models.Match
	mutex   sync.RWMutex
}

func NewMatchRepository() *MatchRepository {
	return &MatchRepository{
		matches: make(map[string]models.Match),
	}
}

func (r *MatchRepository) Create(match models.Match) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.matches[match.ID] = match
	return nil
}

func (r *MatchRepository) Get(id string) (models.Match, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	match, exists := r.matches[id]
	if !exists {
		return models.Match{}, ErrNotFound
	}
	return match, nil
}

func (r *MatchRepository) Update(match models.Match) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.matches[match.ID] = match
	return nil
}

func (r *MatchRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.matches, id)
	return nil
} 