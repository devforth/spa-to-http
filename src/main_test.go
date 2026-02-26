package main

import (
	"errors"
	"go-http-server/param"
	"os"
	"testing"
	"time"

	"github.com/urfave/cli/v2"
)

type stubRunner struct {
	compressCalled chan struct{}
	listenCalled   chan struct{}
}

func (s *stubRunner) CompressFiles() {
	close(s.compressCalled)
}

func (s *stubRunner) Listen() {
	close(s.listenCalled)
}

func TestRunCallsRunner(t *testing.T) {
	stub := &stubRunner{
		compressCalled: make(chan struct{}),
		listenCalled:   make(chan struct{}),
	}

	err := run([]string{"spa-to-http"}, func(_ *param.Params) AppRunner {
		return stub
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	select {
	case <-stub.listenCalled:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("expected Listen to be called")
	}

	select {
	case <-stub.compressCalled:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("expected CompressFiles to be called")
	}
}

func TestRunReturnsErrorOnUnknownFlag(t *testing.T) {
	err := run([]string{"spa-to-http", "--unknown-flag"}, nil)
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}

func TestRunWithParamParserError(t *testing.T) {
	expectedErr := errors.New("parse error")
	err := runWithParamParser([]string{"spa-to-http"}, func(_ *param.Params) AppRunner {
		return &stubRunner{
			compressCalled: make(chan struct{}),
			listenCalled:   make(chan struct{}),
		}
	}, func(_ *cli.Context) (*param.Params, error) {
		return nil, expectedErr
	})
	if err == nil {
		t.Fatal("expected error from parser")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestDefaultNewRunner(t *testing.T) {
	params := &param.Params{
		Directory: ".",
	}
	runner := defaultNewRunner(params)
	if runner == nil {
		t.Fatal("expected runner to be non-nil")
	}
}

func TestMainWithHelp(t *testing.T) {
	origArgs := os.Args
	t.Cleanup(func() {
		os.Args = origArgs
	})
	os.Args = []string{"spa-to-http", "--help"}
	main()
}

func TestMainLogsFatalOnError(t *testing.T) {
	origArgs := os.Args
	origLogFatal := logFatal
	t.Cleanup(func() {
		os.Args = origArgs
		logFatal = origLogFatal
	})

	os.Args = []string{"spa-to-http", "--unknown-flag"}
	called := false
	logFatal = func(_ ...any) {
		called = true
	}

	main()

	if !called {
		t.Fatal("expected logFatal to be called")
	}
}
