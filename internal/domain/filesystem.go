package domain

type FileSystem interface {
	MakeDir(path string) error
	WriteFile(path string, data []byte) error
}