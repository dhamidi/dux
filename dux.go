package dux

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Environment provides read-only access to a key value store mapping strings to strings.
type Environment func(string) string

// SystemEnvironment provides access to the operating system's environment variables.
var (
	SystemEnvironment = os.Getenv
)

// Context provides information about the environment in which Dux runs.
type Context struct {
	BaseDir   string
	OutputDir string
	Logger    *log.Logger
	Env       Environment
	Data      interface{}
	App       *Application
}

// NewContextFromEnvironment creates a new context with settings taking from the given environment.
func NewContextFromEnvironment(env Environment) *Context {
	result := &Context{
		BaseDir:   env("PWD"),
		OutputDir: filepath.Join(env("PWD"), "tmp"),
		Logger:    log.New(os.Stdout, "", 0),
		Env:       env,
	}
	result.App = NewApplication("dux", result)
	return result
}

// Log outputs a log message for the given module.
func (ctx *Context) Log(module, kind, message string) {
	ctx.Logger.Printf("%10s %10s %s", module, kind, message)
}

// MergeJSON deeply merges two raw deserialized JSON objects.
func MergeJSON(a, b interface{}) interface{} {
	result := map[string]interface{}{}

	if a == nil && b != nil {
		return b
	}

	aMap, aMapOk := a.(map[string]interface{})
	bMap, bMapOk := b.(map[string]interface{})
	if aMapOk && bMapOk {
		for key, value := range aMap {
			result[key] = value
		}
		for key, value := range bMap {
			result[key] = MergeJSON(result[key], value)
		}
		return result
	} else if !aMapOk && !bMapOk {
		return b
	}

	return a
}

// GatherData loads template data from the file system into the current context.
func (ctx *Context) GatherData(datafile string) error {
	filename := filepath.Join(ctx.BaseDir, datafile)
	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	result := map[string]interface{}{}
	if err := json.NewDecoder(file).Decode(&result); err != nil {
		return err
	}

	ctx.Data = MergeJSON(ctx.Data, result)
	return nil
}

// BlueprintArgument encodes the data about the arguments accepted by the blueprint.
type BlueprintArgument struct {
	Name string
	Type string
	Doc  *string
}

// BlueprintFileDescription describes which target file should be generated from which template.
type BlueprintFileDescription struct {
	// Template is the name of the template to use.
	Template string

	// Destination is the a template expression that is evaluated
	// to obtain the destination file name.
	Destination string
}

func (fd *BlueprintFileDescription) Render(ctx *Context, blueprint *Blueprint, data interface{}) (string, error) {
	destinationFileTemplate, err := template.New("main").Parse(fd.Destination)
	if err != nil {
		return fd.Destination, fmt.Errorf("cannot parse destination file template %q: %s", fd.Destination, err)
	}
	destinationFileName := bytes.NewBufferString("")
	if err := destinationFileTemplate.Execute(destinationFileName, data); err != nil {
		return fd.Destination, fmt.Errorf("failed to generated filename from %q: %s", fd.Destination, err)
	}

	t, err := blueprint.templates.Clone()
	if err != nil {
		return destinationFileName.String(), fmt.Errorf("failed to clone templates: %s", err)
	}
	outputFilePath := filepath.Join(ctx.OutputDir, destinationFileName.String())
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return destinationFileName.String(), fmt.Errorf("failed to create output file %q: %s", outputFilePath, err)
	}
	defer outputFile.Close()
	if err := t.Execute(outputFile, data); err != nil {
		return destinationFileName.String(), fmt.Errorf("error rendering template %q: %s", outputFilePath, err)
	}
	ctx.Log(blueprint.Name, "create", destinationFileName.String())
	return destinationFileName.String(), nil
}

// Blueprint represents a set of files that need to be created on
// disk, followed by performing editing operations on existing files.
type Blueprint struct {
	// path is the location in which the blueprint is defined.
	path string

	// templates holds all templates used by this blueprint
	templates *template.Template

	// flags accepted by this blueprint
	flags *flag.FlagSet

	// Data holds the data that is provided in the blueprint
	// manifest or extracted from flags provided to the blueprint.
	Data map[string]interface{}

	// Name is the basename of the directory containing the blueprint
	Name string `json:"name"`

	// Description contains a short explanation of what the blueprint is about.
	Description string `json:"description"`

	// Args is a list of arguments that are accepted by this blueprint
	Args []*BlueprintArgument `json:"args"`

	// Files describes the list of files that should be generated
	// for this blueprint.
	Files []*BlueprintFileDescription `json:"files"`
}

// NewBlueprint creates a new blueprint instance.
func NewBlueprint(path string) *Blueprint {
	return &Blueprint{
		path:      path,
		templates: nil,
		Data:      map[string]interface{}{},
	}
}

// BlueprintRenderResult describes the results of rendering a
// blueprint.  In particular, it keeps a list of rendered files.
type BlueprintRenderResult struct {
	Files  []string
	Errors map[string]error
}

// NewBlueprintRenderResult initializes and empty result object.
func NewBlueprintRenderResult() *BlueprintRenderResult {
	return &BlueprintRenderResult{
		Files:  []string{},
		Errors: map[string]error{},
	}
}

// AddFile records the result of rendering a file.
func (r *BlueprintRenderResult) AddFile(path string, err error) *BlueprintRenderResult {
	r.Files = append(r.Files, path)
	if err != nil {
		r.Errors[path] = err
	}
	return r
}

