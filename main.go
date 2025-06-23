package main

// CH1 L3 https://www.boot.dev/lessons/dca1352a-7600-4d1d-bfdf-f9d741282e55
// CH2 L3 https://www.boot.dev/lessons/8279802c-a867-4ba6-9d06-25625bc42544
// CH2 L4 https://www.boot.dev/lessons/6619ebf8-44ab-4a2b-a536-0b17d116c15e
// CH2 L5 https://www.boot.dev/lessons/371be77c-711d-4072-8392-81732ed87512
// CH3 L1 https://www.boot.dev/lessons/7347666d-7967-4c77-84c5-a0306bee8d05
// CH3 L2 https://www.boot.dev/lessons/f0126e90-414e-4a45-b6b6-758d59af012c
// CH3 L3 https://www.boot.dev/lessons/3c66635a-cf05-471e-8ad8-ff3a80a6b177
// CH4 L1 https://www.boot.dev/lessons/a5f72e6a-6af3-4568-9eb7-079a3809a46c
// CH4 L2 https://www.boot.dev/lessons/dbc877bf-a777-416e-ac07-f6ca9559f48c
// CH4 L3 https://www.boot.dev/lessons/b1eb06af-f46e-40c1-a64f-836248122bb0
// CH5 L1 https://www.boot.dev/lessons/096ad14b-a863-4dcf-861d-9085bfc64cf9
// CH5 L2 https://www.boot.dev/lessons/d391e27f-fbc9-4ca0-bc4c-d1a4f912bf16

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/neixir/gator/internal/config"
	"github.com/neixir/gator/internal/database"
	"github.com/neixir/gator/internal/rss"

	_ "github.com/lib/pq"
)

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
		ID:        uuid.New(), // Use the uuid.New() function to generate a new UUID for the user.
		CreatedAt: time.Now(), // created_at and updated_at should be the current time.
		UpdatedAt: time.Now(),
		Name:      username, // Use the provided name.
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

// CH3 L1 + CH5 L1
func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("missing arguments <time_between_reqs>")
	}

	// time_between_reqs is a duration string, like 1s, 1m, 1h, etc.
	// https://pkg.go.dev/time#ParseDuration
	time_between_reqs := cmd.args[0]
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		// time: unknown unit "x" in duration "10x"
		return err
	}

	fmt.Printf("Collecting feeds every %s\n", time_between_reqs)

	// Use a time.Ticker to run your scrapeFeeds function once every time_between_reqs.
	// I used a for loop to ensure that it runs immediately and then every time the ticker ticks:
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		fmt.Println(" ... clock strikes ...")
		scrapeFeeds(s)
	}

	return nil
}

// CH3 L2
func handlerAddfeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("missing arguments <name> <url>")
	}

	// Obtenim nom i url del feed dels arguments
	name := cmd.args[0]
	url := cmd.args[1]

	arg := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	// Pass context.Background() to the query to create an empty Context argument.
	feed, err := s.db.CreateFeed(context.Background(), arg)
	if err != nil {
		return fmt.Errorf("creating feed. %v", err)
	}

	fmt.Println("Created new feed.")
	fmt.Printf("* [%s] %s -- %s\n", user.Name, feed.Name, feed.Url)
	// fmt.Println(feed)

	// CH4 L1
	// It should now automatically create a feed follow record for the current user when they add a feed.
	// Es copy paste de "handleFollow", potser fer-ne metode (TODO)
	argsFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), argsFollow)
	if err != nil {
		return fmt.Errorf("creating feed_follows. %v", err)
	}

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

// CH4 L1
// It takes a single url argument and creates a new feed follow record for the current user.
// It should print the name of the feed and the current user once the record is created
// (which the query we just made should support). You'll need a query to look up feeds by URL.
func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("missing arguments <url>")
	}

	// Obtenim nom i url del feed dels arguments
	url := cmd.args[0]

	// Obtenim el feed segons el que haguem obtingut del fitxer de configuracio
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("the feed does not exist. %v", err)
	}

	arg := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), arg)
	if err != nil {
		return fmt.Errorf("creating feed_follows. %v", err)
	}

	fmt.Println("Created new follow:")
	fmt.Printf("* [%s] %s\n", user.Name, feed.Name)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	followingFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("getting following feeds for [%s] -- %v", user.Name, err)
	}

	fmt.Printf("User %s follows:\n", user.Name)
	for _, feed := range followingFeeds {
		fmt.Printf("* %s\n", feed.Name)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("missing arguments <url>")
	}

	// Obtenim nom i url del feed dels arguments
	url := cmd.args[0]

	// Obtenim el feed
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("the feed does not exist. %v", err)
	}

	arg := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	// Si no el segueix sembla que no dona error
	err = s.db.DeleteFeedFollow(context.Background(), arg)
	if err != nil {
		return fmt.Errorf("deleting feed_follows. %v", err)
	}

	fmt.Printf("User %s is not following \"%s\" anymore.\n", user.Name, feed.Name)

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int
	var err error

	if len(cmd.args) < 1 {
		limit = 2
	} else {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
	}

	arg := database.GetLimitedPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}

	newPosts, err := s.db.GetLimitedPostsForUser(context.Background(), arg)
	if err != nil {
		return fmt.Errorf("getting posts for [%s] -- %v", user.Name, err)
	}

	fmt.Printf("%d new posts.\n", len(newPosts))
	for _, post := range newPosts {
		fmt.Printf("* %s\n", post.Title)
	}

	return nil
}

