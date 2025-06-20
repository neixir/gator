package main

import (
	"fmt"
	"github.com/neixir/gator/internal/config"
)

func main() {
	// Update the main function to:
    // Read the config file.
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config file: %v", err)
	}
	
	// Set the current user to "lane" (actually, you should use your name instead) and update the config file on disk.
	cfg.SetUser("neixir")
	// fmt.Printf("User despres de canviar pero abans de llegir: %v\n", cfg.CurrentUserName)
	
	// Read the config file again...
	cfg, err = config.Read()
	if err != nil {
		fmt.Println("Error reading config file: %v", err)
	}
	
	// ...and print the contents of the config struct to the terminal.
	fmt.Println(cfg.DbUrl)
	fmt.Println(cfg.CurrentUserName)
}