package main

import (
	"log"
	"net/http"

	"threat-intel-engine/internal/api/routes"
	"threat-intel-engine/internal/config"
	"threat-intel-engine/pkg/database"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	db, err := database.NewMongoDB(cfg.MongoURI, cfg.DatabaseName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Setup routes
	router := routes.SetupRoutes(db)

	log.Printf("Threat Intelligence Engine starting on port %s...", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
