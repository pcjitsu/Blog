package main

import (
	"log"
	"os"

	"Blog/internal/config"
)

// state represents the shared application state.
// Right now it only holds the config, but more fields can be added later
// (e.g., database handles, clients, caches, etc.).
type state struct {
	cfg *config.Config
}

func main() {
	// Load the user's CLI configuration from ~/.gatorconfig.json.
	// If this fails, the app can't run, so we exit immediately.
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// Wrap the config inside our state struct.
	// This is passed to every command handler so they can access shared data.
	programState := &state{
		cfg: &cfg,
	}

	// Create a command registry and register supported commands.
	// Each command name maps to a handler function.
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	// The CLI must include a command name. If nothing is provided, exit.
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	// Extract the command name (e.g., "login") and any arguments after it.
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	// Run the matched command handler. All handlers return an error.
	// If the handler fails, we log and exit.
	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}

// brew services start postgresql@17
// psql postgres
// goose postgres "postgres://antonioszikos:@localhost:5432/gator" up
// psql gator
// "postgres://antonioszikos:@localhost:5432/gator?sslmode=disable
