package main

import (
	"errors"
	"fmt"
)



func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return errors.New("command not found")
	}

	err := handler(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) error {
	if name == "" {
		return errors.New("please provide a command name")
	}
	if _, ok := c.handlers[name]; ok {
		return fmt.Errorf("a handler is already registered for the key '%s'", name)
	}
	c.handlers[name] = f

	return nil
}
