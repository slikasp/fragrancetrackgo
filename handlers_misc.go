package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		fmt.Printf("%v\n", err)
		return fmt.Errorf("Could not clear users table\n")
	}

	return nil
}