// HasError returns true if any errors have been recorded in this result
func (r *BlueprintRenderResult) HasError() bool {
	return len(r.Errors) > 0
}

// Error implements the error interface
func (r *BlueprintRenderResult) Error() string {
	out := bytes.NewBufferString("")
	for path, err := range r.Errors {
		fmt.Fprintf(out, "- %s: %s\n", path, err)
	}
	return out.String()
}

// Render renders all templates of the blueprint
func (bp *Blueprint) Render(ctx *Context) *BlueprintRenderResult {
	data := MergeJSON(ctx.Data, bp.Data)
	result := NewBlueprintRenderResult()
	for _, filedesc := range bp.Files {
		path, err := filedesc.Render(ctx, bp, data)
		if err != nil {
			ctx.Log(bp.Name, "error", err.Error())
		}
		result.AddFile(path, err)
	}
	return result
}

// LoadTemplates parses all of the blueprint's templates
func (bp *Blueprint) LoadTemplates() error {
	templates, err := template.ParseGlob(bp.TemplateDir() + "/*")
	if err != nil {
		return err
	}
	bp.templates = templates
	return nil
}

// LoadManifest parses the blueprints serialized description from JSON.
func (bp *Blueprint) LoadManifest() error {
	manifestFilename := filepath.Join(bp.path, "manifest.json")
	manifest, err := os.Open(manifestFilename)
	if err != nil {
		return err
	}
	return json.NewDecoder(manifest).Decode(bp)
}

// ValueTypes maps strings to value constructors
var ValueTypes = map[string]func() interface{}{
	"integer": func() interface{} {
		return int64(0)
	},
	"string": func() interface{} {
		return ""
	},
	"float": func() interface{} {
		return float64(0.0)
	},
	"bool": func() interface{} {
		return false
	},
	"duration": func() interface{} {
		return time.Duration(0)
	},
	"identifier": func() interface{} {
		return &Identifier{}
	},
}

// DefineFlags initializes the flag set for the blueprint based on the
// defined blueprint arguments.
func (bp *Blueprint) DefineFlags() error {
	bp.flags = flag.NewFlagSet(bp.Name, flag.ContinueOnError)
	for _, argument := range bp.Args {
		newValue, found := ValueTypes[argument.Type]
		if !found {
			return fmt.Errorf("Unknown argument type: %s", argument.Type)
		}
		rawValue := newValue()
		doc := ""
		if argument.Doc != nil {
			doc = *argument.Doc
		}
		switch v := rawValue.(type) {
		case int64:
			bp.flags.Int64(argument.Name, 0, doc)
		case float64:
			bp.flags.Float64(argument.Name, 0.0, doc)
		case string:
			bp.flags.String(argument.Name, "", doc)
		case time.Duration:
			bp.flags.Duration(argument.Name, 0*time.Millisecond, doc)
		case bool:
			bp.flags.Bool(argument.Name, false, doc)
		default:
			flagValue, ok := v.(flag.Value)
			if !ok {
				return fmt.Errorf("Cannot handle flag value of type %T", v)
			}
			bp.flags.Var(flagValue, argument.Name, doc)
		}

	}

	return nil
}

// ParseArgs parses the provided arguments according to the blueprint's flag set.
func (bp *Blueprint) ParseArgs(args []string) error {
	if err := bp.flags.Parse(args); err != nil {
		return err
	}
	bp.flags.VisitAll(func(f *flag.Flag) {
		bp.Data[f.Name] = f.Value.(flag.Getter).Get()
	})

	return nil
}

// TemplateDir returns the path to the directory that holds the templates for this blueprint.
func (bp *Blueprint) TemplateDir() string {
	return filepath.Join(bp.path, "templates")
}

// CopyFilesToDestination copies the files generated by this blueprint
// from the staging directory to the destination.  Existing files are
// skipped.
func (bp *Blueprint) CopyFilesToDestination(ctx *Context, result *BlueprintRenderResult) error {
	for _, stagingFilePath := range result.Files {
		bp.copyFileToDestination(ctx, stagingFilePath)

	}
	return nil
}

func (bp *Blueprint) copyFileToDestination(ctx *Context, fileName string) {
	destinationFilePath := filepath.Join(ctx.BaseDir, fileName)
	stagingFilePath := filepath.Join(ctx.OutputDir, fileName)
	destinationFileName := strings.TrimPrefix(destinationFilePath, ctx.BaseDir+"/")
	dest, err := os.OpenFile(destinationFilePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if os.IsExist(err) {
		ctx.Log(bp.Name, "skip", destinationFileName)
		return
	}
	defer dest.Close()
	src, err := os.Open(stagingFilePath)
	if err != nil {
		ctx.Log(bp.Name, "error", err.Error())
		return
	}
	defer src.Close()
	if _, err := io.Copy(dest, src); err != nil {
		ctx.Log(bp.Name, "error", err.Error())
		dest.Close()
		os.Remove(destinationFilePath)
	}
}

// LoadBlueprint loads a blueprint from disk.
func (ctx *Context) LoadBlueprint(blueprintName string) (*Blueprint, error) {
	directoryName := filepath.Join(ctx.BaseDir, "blueprints", blueprintName)
	blueprint := NewBlueprint(directoryName)
	if err := blueprint.LoadTemplates(); err != nil {
		return nil, err
	}
	if err := blueprint.LoadManifest(); err != nil {
		return nil, err
	}
	if err := blueprint.DefineFlags(); err != nil {
		return nil, err
	}
	return blueprint, nil
}
