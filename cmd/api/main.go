package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"soccer-service/internal/api/routes"
	"soccer-service/internal/core/services"
	"soccer-service/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Construct the connection string from Supabase credentials
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		"postgres",                        // default supabase user
		os.Getenv("SUPABASE_DB_PASSWORD"), // from your .env file
		os.Getenv("SUPABASE_PROJECT_ID")+".supabase.co", // your project host
		"5432",     // default postgres port
		"postgres", // default database name
	)

	// Add SSL mode
	connectionString += "?sslmode=require"

	// Connect to the database
	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	// Test the connection
	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("Connected to:", version)

	// Initialize Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize repositories with database connection
	playerRepo := repository.NewPlayerRepository(conn)
	matchRepo := repository.NewMatchRepository(conn)
	eventRepo := repository.NewEventRepository(conn)

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
