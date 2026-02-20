package main

import (
	"fmt"
	"strconv"

	"github.com/slikasp/fragrancetrackgo/handlers"
	"github.com/slikasp/fragrancetrackgo/internal/config"
	"github.com/slikasp/fragrancetrackgo/internal/database"
)

// Wrappers for command line input

// USERS

func cmdHandlerUser(s *config.State, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Incorrect use, expected 'user <login|register|list> [args]'")
	}

	subcmd := cmd.Args[0]
	args := cmd.Args[1:]

	switch subcmd {
	case "register":
		return cmdHandlerUserRegister(s, command{Name: "register", Args: args})
	case "login":
		return cmdHandlerUserLogin(s, command{Name: "login", Args: args})
	case "list":
		return cmdHandlerUserList(s, command{Name: "list", Args: args})
	default:
		return fmt.Errorf("Unknown user command '%s'. Expected one of: login, register, list", subcmd)
	}
}

func cmdHandlerUserRegister(s *config.State, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Incorrect use, expected 'user register <user> <password>'")
	}
	name := cmd.Args[0]
	pass := cmd.Args[1]

	err := handlers.UserRegister(s, name, pass)
	if err != nil {
		return err
	}

	fmt.Printf("User %s registered successfully.\n", name)
	return nil
}

func cmdHandlerUserLogin(s *config.State, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Incorrect use, expected 'user login <user> <password>'")
	}
	name := cmd.Args[0]
	pass := cmd.Args[1]

	err := handlers.UserLogin(s, name, pass)
	if err != nil {
		return err
	}

	fmt.Println("Logged in successfully!")
	return nil
}

func cmdHandlerUserList(s *config.State, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("Incorrect use, expected 'user list'")
	}

	return handlers.UserList(s)
}

// FRAGRANCES

func cmdHandlerFragrance(s *config.State, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Incorrect use, expected 'fragrance <add|remove|update|list> [args]'")
	}

	subcmd := cmd.Args[0]
	args := cmd.Args[1:]

	// TODO: make this dynamic type and list of values
	switch subcmd {
	case "add":
		return cmdHandlerFragranceAdd(s, command{Name: "add", Args: args})
	case "remove":
		return cmdHandlerFragranceRemove(s, command{Name: "remove", Args: args})
	case "update":
		return cmdHandlerFragranceUpdate(s, command{Name: "update", Args: args})
	case "list":
		return cmdHandlerFragranceList(s, command{Name: "list", Args: args})
	default:
		return fmt.Errorf("Unknown fragrance command '%s'. Expected one of: add, remove, update, list", subcmd)
	}
}

func cmdHandlerFragranceAdd(s *config.State, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Incorrect use, expected 'fragrance add <brand> <name>'")
	}
	brand := cmd.Args[0]
	name := cmd.Args[1]

	err := handlers.FragranceAdd(s, brand, name)
	if err != nil {
		return err
	}

	fmt.Printf("Fragrance %s - %s added successfully.\n", brand, name)
	return nil
}

func cmdHandlerFragranceRemove(s *config.State, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Incorrect use, expected 'fragrance remove <id>'")
	}
	v, err := strconv.ParseInt(cmd.Args[0], 10, 32)
	if err != nil {
		return err
	}
	id := int32(v)

	err = handlers.FragranceRemove(s, id)
	if err != nil {
		return err
	}

	fmt.Printf("Fragrance with ID %d removed successfully.\n", id)
	return nil
}

func cmdHandlerFragranceUpdate(s *config.State, cmd command) error {
	if len(cmd.Args) != 3 {
		return fmt.Errorf("Incorrect use, expected 'fragrance update <id> <brand> <name>'")
	}
	v, err := strconv.ParseInt(cmd.Args[0], 10, 32)
	if err != nil {
		return err
	}
	id := int32(v)
	brand := cmd.Args[1]
	name := cmd.Args[2]

	err = handlers.FragranceUpdate(s, id, brand, name)
	if err != nil {
		return err
	}

	fmt.Printf("Fragrance with ID %d updated successfully.\n", id)
	return nil
}

func cmdHandlerFragranceList(s *config.State, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("Incorrect use, expected 'fragrance list'")
	}

	return handlers.FragranceList(s)
}

// DB

func cmdHandlerReset(s *config.State, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Incorrect use, expected 'reset <table>")
	}
	table := cmd.Args[0]

	switch table {
	case "users":
		err := handlers.ResetUsers(s)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("No table named %s", table)
	}

	fmt.Printf("Table %s has been reset.\n", table)
	return nil
}
