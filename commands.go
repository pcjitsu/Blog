package main

import "errors"

// command represents a command with its name and arguments.
type command struct {
	Name string
	Args []string
}

// commands manages the registration and execution of commands.
type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

// Looks at passed command name and runs the corresponding function from the Map of Commands
func (c *commands) run(s *state, cmd command) error {
	//Checks to see if command exists in the map
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	//Runs the corresponding function
	return f(s, cmd)
}
