package param

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/caarlos0/env/v6"
	"os"
)

type Params struct {
	Address            string `env:"ADDRESS"`
	Port               int    `env:"PORT"`
	Gzip               bool   `env:"GZIP"`
	Brotli             bool   `env:"BROTLI"`
	Threshold          int64  `env:"THRESHOLD"`
	Directory          string `env:"DIRECTORY"`
	DirectoryListing   bool   `env:"DIR_LISTING"`
	CacheControlMaxAge int    `env:"CACHE_MAX_AGE"`
	SpaMode            bool   `env:"SPA_MODE"`
}

func parseCli() *Params {
	parser := argparse.NewParser("go-http-server", "Simple http server written in go for spa serving")

	host := parser.String("a", "address", &argparse.Options{})
	port := parser.Int("p", "port", &argparse.Options{})
	gzip := parser.Flag("g", "gzip", &argparse.Options{})
	brotli := parser.Flag("b", "brotli", &argparse.Options{})
	threshold := parser.Int("", "threshold", &argparse.Options{})
	directory := parser.String("d", "directory", &argparse.Options{})
	dirListing := parser.Flag("", "dir-listing", &argparse.Options{})
	cacheControlMaxAge := parser.Int("", "cache-max-age", &argparse.Options{})
	spaMode := parser.Flag("", "spa", &argparse.Options{})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	if *spaMode && *dirListing {
		panic("SPA mode and directory listing cannot be enabled at the same time")
	}

	params := Params{
		Address:            *host,
		Port:               *port,
		Gzip:               *gzip,
		Brotli:             *brotli,
		Threshold:          int64(*threshold),
		Directory:          *directory,
		DirectoryListing:   *dirListing,
		CacheControlMaxAge: *cacheControlMaxAge,
		SpaMode:            *spaMode,
	}

	return &params
}

func parseEnv() *Params {
	params := Params{}

	err := env.Parse(&params)
	if err != nil {
		panic(err)
	}

	fmt.Println(params)

	return &params
}

func coalesce[T comparable](items ...T) T {
	var result T

	for _, item := range items {
		if item != result {
			return item
		}
	}
	return result
}

func GetParams() *Params {
	envParams := parseEnv()
	cliParams := parseCli()

	return &Params{
		Address:            coalesce(cliParams.Address, envParams.Address, "0.0.0.0"),
		Port:               coalesce(cliParams.Port, envParams.Port, 8080),
		Gzip:               coalesce(cliParams.Gzip, envParams.Gzip, false),
		Brotli:             coalesce(cliParams.Brotli, envParams.Brotli, false),
		Threshold:          coalesce(cliParams.Threshold, envParams.Threshold, 1024),
		Directory:          coalesce(cliParams.Directory, envParams.Directory, "."),
		DirectoryListing:   coalesce(cliParams.DirectoryListing, envParams.DirectoryListing, false),
		CacheControlMaxAge: coalesce(cliParams.CacheControlMaxAge, envParams.CacheControlMaxAge, 604800),
		SpaMode:            coalesce(cliParams.SpaMode, envParams.SpaMode, false),
	}
}
