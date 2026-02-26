package param

import (
	"github.com/urfave/cli/v2"
	"path/filepath"
	"strings"
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
	&cli.BoolFlag{
		EnvVars: []string{"LOGGER"},
		Name:    "logger",
		Value:   false,
	},
	&cli.BoolFlag{
		EnvVars: []string{"LOG_PRETTY"},
		Name:    "log-pretty",
		Value:   false,
	},
	&cli.StringSliceFlag{
		EnvVars: []string{"NO_COMPRESS"},
		Name:    "no-compress",
		Value:   nil,
	},
	&cli.StringFlag{
		EnvVars: []string{"BASIC_AUTH"},
		Name:    "basic-auth",
		Value:   "",
	},
	&cli.StringFlag{
		EnvVars: []string{"BASIC_AUTH_REALM"},
		Name:    "basic-auth-realm",
		Value:   "",
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
	Logger                  bool
	LogPretty               bool
	NoCompress              []string
	BasicAuthRaw            string
	BasicAuthUser           string
	BasicAuthPass           string
	BasicAuthRealm          string
	BasicAuthEnabled        bool
	//DirectoryListing        bool
}

func ContextToParams(c *cli.Context) (*Params, error) {
	return ContextToParamsWithAbs(c, filepath.Abs)
}

func ContextToParamsWithAbs(c *cli.Context, abs func(string) (string, error)) (*Params, error) {
	directory, err := abs(c.String("directory"))
	if err != nil {
		return nil, err
	}

	basicAuthRaw := c.String("basic-auth")
	basicAuthRealm := c.String("basic-auth-realm")
	var basicAuthUser string
	var basicAuthPass string
	var basicAuthEnabled bool

	if basicAuthRaw != "" {
		user, pass, ok := strings.Cut(basicAuthRaw, ":")
		if !ok || user == "" || pass == "" {
			return nil, cli.Exit("invalid basic-auth format, expected username:password", 1)
		}
		basicAuthUser = user
		basicAuthPass = pass
		basicAuthEnabled = true
		if basicAuthRealm == "" {
			basicAuthRealm = "Restricted"
		}
	}

	return &Params{
		Address:                 c.String("address"),
		Port:                    c.Int("port"),
		Gzip:                    c.Bool("gzip"),
		Brotli:                  c.Bool("brotli"),
		Threshold:               c.Int64("threshold"),
		Directory:               directory,
		CacheControlMaxAge:      c.Int64("cache-max-age"),
		SpaMode:                 c.Bool("spa"),
		IgnoreCacheControlPaths: c.StringSlice("ignore-cache-control-paths"),
		CacheEnabled:            c.Bool("cache"),
		CacheBuffer:             c.Int("cache-buffer"),
		Logger:                  c.Bool("logger"),
		LogPretty:               c.Bool("log-pretty"),
		NoCompress:              c.StringSlice("no-compress"),
		BasicAuthRaw:            basicAuthRaw,
		BasicAuthUser:           basicAuthUser,
		BasicAuthPass:           basicAuthPass,
		BasicAuthRealm:          basicAuthRealm,
		BasicAuthEnabled:        basicAuthEnabled,
		//DirectoryListing:        c.Bool("directory-listing"),
	}, nil
}
