package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/slikasp/fragrancetrackgo/internal/auth"
	"github.com/slikasp/fragrancetrackgo/internal/config"
	"github.com/slikasp/fragrancetrackgo/internal/database"
)

func UserRegister(s *config.State, name, pass string) error {
	// Check if user exists
	exits, _ := s.Db.GetUser(context.Background(), name)
	if exits.Name == name {
		return fmt.Errorf("User %s already exists", name)
	}
	// Create new user
	hpass, err := auth.HashPassword(pass)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Failed to hash the password.")
	}
	newUser, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		Name:           name,
		HashedPassword: hpass,
	})
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not add user %s to database", name)
	}

	// Set new user to current user
	s.Cfg.SetUser(name)

	log.Printf("User registered: %v:%v\n", newUser.Name, newUser.ID)

	return nil
}

func UserLogin(s *config.State, name, pass string) error {
	// Check if user exists
	user, err := s.Db.GetUser(context.Background(), name)
	if user.Name != name {
		return fmt.Errorf("User %s is not registered", name)
	}

	// Check if password matches
	match, err := auth.CheckPasswordHash(pass, user.HashedPassword)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not check password for user %s", name)
	}
	if match != true {
		return fmt.Errorf("Password does not match")
	}

	err = s.Cfg.SetUser(name)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not set username %s", name)
	}

	log.Printf("User %v logged in.\n", name)
	return nil
}

func UserList(s *config.State) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not get existing users")
	}

	for _, name := range users {
		if name == s.Cfg.CurrentUserName {
			fmt.Printf("%s (current)\n", name)
		} else {
			fmt.Printf("%s\n", name)
		}
	}

	return nil
}
