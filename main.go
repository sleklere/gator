package main

import (
	"fmt"

	"github.com/sleklere/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config file: %v", err)
	}

	s := &state{
		config: cfg,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", loginHandler)

}
