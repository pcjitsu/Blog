package main

import (
	"Blog/internal/config"
	"Blog/internal/database"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq" // Import the postgres driver anonymously to register it
)

// state holds the runtime application state.
// It allows us to pass the database connection and configuration to our command handlers.
type state struct {
	db  *database.Queries // The SQLC-generated database wrapper
	cfg *config.Config    // The loaded configuration
}

func main() {
	// 1. Load the configuration from the .gatorconfig.json file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// 2. Open a connection to the PostgreSQL database using the URL from the config
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close() // Ensure the connection closes when the program exits

	// 3. Initialize the 'dbQueries' struct which contains our SQLC generated methods
	dbQueries := database.New(db)

	// 4. Create the program state instance
	programState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	// 5. Initialize the command registry
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	// 6. Register specific commands and their handler functions
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	// 7. Parse command line arguments. We expect at least the command name (e.g., "login")
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]  // The command name (e.g., "login")
	cmdArgs := os.Args[2:] // Any arguments following the command

	// 8. Run the command based on the name provided
	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
