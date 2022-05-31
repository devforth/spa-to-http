package main

import (
	"go-http-server/app"
	"go-http-server/param"
)

func main() {
	params := param.GetParams()

	newApp := app.NewApp(params)
	go newApp.CompressFiles()
	newApp.Listen()
}
