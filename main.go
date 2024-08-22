package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "radlicalize",
		Usage: "A CLI tool used to clone either remote or local git repos to Radicle.xyz",
		Commands: []*cli.Command{
			{
				Name:    "local",
				Aliases: []string{"l"},
				Usage:   "Use to clone any local repos to Radicle",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "private",
						Aliases: []string{"p"},
						Usage:   "Use this flag if you want the repo to be private.",
					},
				},
				Action: func(ctx *cli.Context) error {
					private := ctx.Bool("private")
					return LocalClone(private)
				},
			},
			{
				Name:    "remote",
				Aliases: []string{"r"},
				Usage:   "Use to clone any remote public repos on Github to Radicle",
				Action: func(ctx *cli.Context) error {
					return RemoteClone()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
