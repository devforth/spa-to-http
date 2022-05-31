package main

import (
	"go-http-server/app"
	"go-http-server/param"
)

func main() {
	params := param.ParseCli()

	newApp := app.NewApp(params)
	go newApp.CompressFiles()
	newApp.Listen()
}
