package testing

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/dhamidi/dux"
)

func AssertFileContents(t *testing.T, fs dux.FileSystem, path string, expectedContents string) {
	t.Helper()
	f, err := fs.Open(path)
	if err != nil {
		dir := filepath.Dir(path)
		filesInDirectory, listErr := fs.List(dir)
		message := bytes.NewBufferString(fmt.Sprintf("AssertFileContents: %s", err))
		if listErr == nil {
			fmt.Fprintf(message, "\nFiles in filesystem:\n")
			for _, f := range filesInDirectory {
				fmt.Fprintf(message, "- %s/%s\n", dir, f)
			}
		}

		t.Fatalf(message.String())
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
