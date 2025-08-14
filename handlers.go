package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/sleklere/gator/internal/database"
)

func loginHandler(s *state, cmd command) error {
	argsLength := len(cmd.args)

	if argsLength == 0 {
		return errors.New("please provide an argument")
	}
	if argsLength != 1 {
		return errors.New("please provide just one argument")
	}

	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		fmt.Printf("user with name '%s' doesn't exist\n", cmd.args[0])
		os.Exit(1)
	}

	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Println("user has been set!")

	return nil
}

func registerHandler(s *state, cmd command) error {
	argsLength := len(cmd.args)

	if argsLength == 0 {
		return errors.New("please provide an argument")
	}
	if argsLength != 1 {
		return errors.New("please provide just a name to register")
	}

	params := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.args[0],
	}

	createdUser, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		fmt.Printf("error creating user with name '%s': %v\n", cmd.args[0], err)
		os.Exit(1)
	}

	fmt.Printf("user created: %v", createdUser)
	s.config.SetUser(createdUser.Name)

	return nil
}
