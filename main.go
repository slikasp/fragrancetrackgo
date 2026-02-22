package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/slikasp/fragrancetrackgo/internal/config"
	"github.com/slikasp/fragrancetrackgo/internal/database"
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
	userDB, err := sql.Open("postgres", cfg.UserDbURL)
	if err != nil {
		log.Fatal(err)
	}
	fragranceDB, err := sql.Open("postgres", cfg.FragranceDbURL)
	if err != nil {
		log.Fatal(err)
	}

	// Create DB query instances
	userQueries := database.New(userDB)
	fragranceQueries := database.New(fragranceDB)

	// Create state to be passed to functions
	stt := &config.State{
		Users:      userQueries,
		Fragrances: fragranceQueries,
		Cfg:        &cfg,
	}

	cmds := registerCommands()

	// Parse command input
	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments provided")
		log.Fatal("Not enough arguments provided")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(stt, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}
