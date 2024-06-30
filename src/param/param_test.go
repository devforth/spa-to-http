package param_test

import (
	"context"
	"flag"
	"go-http-server/param"
	"path/filepath"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestContextToParams(t *testing.T) {
	f := flag.NewFlagSet("a", flag.ContinueOnError)
	e_adress := "example adr"
	e_port := 80
	e_gzip := false
	e_brotli := true
	e_threshold := int64(1024)
	e_directory := "example dir"
	e_cache_control_max_age := int64(2048)
	e_spa := true
	e_ignore_cache_control_paths := "example path1"
	e_cache := true
	e_cache_buffer := 64

	f.String("address", e_adress, "")
	f.Int("port", e_port, "")
	f.Bool("gzip", e_gzip, "")
	f.Bool("brotli", e_brotli, "")
	f.Int64("threshold", e_threshold, "")
	f.String("directory", e_directory, "")
	f.Int64("cache-max-age", e_cache_control_max_age, "")
	f.Bool("spa", e_spa, "")
	f.String("ignore-cache-control-paths", e_ignore_cache_control_paths, "")
	f.Bool("cache", e_cache, "")
	f.Int("cache-buffer", e_cache_buffer, "")

	ctx := cli.NewContext(nil, f, nil)
	ctx.Context = context.WithValue(context.Background(), "key", "val")
	params, err := param.ContextToParams(ctx)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	if params.Address != e_adress {
		t.Errorf("Got %s, expected %s", params.Address, e_adress)
	}

	if params.Port != e_port {
		t.Errorf("Got %d, expected %d", params.Port, e_port)
	}

	if params.Gzip != e_gzip {
		t.Errorf("Got %t, expected %t", params.Gzip, e_gzip)
	}

	if params.Brotli != e_brotli {
		t.Errorf("Got %t, expected %t", params.Brotli, e_brotli)
	}

	if params.Threshold != int64(e_threshold) {
		t.Errorf("Got %d, expected %d", params.Threshold, e_threshold)
	}

	abs_directory, _ := filepath.Abs(e_directory)
	if params.Directory != abs_directory {
		t.Errorf("Got %s, expected %s", params.Directory, e_directory)
	}

	if params.CacheControlMaxAge != e_cache_control_max_age {
		t.Errorf("Got %d, expected %d", params.CacheControlMaxAge, e_cache_control_max_age)
	}

	if params.SpaMode != e_spa {
		t.Errorf("Got %t, expected %t", params.SpaMode, e_spa)
	}

	//TODO
	// fmt.Println(params.IgnoreCacheControlPaths)
	// if params.IgnoreCacheControlPaths[0] != e_ignore_cache_control_paths {
	// 	t.Errorf("Got %s, expected %s", params.IgnoreCacheControlPaths, e_ignore_cache_control_paths)
	// }

	if params.CacheEnabled != e_cache {
		t.Errorf("Got %t, expected %t", params.CacheEnabled, e_cache)
	}

	if params.CacheBuffer != e_cache_buffer {
		t.Errorf("Got %d, expected %d", params.CacheBuffer, e_cache_buffer)
	}
}
