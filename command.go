package main

type Command interface {
	Run([]string) error
}
