package main

import "errors"

// command represents a single command invocation with its name and arguments.
type command struct {
	Name string
	Args []string
}

// commands is the registry that maps command strings (like "login") to functions.
type commands struct {
	registeredCommands map[string]func(*state, command) error
}

// register adds a new command to the map.
// name: the string the user types (e.g., "login")
// f: the function that executes the logic for that command
func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

// run looks up the command by name and executes it.
func (c *commands) run(s *state, cmd command) error {
	// Check if the command exists in our map
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	// Execute the retrieved function, passing in the app state and command args
	return f(s, cmd)
}
