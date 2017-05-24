package main

import (
	"errors"

	"github.com/mitchellh/colorstring"
)

type CLI struct {
	Args     []string
	Commands map[string]func() Command
}

func NewCLI() *CLI {
	return &CLI{
		Commands: map[string]func() Command{
			"token": func() Command {
				return &TokenCommand{}
			},
			"keys": func() Command {
				return &KeyCommand{}
			},
		},
	}
}

func (c *CLI) Run() error {

	err := c.processArgs()
	if err != nil {
		return err
	}
	raw, ok := c.Commands[c.Args[0]]
	if !ok {
		c.printHelp()
		return nil
	}
	command := raw()

	err = command.Run(c.Args)

	if err != nil {
		return err
	}
	return nil
}

func (c *CLI) processArgs() error {

	if len(c.Args) <= 0 {
		c.printHelp()
		return errors.New("Not enough arguments")
	}
	for _, arg := range c.Args {
		if arg == "-h" || arg == "--help" || arg == "-help" || arg == "--h" {
			c.printHelp()
			break
		}
	}
	return nil
}

func (c *CLI) printHelp() {
	var usage = `Mafia is a command line tool to automatically rotate AWS temporary access keys for MFA enabled users. 
It generates MFA token and rotate AWS temporary keys before they expire.

Usage: mafia [command...]
commands:
  token		Generates next MFA token code
  keys		Generates new AWS session keys and updates credential file under [defaul] profile
`
	colorstring.Println("[red]" + usage)
}
