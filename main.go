package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/sleklere/gator/internal/config"
	"github.com/sleklere/gator/internal/database"
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

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Printf("error connecting to DB: %v", err)
	}
	dbQueries := database.New(db)

	s := &state{
		config: cfg,
		db: dbQueries,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", loginHandler)
	cmds.register("register", registerHandler)
	cmds.register("reset", resetHandler)
	cmds.register("users", listUsersHandler)
	cmds.register("agg", aggHandler)
	cmds.register("addfeed", addFeedHandler)
	cmds.register("feeds", listFeedsHandler)
	cmds.register("follow", followHandler)
	cmds.register("following", followingHandler)

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
