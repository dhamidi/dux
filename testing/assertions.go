package testing

import (
	"io/ioutil"
	"testing"

	"github.com/dhamidi/dux"
)

func AssertFileContents(t *testing.T, fs dux.FileSystem, path string, expectedContents string) {
	t.Helper()
	f, err := fs.Open(path)
	if err != nil {
		t.Fatalf("AssertFileContents: %s", err)
	}
	defer f.Close()
	contents, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("AssertFileContents: %s", err)
	}

	if string(contents) != expectedContents {
		t.Fatalf("File %q expected contents:\n%s\n---\nActual contents:\n%s\n",
			path, expectedContents, contents)
	}
}
