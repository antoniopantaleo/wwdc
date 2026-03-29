package filesystem

import (
	"os"
	"path/filepath"
)

type OSFileSystem struct {
	basePath string
}

func NewOSFileSystem(basePath string) *OSFileSystem {
	return &OSFileSystem{basePath: basePath}
}

func (fs *OSFileSystem) MakeDir(path string) error {
	return os.MkdirAll(filepath.Join(fs.basePath, path), 0755)
}

func (fs *OSFileSystem) WriteFile(path string, data []byte) error {
	return os.WriteFile(filepath.Join(fs.basePath, path), data, 0644)
}
