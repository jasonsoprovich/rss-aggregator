package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jasonsoprovich/rss-aggregator/internal/config"
	"github.com/jasonsoprovich/rss-aggregator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}
	return handler(s, cmd)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("login requires a username argument")
	}

	username := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: user does not exist")
		os.Exit(1)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}

	fmt.Printf("User set to: %s\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("register requires a username argument")
	}

	username := cmd.args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: user already exists or could not be created")
		os.Exit(1)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("error setting user in config: %w", err)
	}
	fmt.Printf("User created successfully: %s\n", username)
	log.Printf("Debug - User data: %+v\n", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error resetting users table: %w", err)
	}
	fmt.Println("Users table reset successfully")
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	dbQueries := database.New(db)

	s := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)

	args := os.Args
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "error: not enough arguments, a command is required")
		os.Exit(1)
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	err = cmds.run(s, cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
