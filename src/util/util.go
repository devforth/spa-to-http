package util

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

type FileType int

const (
	FileTypeNotExists = iota
	FileTypeFile
	FileTypeDirectory
)

func GetFileType(path string) FileType {
	stat, err := os.Stat(path)

	if err != nil {
		return FileTypeNotExists
	} else if stat.IsDir() {
		return FileTypeDirectory
	} else {
		return FileTypeFile
	}
}

func GetFileWithInfoAndType(path string) (http.File, fs.FileInfo, FileType) {
	dirPath, fileName := filepath.Split(path)
	fsDir := http.Dir(dirPath)
	file, err := fsDir.Open(fileName)

	if err != nil {
		return nil, nil, FileTypeNotExists
	}

	stat, err := file.Stat()

	if err != nil {
		return nil, nil, FileTypeNotExists
	} else if stat.IsDir() {
		return file, stat, FileTypeDirectory
	} else {
		return file, stat, FileTypeFile
	}
}
