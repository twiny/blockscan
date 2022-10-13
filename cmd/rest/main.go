package main

import (
	"log"
	"os"

	"github.com/twiny/blockscan/cmd/rest/api"

	"github.com/urfave/cli/v2"
)

// main
func main() {
	app := &cli.App{
		Name:     "rest",
		HelpName: "rest",
		Usage:    "Chain Explorer HTTP Server",
		Version:  api.Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "`path` to config file",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			app, err := api.NewAPI(c.String("config"))
			if err != nil {
				return err
			}

			go app.Shutdown()

			return app.Start()
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Println(err)
		return
	}
}
