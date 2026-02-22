package handlers

import (
	"context"
	"fmt"

	"github.com/slikasp/fragrancetrackgo/internal/config"
)

func ResetUsers(s *config.State) error {
	err := s.Users.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Could not clear users table: %v\n", err)
	}

	return nil
}
