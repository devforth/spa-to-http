package param

import (
	"fmt"
	"github.com/akamensky/argparse"
	"os"
)

type CliParams struct {
	Host               string
	Port               int
	Gzip               bool
	Brotli             bool
	Threshold          int64
	Directory          string
	DirectoryListing   bool
	CacheControlMaxAge int
	SpaMode            bool
}

func (cliParams *CliParams) String() string {
	return fmt.Sprintf("CliParams{port:%d gzip:%v brotli:%v}", cliParams.Port, cliParams.Gzip, cliParams.Brotli)
}

func ParseCli() *CliParams {
	parser := argparse.NewParser("go-http-server", "Simple http server written in go for spa serving")

	host := parser.String("", "host", &argparse.Options{Default: "0.0.0.0"})
	port := parser.Int("p", "port", &argparse.Options{Default: 8080})
	gzip := parser.Flag("g", "gzip", &argparse.Options{Default: false})
	threshold := parser.Int("", "threshold", &argparse.Options{Default: 1024})
	brotli := parser.Flag("b", "brotli", &argparse.Options{Default: false})
	directory := parser.String("d", "directory", &argparse.Options{Default: "."})
	dirListing := parser.Flag("", "dir-listing", &argparse.Options{Default: false})
	cacheControlMaxAge := parser.Int("", "cache-control-max-age", &argparse.Options{Default: 604800})
	spaMode := parser.Flag("", "spa", &argparse.Options{Default: false})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	if *spaMode && *dirListing {
		panic("SPA mode and directory listing cannot be enabled at the same time")
	}

	params := CliParams{
		Host:               *host,
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
