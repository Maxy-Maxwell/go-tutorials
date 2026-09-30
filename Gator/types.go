package main

import (
	"fmt"

	config "github.com/Maxy-Maxwell/go-tutorials/Gator/internal/config"
	"github.com/Maxy-Maxwell/go-tutorials/Gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	cmdNameToFunc map[string]func(s *state, cmd command) error
}

func (c *commands) run(s *state, cmd command) error {
	if s == nil {
		return fmt.Errorf("Unable to run command as s (state) is nil")
	}

	err := c.cmdNameToFunc[cmd.name](s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmdNameToFunc[name] = f
}
