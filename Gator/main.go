package main

import (
	"fmt"
	"os"

	config "github.com/Maxy-Maxwell/go-tutorials/Gator/internal"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Bring state into memory
	s := state{config.Read()}
	if s.c == nil {
		return fmt.Errorf("Failed to read config file.")
	}

	commandConfigs := instantiateCommands()

	// Retrieve command line arguments.
	// We ignore os.Args[0] because it will always be the running programs name
	cmdArgs := os.Args[1:]
	if len(cmdArgs) == 0 {
		return fmt.Errorf("No arguments provided: %v", cmdArgs)
	}

	// Parse args into nice command struct
	cmd := command{name: cmdArgs[0], arguments: cmdArgs[1:]}

	// Retrieve command and throw error if not defined
	handler, ok := commandConfigs.cmdNameToFunc[cmd.name]
	if !ok {
		return fmt.Errorf("Failed to find command: %v", cmd.name)
	}

	// Run handler (command with args)
	err := handler(&s, cmd)
	if err != nil {
		return fmt.Errorf("Failed to run command: %v, with args: %v", cmd.name, cmd.arguments)
	}

	return nil
}

func prettyPrintConfig(c *config.Config) {
	fmt.Println("Current_user_name:", c.Current_user_name)
	fmt.Println("Db_url:", c.Db_url)
}

func instantiateCommands() commands {
	cmds := commands{
		cmdNameToFunc: map[string]func(s *state, cmd command) error{
			"login": handlerLogin,
		},
	}

	return cmds
}
