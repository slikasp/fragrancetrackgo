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

	// Load the database
	db, err := sql.Open("postgres", cfg.DbURL)
	dbQueries := database.New(db)

	// Create state to be passed to functions
	stt := &config.State{
		Db:  dbQueries,
		Cfg: &cfg,
	}

	cmds := registerCommands()

	// runREPL()

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
