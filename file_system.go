package dux

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
)

// FileSystem encodes basic operations that can be performed on a file
// system.
type FileSystem interface {
	// Open opens a file for reading.
	//
	// Any errors returned will be of type *FileSystemError
	Open(filename string) (io.ReadCloser, error)

	// Create opens a file for writing
	//
	// Any errors returned will be of type *FileSystemError
	Create(filename string) (io.WriteCloser, error)

	// List returns the names of all files in directory d
	List(dir string) ([]string, error)
}

// FileSystemError wraps errors returned by a FileSystem
type FileSystemError struct {
	Op   string
	Path string
	Err  error
}

// NewFileSystemError creates a new error for the given operation and
// path, optionally wrapping an implementation-specific error.
func NewFileSystemError(op, path string, errors ...error) *FileSystemError {
	err := &FileSystemError{
		Op:   op,
		Path: path,
		Err:  nil,
	}
	if len(errors) > 0 {
		err.Err = errors[0]
	}
	return err
}

// Error implements the error interface
func (err *FileSystemError) Error() string {
	return fmt.Sprintf("%s @ %s: %s", err.Op, err.Path, err.Err)
}

// InMemoryFileSystem implements FileSystem with buffers in RAM.
type InMemoryFileSystem struct {
	files map[string]*bytes.Buffer
}

// NewInMemoryFileSystem constructs a new file system with empty buffers
func NewInMemoryFileSystem() *InMemoryFileSystem {
	return &InMemoryFileSystem{
		files: map[string]*bytes.Buffer{},
	}
}

type nopWriteCloser struct {
	io.Writer
}

// NopWriteCloser adds do-nothing Close method to w.
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return &nopWriteCloser{Writer: w}
}

// Close does nothing and never returns an error.
func (nop *nopWriteCloser) Close() error {
	return nil
}

// Open returns a the buffer at the given path.  If no buffer is
// found, an error is returned.
func (fs *InMemoryFileSystem) Open(filename string) (io.ReadCloser, error) {
	buffer, found := fs.files[filename]
	if !found {
		return nil, NewFileSystemError("open", filename, fmt.Errorf("file not found"))
	}

	return ioutil.NopCloser(buffer), nil
}

// Create creates a new buffer at the given path.  It never returns an error
func (fs *InMemoryFileSystem) Create(filename string) (io.WriteCloser, error) {
	buffer := bytes.NewBufferString("")
	fs.files[filename] = buffer
	return NopWriteCloser(buffer), nil
}

// List returns all file names one hierarchy level below the directory d
func (fs *InMemoryFileSystem) List(d string) ([]string, error) {
	result := []string{}
	for filename := range fs.files {
		if matches, _ := filepath.Match(filepath.Join(d, "*"), filename); matches {
			result = append(result, filepath.Base(filename))
		}
	}
	return result, nil
}
