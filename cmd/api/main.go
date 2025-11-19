package main

import (
	"log"
	"os"

	"github.com/Nikitaannusewicz/carwash-crm/internal/config"
	"github.com/Nikitaannusewicz/carwash-crm/internal/database"
	"github.com/Nikitaannusewicz/carwash-crm/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	// 1. Load config
	cfg := config.LoadConfig()

	logger := log.New(os.Stdout, "CARWASH-API", log.LstdFlags)

	// 2. Connect to DB
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Close database connection when main exits (graceful shutdown)
	defer db.Close()

	logger.Println("Database connection established successfully")
	// 3. Initialize server (injecting the DB)
	srv := server.NewServer(cfg, db)

	// 4. Start server
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
