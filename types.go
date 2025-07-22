package main

import "github.com/sleklere/gator/internal/config"

type state struct {
	config 	*config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

