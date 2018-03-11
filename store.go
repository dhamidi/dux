package dux

import (
	"encoding/json"
	"path/filepath"
)

// Store provides access to a id-indexed store of persistent objects.
type Store interface {
	Get(id string, dest interface{}) error
	Put(id string, src interface{}) error
}

// FileSystemStore stores objects in the provided directory on a filesystem.
//
// Objects are de- and encoded using JSON
type FileSystemStore struct {
	fs  FileSystem
	dir string
}

// NewFileSystemStore creates a new Store for the given directory and file system.
func NewFileSystemStore(dir string, fs FileSystem) *FileSystemStore {
	return &FileSystemStore{
		fs:  fs,
		dir: dir,
	}
}

// Get deserializes the file identified by ID as JSON into dest.
func (s *FileSystemStore) Get(id string, dest interface{}) error {
	f, err := s.fs.Open(filepath.Join(s.dir, id))
	if err != nil {
		return err
	}
	defer f.Close()
	dec := json.NewDecoder(f)

	return dec.Decode(dest)
}

// Put serializes src as JSON and writes it to the file identified by id.
func (s *FileSystemStore) Put(id string, src interface{}) error {
	f, err := s.fs.Create(filepath.Join(s.dir, id))
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)

	return enc.Encode(src)
}
