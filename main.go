package main

// CH1 L3 https://www.boot.dev/lessons/dca1352a-7600-4d1d-bfdf-f9d741282e55
// CH2 L3 https://www.boot.dev/lessons/8279802c-a867-4ba6-9d06-25625bc42544
// CH2 L4 https://www.boot.dev/lessons/6619ebf8-44ab-4a2b-a536-0b17d116c15e
// CH2 L5 https://www.boot.dev/lessons/371be77c-711d-4072-8392-81732ed87512
// CH3 L1 https://www.boot.dev/lessons/7347666d-7967-4c77-84c5-a0306bee8d05
// CH3 L2 https://www.boot.dev/lessons/f0126e90-414e-4a45-b6b6-758d59af012c
// CH3 L3 https://www.boot.dev/lessons/3c66635a-cf05-471e-8ad8-ff3a80a6b177

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/neixir/gator/internal/config"
	"github.com/neixir/gator/internal/database"
	"github.com/neixir/gator/internal/rss"
)

import _ "github.com/lib/pq"

type state struct {
	db  *database.Queries
	cfg *config.Config
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
	_, ok := c.callback[cmd.name]
	if ok {
		err := c.callback[cmd.name](s, cmd)
		if err != nil {
			return err
		}

		return nil
	}
	
	return fmt.Errorf("command not found")

}

// CH1 L3
// This method registers a new handler function for a command name.
func (c *commands) register(name string, f func(*state, command) error) {
	c.callback[name] = f
}

// CH1 L3
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("missing argument <username>")
	}

	username := cmd.args[0]

	// CH2 L3
	// Update the login command handler to error (and exit with code 1) if the given username doesn't exist in the database.
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("the user does not exist. %v", err)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("User %s set.\n", username)

	return nil
}

// CH2 L3
func handlerRegister(s *state, cmd command) error {
	// Ensure that a name was passed in the args.
	if len(cmd.args) == 0 {
		return fmt.Errorf("missing argument <username>")
	}

	username := cmd.args[0]

	// Create a new user in the database.
	// It should have access to the CreateUser query through the state -> db struct.
	arg := database.CreateUserParams{
		ID: uuid.New(),				// Use the uuid.New() function to generate a new UUID for the user.
		CreatedAt: time.Now(),		// created_at and updated_at should be the current time.
		UpdatedAt: time.Now(),
		Name: username,				// Use the provided name.
	}
	// Pass context.Background() to the query to create an empty Context argument.
	newUser, err := s.db.CreateUser(context.Background(), arg)
	if err != nil {
		// Exit with code 1 if a user with that name already exists.
		// TODO Potser millor abans de crear fer GetUser?
		return fmt.Errorf("creating user. %v", err)
	}

	// Set the current user in the config to the given name.
	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("setting user. %v", err)
	}
	
	// Print a message that the user was created, and log the user's data to the console for your own debugging.
	fmt.Printf("Created new user %s.\n", username)
	fmt.Println(newUser)
	
	return nil
}

// CH2 L4
func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("resetting users table. %v", err)
	}
	
	fmt.Println("Users table has been reset.")
	return nil
}

// CH2 L5
func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("getting user list. %v", err)
	}
	
	for _, user := range users {
		username := user.Name
		if s.cfg.CurrentUserName == username {
			username = fmt.Sprintf("%s (current)", username)
		}
		fmt.Printf("* %s\n", username)
	}

	return nil
}

// CH3 L1
// Later this will be our long-running aggregator service. For now, we'll just use it to fetch
// a single feed and ensure our parsing works. It should fetch the feed found at
// https://www.wagslane.dev/index.xml and print the entire struct to the console.
func handlerAgg(s *state, cmd command) error {
	feed, err := rss.FetchFeed("https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	/*
	fmt.Println(feed.Channel.Title)
	for i, item := range feed.Channel.Item {
		fmt.Printf("%02d %v\n", i, item.Title)
	}
	*/
	fmt.Println(feed)

	return nil	
}

//CH3 L2
func handlerAddfeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("missing arguments <name> <url>")
	}

	// Obtenim nom i url del feed dels arguments
	name := cmd.args[0]
	url := cmd.args[1]
	
	// Obtenim l'usuari segons el que haguem obtingut del fitxer de configuracio
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("the user does not exist. %v", err)
	}

	// user_id, err := user.UserID

	arg := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
		Url: url,
		UserID: user.ID,
	}

	// Pass context.Background() to the query to create an empty Context argument.
	newFeed, err := s.db.CreateFeed(context.Background(), arg)
	if err != nil {
		return fmt.Errorf("creating feed. %v", err)
	}

	fmt.Println("Created new feed.")
	fmt.Println(newFeed)
	
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("getting feed list. %v", err)
	}
	
	for _, feed := range feeds {
		// Obtenim User segons id
		// TODO Pper anar be podriem crear un map fora d'aquest for
		// amb id i nom dels usuaris, aixi no hauriem de fer un query cada vegada
		username := "Unknown"
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err == nil {
			username = user.Name
		}

		fmt.Printf("* %s, %s, %v\n", feed.Name, feed.Url, username)
	}

	return nil
	
}

// This will be the function signature of all command handlers.
// func handlerDefault(s *state, cmd command) error {
// }

func main() {
	status := state{}
	
	// CH1 L2-3
    // Read the config file.
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config file: %v", err)
	}

	status.cfg = &cfg

	// CH2 L3
	dbURL := status.cfg.DbUrl
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)
	status.db = dbQueries 

	// CH1 L3 Create a new instance of the commands struct with an initialized map of handler functions.
	listOfCommands := commands{
		callback: make(map[string]func(*state, command) error),
	}
	
	//
	listOfCommands.register("login", handlerLogin)			// CH1 L3
	listOfCommands.register("register", handlerRegister)
	listOfCommands.register("reset", handlerReset)
	listOfCommands.register("users", handlerUsers)
	listOfCommands.register("agg", handlerAgg)				// CH3 L1
	listOfCommands.register("addfeed", handlerAddfeed)		// CH3 L2
	listOfCommands.register("feeds", handlerFeeds)			// CH3 L3

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