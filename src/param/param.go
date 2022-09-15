package param

import (
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		EnvVars: []string{"ADDRESS"},
		Name:    "address",
		Aliases: []string{"a"},
		Value:   "0.0.0.0",
	},
	&cli.IntFlag{
		EnvVars: []string{"PORT"},
		Name:    "port",
		Aliases: []string{"p"},
		Value:   8080,
	},
	&cli.BoolFlag{
		EnvVars: []string{"GZIP"},
		Name:    "gzip",
		Value:   false,
	},
	&cli.BoolFlag{
		EnvVars: []string{"BROTLI"},
		Name:    "brotli",
		Value:   false,
	},
	&cli.Int64Flag{
		EnvVars: []string{"THRESHOLD"},
		Name:    "threshold",
		Value:   1024,
	},
	&cli.StringFlag{
		EnvVars: []string{"DIRECTORY"},
		Name:    "directory",
		Aliases: []string{"d"},
		Value:   ".",
	},
	// TODO
	//&cli.BoolFlag{
	//	EnvVars: []string{"DIRECTORY_LISTING"},
	//	Name:    "directory-listing",
	//	Value:   false,
	//},
	&cli.Int64Flag{
		EnvVars: []string{"CACHE_MAX_AGE"},
		Name:    "cache-max-age",
		Value:   604800,
	},
	&cli.BoolFlag{
		EnvVars: []string{"SPA_MODE"},
		Name:    "spa",
		Value:   true,
	},
	&cli.StringSliceFlag{
		EnvVars: []string{"IGNORE_CACHE_CONTROL_PATHS"},
		Name:    "ignore-cache-control-paths",
		Value:   nil,
	},
	&cli.BoolFlag{
		EnvVars: []string{"CACHE"},
		Name:    "cache",
		Value:   true,
	},
	&cli.IntFlag{
		EnvVars: []string{"CACHE_BUFFER"},
		Name:    "cache-buffer",
		Value:   50 * 1024,
	},
}

type Params struct {
	Address                 string
	Port                    int
	Gzip                    bool
	Brotli                  bool
	Threshold               int64
	Directory               string
	CacheControlMaxAge      int64
	SpaMode                 bool
	IgnoreCacheControlPaths []string
	CacheEnabled            bool
	CacheBuffer             int
	//DirectoryListing        bool

}

func ContextToParams(c *cli.Context) *Params {
	return &Params{
		Address:                 c.String("address"),
		Port:                    c.Int("port"),
		Gzip:                    c.Bool("gzip"),
		Brotli:                  c.Bool("brotli"),
		Threshold:               c.Int64("threshold"),
		Directory:               c.String("directory"),
		CacheControlMaxAge:      c.Int64("cache-control-max-age"),
		SpaMode:                 c.Bool("spa"),
		IgnoreCacheControlPaths: c.StringSlice("ignore-cache-control-paths"),
		CacheEnabled:            c.Bool("cache"),
		CacheBuffer:             c.Int("cache-buffer"),
		//DirectoryListing:        c.Bool("directory-listing"),
	}
}
