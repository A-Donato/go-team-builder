package services

import (
	"soccer-service/internal/core/models"
	"soccer-service/internal/repository"
)

type EventService struct {
	eventRepo *repository.EventRepository
	matchRepo *repository.MatchRepository
}

func NewEventService(eventRepo *repository.EventRepository, matchRepo *repository.MatchRepository) *EventService {
	return &EventService{
		eventRepo: eventRepo,
		matchRepo: matchRepo,
	}
}

func (s *EventService) CreateEvent(event models.Event) error {
	return s.eventRepo.Create(event)
}

func (s *EventService) GetMatchEvents(matchID string) ([]models.Event, error) {
	return s.eventRepo.GetByMatch(matchID)
}
