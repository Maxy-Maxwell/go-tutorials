package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("The login handler expects a single argument - username")
	}

	err := s.c.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}

	fmt.Println("User logged in:", cmd.arguments)

	return nil
}
