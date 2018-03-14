package testing

import (
	"fmt"
	"io"

	"github.com/dhamidi/dux"
)

// FailingFileSystem wraps filesystem and optionally fails for
// preconfigured actions and paths.
type FailingFileSystem struct {
	dux.FileSystem

	failures map[string][]string
}

// NewFailingFileSystem wraps fs and does not fail on any action.
func NewFailingFileSystem(fs dux.FileSystem) *FailingFileSystem {
	return &FailingFileSystem{
		FileSystem: fs,
		failures:   map[string][]string{},
	}
}

// Fail causes the given action to fail for the given filename with a generic error.
//
// Valid values for action are "open", "create" and "write".
//
// - "open" causes an error to be returned when a file is Open()'d in the underlying file system
// - "create" causes an error to be returned when a file is Create()'d in the underlying file system
// - "write" causes an error when writing to the writer returned by Create()
func (fs *FailingFileSystem) Fail(action string, filename string) {
	if action != "open" && action != "create" && action != "write" {
		panic(fmt.Sprintf("Invalid action supplied to %T.Fail: %q", fs, action))
	}
	fs.failures[action] = append(fs.failures[action], filename)
}

// Create forwards the call to the underlying FileSystem, unless a failure has been registered for creating the given file.
func (fs *FailingFileSystem) Create(filename string) (io.WriteCloser, error) {
	for _, f := range fs.failures["create"] {
		if f == filename {
			return nil, fmt.Errorf("Failed to create file %q", filename)
		}
	}

	for _, f := range fs.failures["write"] {
		if f == filename {
			f, err := fs.FileSystem.Create(filename)
			if err != nil {
				return nil, err
			}
			return NewFailingWriteCloser(f), err
		}
	}

	return fs.FileSystem.Create(filename)
}

// Open forwards the call to the underlying FileSystem, unless a failure has been registered for creating the given file.
func (fs *FailingFileSystem) Open(filename string) (io.ReadCloser, error) {
	for _, f := range fs.failures["open"] {
		if f == filename {
			return nil, fmt.Errorf("Failed to open file %q", filename)
		}
	}

	return fs.FileSystem.Open(filename)
}
