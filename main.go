package main

import (
	"os"

	"github.com/mitchellh/colorstring"
)

func main() {

	cli := NewCLI()
	cli.Args = os.Args[1:]

	err := cli.Run()

	if err != nil {
		colorstring.Println("[red]" + err.Error())
	}

	os.Exit(0)

}
