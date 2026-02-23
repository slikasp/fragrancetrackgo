package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/slikasp/fragrancetrackgo/internal/config"
	"github.com/slikasp/fragrancetrackgo/internal/database/localDatabase"
	"github.com/slikasp/fragrancetrackgo/internal/database/remoteDatabase"
)

func main() {
	logFile, err := setupLogging("app.log")
	if err != nil {
		log.Fatalf("failed to setup logging: %v", err)
	}
	defer logFile.Close()

	// Read config
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// Load databases
	userDB, err := sql.Open("postgres", cfg.LocalDbURL)
	if err != nil {
		log.Fatal(err)
	}
	fragranceDB, err := sql.Open("postgres", cfg.RemoteDbURL)
	if err != nil {
		log.Fatal(err)
	}

	// Create DB query instances
	userQueries := localDatabase.New(userDB)
	fragranceQueries := remoteDatabase.New(fragranceDB)

	// Create a webapp which holds current state and HTML templates
	app, err := newWebApp(userQueries, fragranceQueries)
	if err != nil {
		log.Fatalf("failed to initialize templates: %v", err)
	}

	// Start web server
	err = app.serve(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
