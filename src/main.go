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

var logFatal = log.Fatal

func defaultNewRunner(params *param.Params) AppRunner {
	newApp := app.NewApp(params)
	return &newApp
}

func run(args []string, newRunner func(*param.Params) AppRunner) error {
	return runWithParamParser(args, newRunner, nil)
}

func runWithParamParser(
	args []string,
	newRunner func(*param.Params) AppRunner,
	parseParams func(*cli.Context) (*param.Params, error),
) error {
	if newRunner == nil {
		newRunner = defaultNewRunner
	}
	if parseParams == nil {
		parseParams = param.ContextToParams
	}

	cliApp := &cli.App{
		Name:  "spa-to-http",
		Flags: param.Flags,
		Action: func(c *cli.Context) error {
			params, err := parseParams(c)
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
		logFatal(err)
	}
}
