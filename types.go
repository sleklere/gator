package main

import (
	"github.com/sleklere/gator/internal/config"
	"github.com/sleklere/gator/internal/database"
)

type state struct {
	db *database.Queries
	config 	*config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

