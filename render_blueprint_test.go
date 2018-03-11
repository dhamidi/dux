package dux_test

import (
	"path/filepath"
	"testing"

	h "github.com/dhamidi/dux/testing"
)

func TestApp_RenderBlueprint_creates_the_files_specified_by_the_blueprint(t *testing.T) {
	app := h.NewApp()
	cmd := h.RenderBlueprint()
	if err := app.Execute(cmd); err != nil {
		t.Fatalf("renderBlueprint: %s", err)
	}

	for _, filename := range h.ExampleBlueprintFilenames() {
		path := filepath.Join(cmd.Destination, filename)
		_, err := app.FileSystem.Open(path)
		if err != nil {
			t.Fatalf("Expecting file %q to exist: %s", path, err)
		}
	}

}