// This will be the function signature of all command handlers.
// func handlerDefault(s *state, cmd command) error {
// }
// func handlerDefault(s *state, cmd command, user database.User) error {
// }

// CH4 L2
// Obtenim l'usuari segons el que haguem obtingut del fitxer de configuracio
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("the user does not exist. %v", err)
		}

		return handler(s, cmd, user)
	}
}

// CH5 L1-L2
func scrapeFeeds(s *state) error {
	// Get the next feed to fetch from the DB
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("getting next feed to fetch. %v", err)
	}
	// fmt.Printf("- Nextfeed: %v / %v (last fetched %v)\n", nextFeed.Name, nextFeed.Url, nextFeed.LastFetchedAt)

	// Mark it as fetched
	argsMark := database.MarkFeedFetchedParams{
		ID:            nextFeed.ID,
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	err = s.db.MarkFeedFetched(context.Background(), argsMark)
	if err != nil {
		return fmt.Errorf("marking feed as fetched. %v", err)
	}

	// Fetch the feed using the URL (we already wrote this function)
	fmt.Printf("# Fetching %s", nextFeed.Name)
	feed, err := rss.FetchFeed(nextFeed.Url)
	if err != nil {
		return err
	}

	// Iterate over the items in the feed and print their titles to the console.

	// Update your scraper to save posts. Instead of printing out the titles of the posts, save them to the database!
	fmt.Printf(": %d items.\n", len(feed.Channel.Item))
	for _, item := range feed.Channel.Item {
		fmt.Printf("* ADDING --> %s (%v)\n", item.Title, item.PubDate)

		// converteix item.PubDate (string) a time.Time
		// fmt.Printf("item.Pubdate: %s\n", item.PubDate)
		//pubDate, err := time.Parse("Mon, 23 Jun 2025 07:01:00 +0000", item.PubDate)
		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return fmt.Errorf("failed to parse date: %w", err)
		}

		argsCreatePost := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: sql.NullTime{Time: pubDate, Valid: true},
			FeedID:      uuid.NullUUID{UUID: nextFeed.ID, Valid: true},
		}

		_, err = s.db.CreatePost(context.Background(), argsCreatePost)
		if err != nil {
			// If you encounter an error where the post with that URL already exists, just ignore it. That will happen a lot.
			// If it's a different error, you should probably log it.
			// pq: duplicate key value violates unique constraint "posts_url_key"
			if !strings.Contains(err.Error(), "unique constraint \"posts_url_key\"") {
				//return fmt.Errorf("creating post -- %v", err)
				fmt.Println("Error creating post -- %v", err)
			}
		}

	}
	fmt.Println("")

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
	listOfCommands.register("login", handlerLogin) // CH1 L3
	listOfCommands.register("register", handlerRegister)
	listOfCommands.register("reset", handlerReset)
	listOfCommands.register("users", handlerUsers)
	listOfCommands.register("agg", handlerAgg)                                 // CH3 L1 + CH5 L1
	listOfCommands.register("addfeed", middlewareLoggedIn(handlerAddfeed))     // CH3 L2 + CH4 L2
	listOfCommands.register("feeds", handlerFeeds)                             // CH3 L3
	listOfCommands.register("follow", middlewareLoggedIn(handlerFollow))       // CH4 L1 + CH4 L2
	listOfCommands.register("following", middlewareLoggedIn(handlerFollowing)) // CH4 L1 + CH4 L2
	listOfCommands.register("unfollow", middlewareLoggedIn(handlerUnfollow))   // CH4 L3
	listOfCommands.register("browse", middlewareLoggedIn(handlerBrowse))       // CH5 L2

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
