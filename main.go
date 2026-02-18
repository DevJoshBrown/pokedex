package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var commands = map[string]cliCommand{}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

// Represents the JSON Structure from the PokeAPI
type RespLocationAreas struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// Represents the state of the CLI Tool
type config struct {
	Next     string
	Previous *string
}

func commandExit(cfg *config) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\n")

	for _, c := range commands {
		fmt.Printf("%v: %v\n", c.name, c.description)
	}
	return nil
}

func commandMap(cfg *config) error {

	// THE HTTP GET REQUEST TO THE POKEAPI.
	res, err := http.Get(cfg.Next)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// READ THE BODY OF THE REQUEST.
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Create a local instance of the "config" type
	// to hold the incoming data from the API
	locationsResponse := RespLocationAreas{}
	err = json.Unmarshal(body, &locationsResponse)
	if err != nil {
		return err
	}

	// Transfer the URLs from the "Fresh Data"
	// into the "persistent memory" (the cfg pointer).
	cfg.Next = locationsResponse.Next
	cfg.Previous = locationsResponse.Previous

	for _, location := range locationsResponse.Results {
		fmt.Println(location.Name)
	}

	return nil

}

func commandMapb(cfg *config) error {
	if cfg.Previous == nil {
		fmt.Println("You're on the first page")
		return nil
	} else {

	}

	res, err := http.Get(*cfg.Previous)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// READ THE BODY OF THE REQUEST.
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Create a local instance of the "config" type
	// to hold the incoming data from the API
	locationsResponse := RespLocationAreas{}
	err = json.Unmarshal(body, &locationsResponse)
	if err != nil {
		return err
	}

	// Transfer the URLs from the "Fresh Data"
	// into the "persistent memory" (the cfg pointer).
	cfg.Next = locationsResponse.Next
	cfg.Previous = locationsResponse.Previous

	for _, location := range locationsResponse.Results {
		fmt.Println(location.Name)
	}

	return nil

}

func main() {

	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			callback:    commandExit,
		},

		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},

		"map": {
			name:        "map",
			description: "displayes the names of 20 location areas in the Pokemon world",
			callback:    commandMap,
		},

		"mapb": {
			name:        "mapb",
			description: "displays the previous 20 location areas in the Pokemon world",
			callback:    commandMapb,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	// When you later call scanner.Scan it will block and wait for input until the user presses enter.

	//
	cfg := &config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: nil,
	}

	for {
		fmt.Print("Pokedex > ")
		for scanner.Scan() {
			clean := cleanInput(scanner.Text())
			if len(clean) > 0 {
				cmd, exists := commands[clean[0]]
				if exists {
					cmd.callback(cfg)
					break
				} else {
					fmt.Print("Unkown command\n")
					break
				}

			}
		}
	}
}
