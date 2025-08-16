package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
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
	if len(cmd.args) != 1 {
		return errors.New("please just provide one argument (time_between_reqs)")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("collecting feeds every %s\n", cmd.args[0])

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func addFeedHandler(s *state, cmd command, user database.User) error {
	argsLength := len(cmd.args)

	if argsLength < 2 {
		return errors.New("please provide a name and a url for the feed")
	}

	if argsLength > 2 {
		fmt.Printf("ignoring arguments after %s\n", cmd.args[1])
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
	}

	fmt.Printf("✓ feed created successfully!\n")
	fmt.Printf("  Name: %s\n", feed.Name)
	fmt.Printf("  URL: %s\n", feed.Url)
	fmt.Printf("  ID: %s\n", feed.ID)
	fmt.Printf("  User ID: %s\n", feed.UserID)

	_, err = createFeedFollowByUrl(s, feed.Url, user)
	if err != nil {
		return err
	}

	return nil
}

func listFeedsHandler(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("feeds: ")
	for _, f := range feeds {
		fmt.Printf("name: %s\n", f.Name)
		fmt.Printf("url: %s\n", f.Url)
		fmt.Printf("user: %s\n", f.UserName)
		fmt.Println("--------")
	}

	return nil
}

func followHandler(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("please just provide one argument (the feed's url)")
	}

	feedFollow, err := createFeedFollowByUrl(s, cmd.args[0], user)
	if err != nil {
		return err
	}

	fmt.Printf("user '%s' is now following feed '%s'\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func followingHandler(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsByUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, f := range follows {
		fmt.Printf("user: %s, feed name: %s, feed url: %s\n", f.UserName, f.Url, f.Name)
	}

	return nil
}

func unfollowHandler(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("please provide just one command (the feed's url)")
	}

	params := database.DeleteFeedFollowByUserAndFeedUrlParams{
		ID: user.ID,
		Url: cmd.args[0],
	}
	err := s.db.DeleteFeedFollowByUserAndFeedUrl(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Println("successfully unfollowed feed")

	return nil
}

func browseHandler(s *state, cmd command, user database.User) error {
	var limit int32 = 2

	if len(cmd.args) > 0 && cmd.args[0] != "" {
		parsedLimit, err := strconv.ParseInt(cmd.args[0], 10, 32)
		if err != nil {
			return err
		}

		limit = int32(parsedLimit)
	}

	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: limit,
	}

	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("posts len: %d", len(posts))

	for _, p := range posts {
		description := "No description"
		if p.Description.Valid {
			description = p.Description.String
		}
		fmt.Printf("* %s (%s) - %s\n", p.Title, p.PublishedAt.Format("Jan 2"), description)
	}

	return nil
}

func createFeedFollowByUrl(s *state, url string, user database.User) (database.CreateFeedFollowRow, error) {
	var feedFollow database.CreateFeedFollowRow

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return feedFollow, err
	}

	params := database.CreateFeedFollowParams{
		ID: uuid.New(),
		UserID: user.ID,
		FeedID: feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	feedFollow, err = s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return feedFollow, err
	}

	return feedFollow, nil
}

