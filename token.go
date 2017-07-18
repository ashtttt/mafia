package main

import (
	"flag"
	"os/signal"
	"syscall"

	"errors"
	"os"

	"time"

	"github.com/mitchellh/colorstring"
)

type TokenCommand struct {
}

func (t *TokenCommand) Run(args []string) error {
	flags := flag.NewFlagSet("token", flag.ContinueOnError)

	flags.Usage = func() {
		colorstring.Println("[red]" + t.commandHelp())
	}
	fallow := flags.Bool("fallow", false, "")
	profile := flags.String("profile", "", "")

	flags.Parse(args[1:])
	args = flags.Args()

	if len(args) > 0 {
		flags.Usage()
	}

	if !*fallow {
		code, err := getOpt(*profile)
		if err != nil {
			return err
		}
		colorstring.Println("[green]" + code)

	} else {

		c := make(chan string)
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			for {
				code, _ := getOpt(*profile)
				c <- code
				time.Sleep(30 * time.Second)
			}
		}()

		for {
			select {
			case code := <-c:
				colorstring.Printf("\n")
				colorstring.Println("[green]" + code)
			case <-time.After(1 * time.Second):
				colorstring.Printf("[yellow].")
			case <-sigs:
				return errors.New("process interrupted. Exiting")
			}
		}
	}

	return nil
}

func (t *TokenCommand) commandHelp() string {

	var usage = `Usage: mafia token -<option>
options:
  -fallow=false  Set to ture to contineously disply token for every 30 sec, defualt is false
  -profile=none  Configured profile name to find mfa device secret. Default set to none
`
	return usage
}
