package main

import (
	"context"
	"fmt"
	"time"

	"Blog/internal/database"
	"github.com/google/uuid"
)

// handlerRegister handles the "register" command.
// It creates a new user in the DB and logs them in locally.
// Usage: register <name>
func handlerRegister(s *state, cmd command) error {
	// Validate input: ensure exactly one argument (the username) is provided
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)
	}

	name := cmd.Args[0]

	// Call the generated CreateUser method to insert into the database.
	// We generate a new UUID and current timestamps here.
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	// Update the local config file to remember this user is logged in
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User created successfully:")
	printUser(user)
	return nil
}

// handlerLogin handles the "login" command.
// It checks if a user exists and updates the local config.
// Usage: login <name>
func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]

	// Verify the user exists in the database
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	// Update the config file to switch the current user
	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User switched successfully!")
	return nil
}

// printUser is a helper function to print user details nicely
func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
