package main

import (
	"context"
	"fmt"
	"os"

	"github.com/slikasp/fragrancetrackgo/internal/config"
	"github.com/slikasp/fragrancetrackgo/internal/database"
)

func middlewareLoggedIn(handler func(s *config.State, cmd command, user database.User) error) func(*config.State, command) error {
	return func(s *config.State, cmd command) error {
		user, err := s.Users.GetUserByID(context.Background(), s.Cfg.CurrentUser)
		if err != nil {
			fmt.Printf("User %s does not exist.\n", s.Cfg.CurrentUser)
			os.Exit(1)
		}
		if user.ID != s.Cfg.CurrentUser {
			fmt.Printf("User %s not logged in.\n", s.Cfg.CurrentUser)
			os.Exit(1)
		}

		return handler(s, cmd, user)
	}
}

func middlewareLoggedInAdmin(handler func(s *config.State, cmd command, user database.User) error) func(*config.State, command) error {
	return func(s *config.State, cmd command) error {
		user, err := s.Users.GetUserByID(context.Background(), s.Cfg.CurrentUser)
		if err != nil {
			fmt.Printf("User %s does not exist.\n", s.Cfg.CurrentUser)
			os.Exit(1)
		}
		if user.ID != s.Cfg.CurrentUser {
			fmt.Printf("User %s not logged in.\n", s.Cfg.CurrentUser)
			os.Exit(1)
		}
		if user.IsAdmin != true {
			fmt.Printf("User %s is not admin.\n", s.Cfg.CurrentUser)
			os.Exit(1)
		}

		return handler(s, cmd, user)
	}
}
