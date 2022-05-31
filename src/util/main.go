package util

import "os"

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
