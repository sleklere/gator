package main

import (
	"fmt"
	"os"

	"github.com/sleklere/gator/internal/config"
)

func main() {

	args := os.Args

	if len(args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config file: %v\n", err)
	}

	s := &state{
		config: cfg,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", loginHandler)

	cmdName := args[1]
	cmdArgs := args[2:]

	loginCmd := command{
		name: cmdName,
		args: cmdArgs,
	}

	err = cmds.run(s, loginCmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
