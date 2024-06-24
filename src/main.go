package main

import (
	"github.com/urfave/cli/v2"
	"go-http-server/app"
	"go-http-server/param"
	"log"
	"os"
)

func main() {
	cliApp := &cli.App{
		Name:  "spa-to-http",
		Flags: param.Flags,
		Action: func(c *cli.Context) error {
			params, err := param.ContextToParams(c)
			if err != nil {
				return err
			}

			newApp := app.NewApp(params)
			go newApp.CompressFiles()
			newApp.Listen()

			return nil
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
