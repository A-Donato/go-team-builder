package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"soccer-service/internal/api/routes"
	"soccer-service/internal/core/services"
	"soccer-service/internal/repository"
)

func main() {
	// Initialize Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize repositories
	playerRepo := repository.NewPlayerRepository()
	matchRepo := repository.NewMatchRepository()
	eventRepo := repository.NewEventRepository()

	// Initialize services
	playerService := services.NewPlayerService(playerRepo)
	matchService := services.NewMatchService(matchRepo, playerRepo)
	eventService := services.NewEventService(eventRepo, matchRepo)
	rankingService := services.NewRankingService(playerRepo, matchRepo)

	// Setup routes
	routes.SetupRoutes(e, playerService, matchService, eventService, rankingService)

	// Start server
	log.Printf("Server starting on port 8080")
	log.Fatal(e.Start(":8080"))
} 