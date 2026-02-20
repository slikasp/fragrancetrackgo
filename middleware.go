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
		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			fmt.Printf("User %s does not exist.\n", s.Cfg.CurrentUserName)
			os.Exit(1)
		}
		if user.Name != s.Cfg.CurrentUserName {
			fmt.Printf("User %s not logged in.\n", s.Cfg.CurrentUserName)
			os.Exit(1)
		}

		return handler(s, cmd, user)
	}
}

func middlewareLoggedInAdmin(handler func(s *config.State, cmd command, user database.User) error) func(*config.State, command) error {
	return func(s *config.State, cmd command) error {
		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			fmt.Printf("User %s does not exist.\n", s.Cfg.CurrentUserName)
			os.Exit(1)
		}
		if user.Name != s.Cfg.CurrentUserName {
			fmt.Printf("User %s not logged in.\n", s.Cfg.CurrentUserName)
			os.Exit(1)
		}
		if user.IsAdmin != true {
			fmt.Printf("User %s is not admin.\n", s.Cfg.CurrentUserName)
			os.Exit(1)
		}

		return handler(s, cmd, user)
	}
}
