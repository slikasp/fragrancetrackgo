package main

import (
	"errors"
)

// Command Dispatcher/Registry

// Command structure, e.g. login user [...] where login->Name, user->Args[0]
type command struct {
	Name string
	Args []string
}

// Command registry which links the command name to a respective action (function) of the command
type commands struct {
	registeredCommands map[string]func(*state, command) error
}

// Command dispatcher, takes the command object and looks up if c.Name command exists,
// then tries to execute the associated command function with provided c.Args
func (c *commands) run(s *state, cmd command) error {
	// Looks up the function by command name in the registeredCommands map
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	// return executes the found function
	return f(s, cmd)
}

// Allows adding new capabilities (commands) during runtime without having to modify the execution logic
func (c *commands) register(name string, f func(*state, command) error) {
	// adds a command with provided 'name' and it's function 'f'
	c.registeredCommands[name] = f
}

func registerCommands() commands {
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	// No login required
	// Misc
	cmds.register("reset", handlerReset)
	// Users
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("users", handlerUsers)
	// Frags

	return cmds
}
