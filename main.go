package main

// CH1 L3 https://www.boot.dev/lessons/dca1352a-7600-4d1d-bfdf-f9d741282e55

import (
	"fmt"
	"os"
	"github.com/neixir/gator/internal/config"
)

type state struct {
	config *config.Config
}

// CH1 L3
// A command contains a name and a slice of string arguments.
// For example, in the case of the login command, the name would be "login"
// and the handler will expect the arguments slice to contain one string, the username.
type command struct {
	name string
	args []string
}

// CH1 L3
// Create a commands struct. This will hold all the commands the CLI can handle.
type commands struct {
	// This will be a map of command names to their handler functions.
	callback map[string]func(*state, command) error
}

// CH1 L3
// This method runs a given command with the provided state if it exists.
func (c *commands) run(s *state, cmd command) error {
	err := c.callback[cmd.name](s, cmd)
	if err != nil {
		return err
	}

	return nil
}

// CH1 L3
// This method registers a new handler function for a command name.
func (c *commands) register(name string, f func(*state, command) error) {
	c.callback[name] = f
}

// This will be the function signature of all command handlers.
// func handlerDefault(s *state, cmd command) error {
// }

// CH1 L3
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("missing argument <username>")
	}

	username := cmd.args[0]
	err := s.config.SetUser(username)
	if err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}

	fmt.Printf("User %s set", username)

	return nil
}

func main() {
	status := state{}
	
	// CH1 L2-3
    // Read the config file.
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config file: %v", err)
	}

	status.config = &cfg
	
	// CH1 L3 Create a new instance of the commands struct with an initialized map of handler functions.
	listOfCommands := commands{
		callback: make(map[string]func(*state, command) error),
	}
	
	// CH1 L3 Register a handler function for the login command.
	listOfCommands.register("login", handlerLogin)

	// CH1 L3 Use os.Args to get the command-line arguments passed in by the user.
	if len(os.Args) < 2 {
		fmt.Println("Please provide a command.")
		os.Exit(1)
	}
	
	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}
	
	// Run the command
	err = listOfCommands.run(&status, cmd)
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
		os.Exit(1)
	}
}