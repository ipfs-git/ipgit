package main

import (
	"os"

	"github.com/ipfs-git/ipgit/commands"
)

func main() {
	app := commands.NewApp()

	os.Exit(app.Run(os.Args, os.Stdout))
}
