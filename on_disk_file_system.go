package dux

import (
	"io"
	"os"
	"path/filepath"
)

// OnDiskFileSystem implements FileSystem with files from the Operating System.
type OnDiskFileSystem struct {
}

// NewOnDiskFileSystem creates a new OnDiskFileSystem anchored at the
// process's current directory.
func NewOnDiskFileSystem() *OnDiskFileSystem {
	return &OnDiskFileSystem{}
}

// Open implements FileSystem by delegating to os.Open
func (fs *OnDiskFileSystem) Open(filename string) (io.ReadCloser, error) {
	f, err := os.Open(filename)
	return f, err
}

// Create implements FileSystem by delegating to os.Create.
//
// Before creating the file, intermediate directories are created using os.MkdirAll.
//
// Directories are created with permissions 0755 and files are created
// with permissions 0644.
func (fs *OnDiskFileSystem) Create(filename string) (io.WriteCloser, error) {
	os.MkdirAll(filepath.Dir(filename), 0755)
	f, err := os.Create(filename)
	return f, err
}

// List returns the names of all files in the given directory.
//
// Unlike the command `ls`, no entries for "." and ".." are returned.
func (fs *OnDiskFileSystem) List(dir string) ([]string, error) {
	f, err := os.Open(dir)
	if err != nil {
		return []string{}, err
	}

	names, err := f.Readdirnames(0)
	return names, err
}
