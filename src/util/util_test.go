package util_test

import (
	"errors"
	"go-http-server/util"
	"io"
	"io/fs"
	"net/http"
	"testing"
)

func TestGetFileWithInfoAndType(t *testing.T) {
	f1, f_info1, f_type1 := util.GetFileWithInfoAndType("afadf")
	if f1 != nil {
		t.Errorf("Expected f1 = nil, got f1 = %s", f1)
	}
	if f_info1 != nil {
		t.Errorf("Expected f_info1 = nil, got f_info1 = %s", f_info1)
	}
	if f_type1 != 0 {
		t.Errorf("Expected f_type1 = 0, got f_type1 = %d", f_info1)
	}

	f2, f_info2, f_type2 := util.GetFileWithInfoAndType("util.go")
	if f2 == nil {
		t.Errorf("Expected f2 != nil, got f2 = nil")
	}
	if f_info2.Name() != "util.go" {
		t.Errorf("Expected f_info2.Name() == 'util.go', got f_info2.Name() = nil")
	}
	if f_type2 != 1 {
		t.Errorf("Expected f_type2 = 1, got f_type1 = %d", f_info1)
	}

	f3, f_info3, f_type3 := util.GetFileWithInfoAndType("../util")
	if f3 == nil {
		t.Errorf("Expected f3 != nil, got f3 = nil")
	}
	if f_info3.Name() != "util" {
		t.Errorf("Expected f_info3.Name() == 'util', got f_info3.Name() = nil")
	}
	if f_type3 != 2 {
		t.Errorf("Expected f_type3 = 2, got f_type3 = %d", f_info1)
	}
}

func TestGetFileType(t *testing.T) {
	table := []struct {
		name     string
		input    string
		expected util.FileType
	}{
		{"wrong dir input", "../dasds", 0},
		{"exsisting file", "util.go", 1},
		{"directory", "../util", 2},
	}

	for _, tc := range table {
		actual := util.GetFileType(tc.input)
		if actual != tc.expected {
			t.Error("expected: ", tc.expected, "\ngot: ", actual)
		}
	}
}

type badStatFile struct{}

func (f *badStatFile) Close() error                       { return nil }
func (f *badStatFile) Read(_ []byte) (int, error)         { return 0, io.EOF }
func (f *badStatFile) Seek(_ int64, _ int) (int64, error) { return 0, nil }
func (f *badStatFile) Readdir(_ int) ([]fs.FileInfo, error) {
	return nil, nil
}
func (f *badStatFile) Stat() (fs.FileInfo, error) {
	return nil, errors.New("stat failed")
}

func TestGetFileWithInfoAndTypeWithOpenerStatError(t *testing.T) {
	opener := func(_, _ string) (http.File, error) {
		return &badStatFile{}, nil
	}

	file, info, fileType := util.GetFileWithInfoAndTypeWithOpener("irrelevant", opener)
	if file != nil {
		t.Errorf("Expected file = nil, got %v", file)
	}
	if info != nil {
		t.Errorf("Expected info = nil, got %v", info)
	}
	if fileType != util.FileTypeNotExists {
		t.Errorf("Expected FileTypeNotExists, got %v", fileType)
	}
}
