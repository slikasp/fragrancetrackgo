package main

import (
	"context"
	"fmt"
	"os"

	"github.com/slikasp/fragrancetrackgo/internal/database"
)

func loggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		// Check if user exists
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			fmt.Printf("User %s does not exist\n", s.cfg.CurrentUserName)
			os.Exit(1)
		}
		if user.Name != s.cfg.CurrentUserName {
			fmt.Printf("User %s not logged in\n", s.cfg.CurrentUserName)
			os.Exit(1)
		}
		return handler(s, cmd, user)
	}
}
