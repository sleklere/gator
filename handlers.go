package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/sleklere/gator/internal/database"
	"github.com/sleklere/gator/internal/feed"
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

func resetHandler(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		os.Exit(1)
		return err
	}
	fmt.Println("database reset successful")
	os.Exit(0)
	return nil
}

func listUsersHandler(s *state, cmd command) error {
	users, err := s.db.ListUsers(context.Background())
	if err != nil {
		return err
	}
	for _, u := range users {
		line := "* %s"
		if u.Name == s.config.CurrentUserName {
			line += " (current)"
		}
		fmt.Printf(line, u.Name)
		fmt.Print("\n")
	}

	return nil
}

func aggHandler(s *state, cmd command) error {
	feedURL := "https://www.wagslane.dev/index.xml"
	feed, err := feed.FetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}

	fmt.Println("-----------------------")
	fmt.Print(feed)

	return nil
}

func addFeedHandler(s *state, cmd command) error {
	argsLength := len(cmd.args)

	if argsLength < 2 {
		return errors.New("please provide a name and a url for the feed")
	}

	if argsLength > 2 {
		fmt.Printf("ignoring arguments after %s\n", cmd.args[1])
	}

	user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
	}

	params := database.CreateFeedParams{
		ID: uuid.New(),
		Name: cmd.args[0],
		Url: cmd.args[1],
		UserID: user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
		// fmt.Printf("error creating feed: %v\n", err)
		// os.Exit(1)
	}

	fmt.Printf("✓ feed created successfully!\n")
	fmt.Printf("  Name: %s\n", feed.Name)
	fmt.Printf("  URL: %s\n", feed.Url)
	fmt.Printf("  ID: %s\n", feed.ID)
	fmt.Printf("  User ID: %s\n", feed.UserID)

	return nil
}
