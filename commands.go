package main

import (
	"errors"

	"github.com/slikasp/fragrancetrackgo/internal/config"
)

// Command Dispatcher/Registry

// Command structure where Name is the top-level command and Args are the remaining tokens.
type command struct {
	Name string
	Args []string
}

// Command registry which links the command name to a respective action (function) of the command
type commands struct {
	registeredCommands map[string]func(*config.State, command) error
}

// Command dispatcher, takes the command object and looks up if c.Name command exists,
// then tries to execute the associated command function with provided c.Args
func (c *commands) run(s *config.State, cmd command) error {
	// Looks up the function by command name in the registeredCommands map
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	// return executes the found function
	return f(s, cmd)
}

// Allows adding new capabilities (commands) during runtime without having to modify the execution logic
func (c *commands) register(name string, f func(*config.State, command) error) {
	// adds a command with provided 'name' and it's function 'f'
	c.registeredCommands[name] = f
}

func registerCommands() commands {
	cmds := commands{
		registeredCommands: make(map[string]func(*config.State, command) error),
	}

	// No login required
	// User commands
	cmds.register("user", cmdHandlerUser)

	cmds.register("fragrance", middlewareLoggedInAdmin(cmdHandlerFragrance))
	// Fragrances (should only be used by admins, maybe via API)
	// add brand name
	// remove brand name
	// update

	// Scores (used by users, need login)
	// add brand name score (comment)
	// update

	// dev
	cmds.register("reset", middlewareLoggedInAdmin(cmdHandlerReset))

	return cmds
}
