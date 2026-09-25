package main

import (
	"fmt"
)

func main() {
	conf := config{}
	// Initialize constants
	if err := initConstants(&conf); err != nil {
		fmt.Println("Failed to initialize constants:", err.Error())
	}

	// Start program
	runRepl(&conf)
}
