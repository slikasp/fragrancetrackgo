package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/slikasp/fragrancetrackgo/internal/auth"
	"github.com/slikasp/fragrancetrackgo/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Incorrect use, expected 'login <user> <password>")
	}
	name := cmd.Args[0]
	pass := cmd.Args[1]

	// Check if user exists
	user, err := s.db.GetUser(context.Background(), name)
	if user.Name != name {
		fmt.Printf("User %s is not registered\n", name)
		os.Exit(1)
	}

	// Check if password matches
	match, err := auth.CheckPasswordHash(pass, user.HashedPassword)
	if match != true {
		fmt.Printf("Password does not match\n")
		os.Exit(1)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("Could not set username %s", name)
	}

	fmt.Println("User switched successfully!")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Incorrect use, expected 'register <user> <password>")
	}
	name := cmd.Args[0]
	pass := cmd.Args[1]

	// Check if user exists
	exits, _ := s.db.GetUser(context.Background(), name)
	if exits.Name == name {
		fmt.Printf("User %s already exists\n", name)
		os.Exit(1)
	}
	// Create new user
	hpass, err := auth.HashPassword(pass)
	if err != nil {
		return fmt.Errorf("Failed to hash the password %s", err)
	}
	newUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		Name:           name,
		HashedPassword: hpass,
	})
	if err != nil {
		// fmt.Printf("%v", err)
		return fmt.Errorf("Could not add user %s to database", name)
	}

	// Set new user to current user
	s.cfg.SetUser(name)

	fmt.Println("User registered successfully!")
	log.Printf("%v\n", newUser)

	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Could not get existing users")
	}

	for _, name := range users {
		if name == s.cfg.CurrentUserName {
			fmt.Printf("%s (current)\n", name)
		} else {
			fmt.Printf("%s\n", name)
		}
	}

	return nil
}
