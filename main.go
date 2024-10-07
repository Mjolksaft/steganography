package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	mainMenu int = iota
	secondMenu
)

type CLIcommand struct {
	Name        string
	Description string
	Callback    func()
}

var currentMenu = secondMenu

var Menus = []map[string]CLIcommand{
	{
		"help": {
			Name:        "Help",
			Description: "describes what it does ",
			Callback:    helpCommand,
		},
		"test1": {
			Name:        "Help",
			Description: "describes what it does ",
			Callback:    helpCommand,
		},
		"poopy": {
			Name:        "Help",
			Description: "describes what it does ",
			Callback:    helpCommand,
		},
	},
	{
		"help2": {
			Name:        "Help",
			Description: "describes what it does ",
			Callback:    helpCommand,
		},
		"test2": {
			Name:        "Help",
			Description: "describes what it does ",
			Callback:    helpCommand,
		},
		"poopy2": {
			Name:        "Help",
			Description: "describes what it does ",
			Callback:    helpCommand,
		},
	},
}

func main() {
	ClILoop()

}

func helpCommand() {

}

func ClILoop() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("hello world > ")
	_, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("error reading string: %w", err)
	}

}
