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
	var usage = `Usage: mafia [command...]
Commands:
  token		Generates next OTP Token code
  keys		Generates AWS keys by calling sts GetSessionToken 
`
	colorstring.Println("[red]" + usage)
}
