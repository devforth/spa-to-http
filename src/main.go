package main

import (
	"go-http-server/app"
	"go-http-server/param"
)

func main() {
	params := param.ParseCli()

	newApp := app.NewApp(params)
	newApp.CompressFiles()
	newApp.Listen()
}
