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
	exits, _ := s.Users.GetUserByName(context.Background(), name)
	if exits.Name == name {
		return fmt.Errorf("User %s already exists", name)
	}
	// Create new user
	hpass, err := auth.HashPassword(pass)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Failed to hash the password.")
	}
	newUser, err := s.Users.CreateUser(context.Background(), database.CreateUserParams{
		Name:           name,
		HashedPassword: hpass,
	})
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not add user %s to database", name)
	}

	// Set new user to current user
	s.Cfg.SetUser(newUser.ID)

	log.Printf("User registered: %v:%v\n", newUser.Name, newUser.ID)

	return nil
}

func UserLogin(s *config.State, name, pass string) error {
	// Check if user exists
	user, err := s.Users.GetUserByName(context.Background(), name)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not find user %s", name)
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

	err = s.Cfg.SetUser(user.ID)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not set username %s", name)
	}

	log.Printf("User %v logged in.\n", name)
	return nil
}

func UserList(s *config.State) error {
	users, err := s.Users.GetUsers(context.Background())
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not get existing users")
	}

	currentUser, err := s.Users.GetUserByID(context.Background(), s.Cfg.CurrentUser)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not confirm current user")
	}

	for _, name := range users {
		if name == currentUser.Name {
			fmt.Printf("%s (current)\n", name)
		} else {
			fmt.Printf("%s\n", name)
		}
	}

	return nil
}
