package param_test

import (
	"context"
	"errors"
	"flag"
	"go-http-server/param"
	"path/filepath"
	"reflect"
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
	e_ignore_cache_control_paths := []string{"example path1", "example path2"}
	e_cache := true
	e_cache_buffer := 64
	e_logger := true
	e_log_pretty := true
	e_no_compress := []string{".map", ".zip"}

	f.String("address", e_adress, "")
	f.Int("port", e_port, "")
	f.Bool("gzip", e_gzip, "")
	f.Bool("brotli", e_brotli, "")
	f.Int64("threshold", e_threshold, "")
	f.String("directory", e_directory, "")
	f.Int64("cache-max-age", e_cache_control_max_age, "")
	f.Bool("spa", e_spa, "")
	f.Var(cli.NewStringSlice(e_ignore_cache_control_paths...), "ignore-cache-control-paths", "")
	f.Bool("cache", e_cache, "")
	f.Int("cache-buffer", e_cache_buffer, "")
	f.Bool("logger", e_logger, "")
	f.Bool("log-pretty", e_log_pretty, "")
	f.Var(cli.NewStringSlice(e_no_compress...), "no-compress", "")

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

	if !reflect.DeepEqual(params.IgnoreCacheControlPaths, e_ignore_cache_control_paths) {
		t.Errorf("Got %v, expected %v", params.IgnoreCacheControlPaths, e_ignore_cache_control_paths)
	}

	if params.CacheEnabled != e_cache {
		t.Errorf("Got %t, expected %t", params.CacheEnabled, e_cache)
	}

	if params.CacheBuffer != e_cache_buffer {
		t.Errorf("Got %d, expected %d", params.CacheBuffer, e_cache_buffer)
	}

	if params.Logger != e_logger {
		t.Errorf("Got %t, expected %t", params.Logger, e_logger)
	}

	if params.LogPretty != e_log_pretty {
		t.Errorf("Got %t, expected %t", params.LogPretty, e_log_pretty)
	}

	if !reflect.DeepEqual(params.NoCompress, e_no_compress) {
		t.Errorf("Got %v, expected %v", params.NoCompress, e_no_compress)
	}
}

func TestContextToParamsWithAbsError(t *testing.T) {
	f := flag.NewFlagSet("a", flag.ContinueOnError)
	f.String("directory", "example", "")

	ctx := cli.NewContext(nil, f, nil)
	ctx.Context = context.WithValue(context.Background(), "key", "val")

	expectedErr := errors.New("abs failed")
	params, err := param.ContextToParamsWithAbs(ctx, func(_ string) (string, error) {
		return "", expectedErr
	})
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, got %v", expectedErr, err)
	}
	if params != nil {
		t.Fatalf("Expected params to be nil on error, got %v", params)
	}
}
