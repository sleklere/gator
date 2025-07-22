package main

import (
	"errors"
	"fmt"
)

func loginHandler(s *state, cmd command) error {
	argsLength := len(cmd.args)
	if argsLength == 0 {
		return errors.New("please provide an argument")
	}
	if argsLength != 1 {
		return errors.New("please provide just one argument")
	}

	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Println("user has been set!")

	return nil
}
