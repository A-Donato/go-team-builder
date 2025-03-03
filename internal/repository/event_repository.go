package repository

import (
	"soccer-service/internal/core/models"
	"sync"
)

type EventRepository struct {
	events map[string]models.Event
	mutex  sync.RWMutex
}

func NewEventRepository() *EventRepository {
	return &EventRepository{
		events: make(map[string]models.Event),
	}
}

func (r *EventRepository) Create(event models.Event) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.events[event.ID] = event
	return nil
}

func (r *EventRepository) GetByMatch(matchID string) ([]models.Event, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	var matchEvents []models.Event
	for _, event := range r.events {
		if event.MatchID == matchID {
			matchEvents = append(matchEvents, event)
		}
	}
	return matchEvents, nil
} 