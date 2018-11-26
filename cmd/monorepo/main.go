package main

import (
	"fmt"
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
			Name:   "root",
			Usage:  "Path to the root of the monorepo",
			EnvVar: "MONOREPO_ROOT",
		},
	}

	// Subcommands
	app.Commands = []cli.Command{
		{
			Name:      "list",
			ShortName: "ls",
			Action: func(ctx *cli.Context) error {
				return withRepo(ctx, func(r *repo.Repo) error {
					return r.List()
				})
			},
		},
		{
			Name:      "pull",
			ArgsUsage: "[subrepo] [path] [ref]",
			Action: func(ctx *cli.Context) error {
				args := ctx.Args()
				url := args.Get(0)
				subkey := args.Get(1)
				ref := args.Get(2)

				// Pull all
				if url == "" && subkey == "" && ref == "" {
					return withRepo(ctx, func(r *repo.Repo) error {
						return r.Pull(ctx.Bool("force"))
					})
				}

				// Update existing
				if url != "" && subkey == "" && ref == "" {
					return withRepo(ctx, func(r *repo.Repo) error {
						return r.PullSub(url, ctx.Bool("force"))
					})
				}

				// Add new
				if url != "" && subkey != "" {
					return withRepo(ctx, func(r *repo.Repo) error {
						return r.Add(url, subkey, ref)
					})
				}

				return fmt.Errorf("Invalid args")
			},
		},
		{
			Name: "add",
			Action: func(ctx *cli.Context) error {
				return withRepo(ctx, func(r *repo.Repo) error {
					return r.Add("https://github.com/urfave/cli", "a", "master")
				})
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
		fmt.Println("Failed:", err.Error())
		os.Exit(1)
	}
}

func withRepo(ctx *cli.Context, fn func(r *repo.Repo) error) error {
	// Open repo
	r, err := repo.Open(getRoot(ctx))
	if err != nil {
		return err
	}
	// Do func
	if err := fn(r); err != nil {
		return err
	}
	return nil
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
