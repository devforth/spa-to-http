package main

import (
	"github.com/urfave/cli/v2"
	"go-http-server/app"
	"go-http-server/param"
	"log"
	"os"
)

type AppRunner interface {
	CompressFiles()
	Listen()
}

func defaultNewRunner(params *param.Params) AppRunner {
	newApp := app.NewApp(params)
	return &newApp
}

func run(args []string, newRunner func(*param.Params) AppRunner) error {
	if newRunner == nil {
		newRunner = defaultNewRunner
	}

	cliApp := &cli.App{
		Name:  "spa-to-http",
		Flags: param.Flags,
		Action: func(c *cli.Context) error {
			params, err := param.ContextToParams(c)
			if err != nil {
				return err
			}

			newApp := newRunner(params)
			go newApp.CompressFiles()
			newApp.Listen()

			return nil
		},
	}

	return cliApp.Run(args)
}

func main() {
	if err := run(os.Args, nil); err != nil {
		log.Fatal(err)
	}
}
