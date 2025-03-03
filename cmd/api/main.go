package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"soccer-service/internal/api/routes"
	"soccer-service/internal/core/services"
	"soccer-service/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Verify required environment variables
	requiredEnvVars := []string{
		"SUPABASE_HOST",
		"SUPABASE_PROJECT_ID",
		"SUPABASE_DB_PASSWORD",
	}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Required environment variable %s is not set", envVar)
		}
	}

	host := os.Getenv("SUPABASE_HOST")
	projectID := os.Getenv("SUPABASE_PROJECT_ID")
	dbPass := os.Getenv("SUPABASE_DB_PASSWORD")

	// Log connection details for debugging
	log.Printf("Connecting to database using transaction pooler...")

	// Format user with project ID for pooler connection
	poolUser := fmt.Sprintf("postgres.%s", projectID)

	// Construct the connection string using the transaction pooler configuration
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:6543/postgres?sslmode=require&default_query_exec_mode=simple_protocol",
		poolUser,
		dbPass,
		host,
	)

	log.Printf("Attempting database connection with SSL mode...")

	// Connect to the database with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := pgx.ParseConfig(connectionString)
	if err != nil {
		log.Fatalf("Failed to parse connection string: %v", err)
	}

	// Configure SSL mode
	config.RuntimeParams = map[string]string{
		"sslmode": "require",
	}

	log.Printf("Connecting to database with SSL verification...")
	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		log.Printf("Connection error details: %v", err)
		log.Fatalf("Failed to connect to the database. Please check your network connection and credentials.")
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
