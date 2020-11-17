package main

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path"

	"gotime/cmd/commands/add"
	"gotime/cmd/commands/current"
	"gotime/cmd/commands/list"
	"gotime/cmd/commands/start"
	"gotime/cmd/commands/stop"
	"gotime/internal/timeline"

	"github.com/spf13/cobra"
)

func main() {
	ctx := context.Background()

	root := &cobra.Command{
		Use: "gotime <command>",
	}

	root.Version = "1.1.0"

	home, err := os.UserHomeDir()
	if err != nil {
		home = "/"
	}

	p := path.Join(home, ".gotime", "timeline.yaml")

	username := "guest"
	usr, err := user.Current()
	if err == nil {
		username = usr.Username
	}

	tl, err := timeline.New(p, username)
	if err != nil {
		handleErr("%v", err)
	}

	defer func() {
		err := tl.Close()
		if err != nil {
			fmt.Println(err)
		}

	}()

	root.AddCommand(
		start.New(tl),
		stop.New(tl),
		add.New(tl),
		current.New(tl),
		list.New(tl),
	)

	if err := root.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleErr(f string, a ...interface{}) {
	fmt.Printf(f, a...)
	os.Exit(1)
}
