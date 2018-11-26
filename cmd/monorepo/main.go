package main

import (
	"os"

	"github.com/GitbookIO/monorepo"
	"github.com/GitbookIO/monorepo/repo"

	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = monorepo.VERSION
	app.Name = "monorepo"
	app.Usage = "A big home for small repos"
	app.Authors = []cli.Author{
		{
			Name:  "Aaron O'Mullan",
			Email: "aaron@gitbook.com",
		},
	}

	// Flags
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "force",
			Usage:  "Force action, may result in git force-pushes",
			EnvVar: "MONOREPO_FORCE",
		},
		cli.StringFlag{
			Name: "root",
			Usage: "Path to the root of the monorepo",
			EnvVar: "MONOREPO_ROOT",
		}
	}

	// Subcommands
	app.Commands = []cli.Command{
		{
			Name:      "list",
			ShortName: "ls",
			Action: func(ctx *cli.Context) error {
				return nil
			},
		},
		{
			Name: "pull",
			Action: func(ctx *cli.Context) error {
				return nil
			},
		},
		{
			Name: "add",
			Action: func(ctx *cli.Context) error {

			},
		},
		{
			Name: "rm",
			Action: func(ctx *cli.Context) error {
				return nil
			},
		},
	}

	// Load envs from .env file (relative to where command is launched)
	// ignore any errors from now
	godotenv.Load()

	// Run and exit on failure
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func getRoot(ctx *cli.Context) string {
	// Check CLI arg
	if v := ctx.String("root"); v != "" {
		return v
	}
	// Fallback to CWD
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}
