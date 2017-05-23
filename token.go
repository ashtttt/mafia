package main

import (
	"encoding/base32"
	"flag"
	"os/signal"
	"syscall"

	"errors"
	"os"

	"time"

	"github.com/ashtttt/mafia/topt"
	"github.com/mitchellh/colorstring"
)

type TokenCommand struct {
}

func (t *TokenCommand) Run(args []string) error {
	flags := flag.NewFlagSet("token", flag.ContinueOnError)

	flags.Usage = func() {
		colorstring.Println("[red]" + t.commandHelp())
	}
	goOn := flags.Bool("force-continue", false, "")
	flags.Parse(args[1:])
	args = flags.Args()

	if len(args) > 0 {
		flags.Usage()
	}

	secret := os.Getenv("TOPT_SECRET")
	if len(secret) <= 0 {
		return errors.New("TOPT_SECRET environment variable is not set. Please do so")
	}
	sec, _ := base32.StdEncoding.DecodeString(secret)
	opt := topt.NewTOPT()

	if !*goOn {
		code := opt.TokenCode(sec)
		colorstring.Println("[green]" + code)

	} else {

		c := make(chan string)
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			for {
				code := opt.TokenCode(sec)
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

	var usage = `Usage: mafia token
Options:
  -force-continue=false  Set to ture to contineously disply token every 30 sec, defualt is false
`
	return usage
}
