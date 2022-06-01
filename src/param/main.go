package param

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/caarlos0/env/v6"
	"golang.org/x/exp/slices"
	"gopkg.in/guregu/null.v4"
	"os"
	"strings"
)

type Params struct {
	Address            string    `env:"ADDRESS"`
	Port               int       `env:"PORT"`
	Gzip               null.Bool `env:"GZIP"`
	Brotli             null.Bool `env:"BROTLI"`
	Threshold          int64     `env:"THRESHOLD"`
	Directory          string    `env:"DIRECTORY"`
	DirectoryListing   null.Bool `env:"DIR_LISTING"`
	CacheControlMaxAge int       `env:"CACHE_MAX_AGE"`
	SpaMode            null.Bool `env:"SPA_MODE"`
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
	spaModeString := parser.String("", "spa", &argparse.Options{})

	args := os.Args
	spaMode := null.NewBool(false, false)

	if slices.Contains(args, "--spa") {
		index := slices.Index(args, "--spa")

		if len(args) <= index+1 || strings.Index(args[index+1], "-") == 0 {
			args = slices.Delete(args, index, index+1)

			spaMode = null.NewBool(true, true)
		}
	}

	fmt.Println(args)

	err := parser.Parse(args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if *spaModeString != "" {
		if *spaModeString == "false" || *spaModeString == "0" {
			spaMode = null.NewBool(false, true)
		} else {
			spaMode = null.NewBool(true, true)
		}
	}

	if spaMode.ValueOrZero() && *dirListing {
		panic("SPA mode and directory listing cannot be enabled at the same time")
	}

	params := Params{
		Address:            *host,
		Port:               *port,
		Gzip:               null.BoolFromPtr(gzip),
		Brotli:             null.BoolFromPtr(brotli),
		Threshold:          int64(*threshold),
		Directory:          *directory,
		DirectoryListing:   null.BoolFromPtr(dirListing),
		CacheControlMaxAge: *cacheControlMaxAge,
		SpaMode:            spaMode,
	}

	return &params
}

func parseEnv() *Params {
	params := Params{}

	err := env.Parse(&params)
	if err != nil {
		panic(err)
	}

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
		Gzip:               coalesce(cliParams.Gzip, envParams.Gzip, null.BoolFrom(false)),
		Brotli:             coalesce(cliParams.Brotli, envParams.Brotli, null.BoolFrom(false)),
		Threshold:          coalesce(cliParams.Threshold, envParams.Threshold, 1024),
		Directory:          coalesce(cliParams.Directory, envParams.Directory, "."),
		DirectoryListing:   coalesce(cliParams.DirectoryListing, envParams.DirectoryListing, null.BoolFrom(false)),
		CacheControlMaxAge: coalesce(cliParams.CacheControlMaxAge, envParams.CacheControlMaxAge, 604800),
		SpaMode:            coalesce(cliParams.SpaMode, envParams.SpaMode, null.BoolFrom(true)),
	}
}
