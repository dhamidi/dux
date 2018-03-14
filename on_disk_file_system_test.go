package dux_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/dhamidi/dux"
)

func TestOnDiskFileSystem_CreateAndOpen(t *testing.T) {
	os.RemoveAll("a")
	defer os.RemoveAll("a")
	fs := dux.NewOnDiskFileSystem()
	f, err := fs.Create("a/b/c")
	if err != nil {
		t.Fatalf("fs.Create: %s", err)
	}
	fmt.Fprintf(f, "hello, world")
	f.Close()

	in, err := fs.Open("a/b/c")
	if err != nil {
		t.Fatalf("fs.Open: %s", err)
	}
	defer in.Close()
	contents, err := ioutil.ReadAll(in)
	if err != nil {
		t.Fatalf("ioutil.ReadAll: %s", err)
	}
	in.Close()

	if got, want := string(contents), "hello, world"; got != want {
		t.Fatalf("Expected file contents %q, got %q", want, got)
	}
}

func TestOnDiskFileSystem_List_does_not_return_dot_nor_dotdot(t *testing.T) {
	fs := dux.NewOnDiskFileSystem()
	names, err := fs.List(".")
	if err != nil {
		t.Fatalf("fs.List: %s", err)
	}

	for _, name := range names {
		if name == "." || name == ".." {
			t.Fatalf("Invalid name %q found in %#v", name, names)
		}
	}
}
