package main

import (
	"errors"
	"flag"

	"github.com/mitchellh/colorstring"
)

type ConfigCommand struct {
}

func init() {
	initDB()
}

func (c *ConfigCommand) Run(args []string) error {
	flags := flag.NewFlagSet("config", flag.ContinueOnError)

	flags.Usage = func() {
		colorstring.Println("[red]" + c.commandHelp())
	}
	profileName := flags.String("profile-name", "", "")
	secretKey := flags.String("secret-key", "", "")

	flags.Parse(args[1:])
	args = flags.Args()

	if len(*profileName) <= 0 {
		flags.Usage()
		return errors.New("Profile name is required")
	} else if len(*secretKey) <= 0 {
		flags.Usage()
		return errors.New("Secret Key is required")
	}

	if len(args) > 0 {
		flags.Usage()
	}

	err := updateDB(*profileName, *secretKey)

	if err != nil {
		return err
	}
	colorstring.Printf("[green]%s \n", "Profile saved!")
	return nil
}

func (c *ConfigCommand) commandHelp() string {
	var usage = `Usage: mafia config [options]
options:
  -profile-name  AWS profile name to identify secret key.(Required)
  -secret-key	 MFA virtual device secret key.(Required)
`
	return usage
}
