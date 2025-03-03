package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewDBConnection() (*sql.DB, error) {
	config := DBConfig{
		Host:     os.Getenv("SUPABASE_HOST"),
		Port:     os.Getenv("SUPABASE_PORT"),
		User:     os.Getenv("SUPABASE_USER"),
		Password: os.Getenv("SUPABASE_PASSWORD"),
		DBName:   os.Getenv("SUPABASE_DB_NAME"),
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	return db, nil
}
