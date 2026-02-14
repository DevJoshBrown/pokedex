package main

import (
	"bufio"
	"fmt"
	"os"
)

var commands = map[string]cliCommand{}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\n")

	for _, c := range commands {
		fmt.Printf("%v: %v\n", c.name, c.description)
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
	}

	scanner := bufio.NewScanner(os.Stdin)
	// When you later call scanner.Scan it will block and wait for input until the user presses enter.
	for {
		fmt.Print("Pokedex > ")
		for scanner.Scan() {
			clean := cleanInput(scanner.Text())
			if len(clean) > 0 {
				cmd, exists := commands[clean[0]]
				if exists {
					cmd.callback()
					break
				} else {
					fmt.Print("Unkown command\n")
					break
				}

			}
		}
	}
}
