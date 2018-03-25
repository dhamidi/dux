package dux

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Store provides access to a id-indexed store of persistent objects.
type Store interface {
	Get(id string, dest interface{}) error
	Put(id string, src interface{}) error
	List(pattern string) ([]string, error)
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
	f, err := s.fs.Open(filepath.Join(s.dir, id) + ".json")
	if err != nil {
		return err
	}
	defer f.Close()
	dec := json.NewDecoder(f)

	return dec.Decode(dest)
}

// Put serializes src as JSON and writes it to the file identified by id.
func (s *FileSystemStore) Put(id string, src interface{}) error {
	f, err := s.fs.Create(filepath.Join(s.dir, id) + ".json")
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)

	return enc.Encode(src)
}

// List searches for all json files matching the given glob in the base directory of the file system store.
//
// The json extension is removed from all filenames before they are returned
func (s *FileSystemStore) List(pattern string) ([]string, error) {
	filenames, err := filepath.Glob(filepath.Join(s.dir, pattern+".json"))
	result := []string{}
	if err != nil {
		return result, err
	}

	prefix := s.dir + string(os.PathSeparator)
	for _, filename := range filenames {
		withoutExtension := strings.TrimSuffix(filename, ".json")
		withoutDirectory := strings.TrimPrefix(withoutExtension, prefix)
		result = append(result, withoutDirectory)
	}

	return result, nil
}
