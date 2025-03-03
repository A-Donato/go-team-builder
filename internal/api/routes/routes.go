package routes

import (
	"soccer-service/internal/api/handlers"
	"soccer-service/internal/core/services"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, 
	playerService *services.PlayerService,
	matchService *services.MatchService,
	eventService *services.EventService,
	rankingService *services.RankingService) {

	// Create handlers
	playerHandler := handlers.NewPlayerHandler(playerService)

	// API v1 group
	v1 := e.Group("/api/v1")

	// Player routes
	players := v1.Group("/players")
	players.POST("", playerHandler.Create)
	players.GET("", playerHandler.List)
	players.GET("/:id", playerHandler.Get)
	players.PUT("/:id", playerHandler.Update)
	players.DELETE("/:id", playerHandler.Delete)

	// Additional routes will be added here for matches, events, and rankings
} 